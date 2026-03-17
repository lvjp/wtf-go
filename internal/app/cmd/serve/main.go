package serve

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/lvjp/wtf-go/internal/app/api/misc"
	"github.com/lvjp/wtf-go/internal/pkg/cmd/util"
	"github.com/lvjp/wtf-go/pkg/buildinfo"

	fiberzerolog "github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

func Run(ctx *util.Context) error {
	var cancel context.CancelFunc
	ctx.Context, cancel = context.WithCancel(ctx.Context)
	defer cancel()

	server := newFiberApp(&ctx.Logger)

	apiGroup := server.Group("/api/v0")
	misc.Route(apiGroup.Group("/misc"), misc.NewService())

	prometheus.MustRegister(newCollector())
	server.Get("/metrics", promhttp.Handler())

	var serverErr error
	go func() {
		defer cancel()

		serverErr = server.Listen(*ctx.Config.Server.ListenAddress)
	}()

	<-ctx.Done()
	ctx.Logger.Info().Msg("Server shutdown sequence started")

	if serverErr != nil && !errors.Is(serverErr, http.ErrServerClosed) {
		return fmt.Errorf("ListenAndServe error: %v", serverErr)
	}

	if err := server.Shutdown(); err != nil {
		return fmt.Errorf("could not shutdown server: %v", err)
	}

	return nil
}

func newFiberApp(logger *zerolog.Logger) *fiber.App {
	app := fiber.New()

	app.Hooks().OnListen(func(listenData fiber.ListenData) error {
		u := url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(listenData.Host, listenData.Port),
		}

		if listenData.TLS {
			u.Scheme = "https"
		}

		logger.Info().
			Stringer("endpoint", &u).
			Msg("Listening")

		return nil
	})

	app.Hooks().OnPostShutdown(func(err error) error {
		logger.Info().Err(err).Msg("Fiber shutdown done")
		return nil
	})

	app.Use(requestid.New())

	app.Use(func(c fiber.Ctx) error {
		ctx := logger.With().
			Str("requestId", requestid.FromContext(c)).
			Logger().
			WithContext(c.Context())
		c.SetContext(ctx)

		return c.Next()
	})

	app.Use(fiberzerolog.New(fiberzerolog.Config{
		GetLogger: func(c fiber.Ctx) zerolog.Logger {
			return *zerolog.Ctx(c.Context())
		},
	}))

	app.Use(cors.New(cors.Config{
		AllowOrigins:        []string{"*"},
		AllowMethods:        []string{fiber.MethodGet},
		ExposeHeaders:       []string{"X-Request-ID"},
		AllowPrivateNetwork: true,
	}))

	return app
}

func newCollector() prometheus.Collector {
	bi := buildinfo.Get()

	info := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "info",
		Help: "Information about wtf-go build",
		ConstLabels: prometheus.Labels{
			"revision":      bi.Revision,
			"revision_time": bi.RevisionTime,
			"modified":      strconv.FormatBool(bi.Modified),
		},
	})
	info.Set(1)

	start_date := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "start_date_timestamp",
		Help: "The date on which the server started expressed as an UTC Unix timestamp",
	})
	start_date.SetToCurrentTime()

	registry := prometheus.NewRegistry()
	registry.MustRegister(
		info,
		start_date,
	)

	return prometheus.WrapCollectorWithPrefix("wtf_go_", registry)
}
