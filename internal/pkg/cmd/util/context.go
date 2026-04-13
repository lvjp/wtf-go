package util

import (
	"context"
	"fmt"
	"io"
	stdlog "log"
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

type ContextBuilder struct {
	cmd *cobra.Command

	configPath *string
	verbose    *bool
}

func (cb *ContextBuilder) WithCobraCommand(cmd *cobra.Command) *ContextBuilder {
	cb.cmd = cmd
	return cb
}

func (cb *ContextBuilder) WithVerbose(verbose *bool) *ContextBuilder {
	cb.verbose = verbose
	return cb
}

func (cb *ContextBuilder) WithConfigPath(configPath *string) *ContextBuilder {
	cb.configPath = configPath
	return cb
}

func (cb *ContextBuilder) Build() (*Context, error) {
	ret := &Context{
		Context: cb.cmd.Context(),

		Input:  cb.cmd.InOrStdin(),
		Output: cb.cmd.OutOrStdout(),
		Error:  cb.cmd.ErrOrStderr(),
	}

	if err := cb.buildConfig(ret); err != nil {
		return nil, err
	}

	if err := cb.buildLogger(ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (cb *ContextBuilder) buildConfig(ctx *Context) error {
	if cb.configPath == nil {
		ctx.Config = &config.Config{}
		ctx.Config.SetDefaults()
	} else {
		var err error
		ctx.Config, err = config.LoadFromFile(*cb.configPath)
		if err != nil {
			return err
		}
	}

	if cb.verbose != nil && *cb.verbose {
		ctx.Config.Log.Level = "debug"
	}

	return nil
}

func (cb *ContextBuilder) buildLogger(ctx *Context) error {
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
