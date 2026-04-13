package config

import (
	"github.com/lvjp/wtf-go/internal/app/cmd/config/dump"
	"github.com/lvjp/wtf-go/internal/pkg/cmd/util"

	"github.com/spf13/cobra"
)

func NewDumpCmd() *cobra.Command {
	var ctxBuilder util.ContextBuilder

	cmd := &cobra.Command{
		Use:   "dump",
		Short: "Dump the current configuration",
		Long: `Dump command will load the current configuration which implies a validation.
And then dump the computed configuration.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := ctxBuilder.Build()
			if err != nil {
				return err
			}

			return dump.Run(ctx)
		},
	}

	flags := cmd.Flags()

	ctxBuilder.
		WithCobraCommand(cmd).
		WithVerbose(util.NewVerboseFlag(flags)).
		WithConfigPath(util.NewConfigFlag(flags))

	return cmd
}
