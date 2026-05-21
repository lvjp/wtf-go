package serve

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/lvjp/wtf-go/internal/app/config"
	"github.com/lvjp/wtf-go/internal/pkg/cmd/util"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(
		m,
		goleak.IgnoreAnyFunction("github.com/valyala/fasthttp.updateServerDate.func1"),
		goleak.IgnoreAnyFunction("github.com/valyala/fasthttp.(*workerPool).Start.func2"),
	)
}

// Run registers a collector on the global Prometheus registry, so it can only
// be invoked once per process; this single test exercises the concurrent
// shutdown path that previously raced on the listen error.
func TestRun_gracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	uctx := &util.Context{
		Context: ctx,
		Logger:  zerolog.Nop(),
		Config: &config.Config{
			Server: config.Server{ListenAddress: "127.0.0.1:"},
			Auth:   config.Auth{TokenTTL: time.Hour},
		},
	}

	ready := make(chan string, 1)
	runErr := make(chan error, 1)
	go func() {
		runErr <- Run(uctx, ready)
	}()

	select {
	case endpoint := <-ready:
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"/metrics", nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	case <-time.After(5 * time.Second):
		t.Fatal("server never became ready")
	}

	cancel()

	select {
	case err := <-runErr:
		require.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}
