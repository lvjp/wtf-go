package cmd

import (
	"github.com/spf13/cobra"

	"github.com/lvjp/wtf-go/internal/app/cmd/serve"
	"github.com/lvjp/wtf-go/internal/pkg/cmd/util"
)

func NewServerCmd() *cobra.Command {
	var ctxBuilder util.ContextBuilder

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve the wtf-go API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := ctxBuilder.Build()
			if err != nil {
				return err
			}

			return serve.Run(ctx, nil)
		},
	}

	flags := cmd.Flags()

	ctxBuilder.
		WithCobraCommand(cmd).
		WithVerbose(util.NewVerboseFlag(flags)).
		WithConfigPath(util.NewConfigFlag(flags))

	return cmd
}
