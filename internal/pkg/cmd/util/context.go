package util

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
	"os"
	"time"

	"github.com/lvjp/wtf-go/internal/app/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type Context struct {
	context.Context

	Input  io.Reader
	Output io.Writer
	Error  io.Writer
	Logger zerolog.Logger

	Config *config.Config
}

// NewContext will print error on os.Stderr and exit with code 1 if any error occurs during
// initialization.
func NewContext(cmd *cobra.Command, configPath string, verbose bool) *Context {
	ret := &Context{
		Context: cmd.Context(),

		Input:  cmd.InOrStdin(),
		Output: cmd.OutOrStdout(),
		Error:  cmd.ErrOrStderr(),
	}

	ret.CheckErr(ret.initConfig(configPath, verbose), 1)
	ret.CheckErr(ret.initLogger(), 1)

	return ret
}

func (ctx *Context) CheckErr(err error, code int) {
	if err == nil {
		return
	}

	fmt.Fprintln(ctx.Error, "Error:", err)
	os.Exit(code)
}

func (ctx *Context) initConfig(configPath string, verbose bool) error {
	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		return err
	}

	if verbose {
		cfg.Log.Level = "debug"
	}

	ctx.Config = cfg
	return nil
}

func (ctx *Context) initLogger() error {
	writer := ctx.Error

	var unknowFormat bool

	switch ctx.Config.Log.Format {
	case "json":
		// default is json, do nothing
	case "console":
		writer = zerolog.ConsoleWriter{
			Out:        writer,
			TimeFormat: time.RFC3339,
		}
	default:
		unknowFormat = true
	}

	level, err := zerolog.ParseLevel(ctx.Config.Log.Level)
	if err != nil {
		return fmt.Errorf("log level parsing: %v", err)
	}

	ctx.Logger = zerolog.New(writer).With().Timestamp().Logger()
	ctx.Context = ctx.Logger.WithContext(ctx.Context)

	log.Logger = ctx.Logger.With().Str("component", "default logger").Logger()
	zerolog.DefaultContextLogger = &log.Logger
	zerolog.SetGlobalLevel(level)

	if unknowFormat {
		ctx.Logger.Warn().
			Str("format", ctx.Config.Log.Format).
			Msg("unknown log format, defaulting to json")
	}

	// Remove date/time flags which are already present in zerolog output
	stdlog.SetFlags(stdlog.Flags() & ^(stdlog.Ldate | stdlog.Ltime | stdlog.Lmicroseconds))
	stdlog.SetOutput(ctx.Logger.With().Str("level", "stdlog").Logger())

	return nil
}
