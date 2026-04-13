package cmd

import (
	"github.com/lvjp/wtf-go/internal/app/cmd/healthcheck"
	"github.com/lvjp/wtf-go/internal/pkg/cmd/util"

	"github.com/spf13/cobra"
)

func NewHealthCheckCmd() *cobra.Command {
	var ctxBuilder util.ContextBuilder

	cmd := &cobra.Command{
		Use:   "healthcheck",
		Short: "Check the health of the server",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := ctxBuilder.Build()
			if err != nil {
				return err
			}

			return healthcheck.Run(ctx)
		},
	}

	flags := cmd.Flags()

	ctxBuilder.
		WithCobraCommand(cmd).
		WithVerbose(util.NewVerboseFlag(flags)).
		WithConfigPath(util.NewConfigFlag(flags))

	return cmd
}
