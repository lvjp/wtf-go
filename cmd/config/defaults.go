package config

import (
	"github.com/lvjp/wtf-go/internal/app/cmd/config/defaults"
	"github.com/lvjp/wtf-go/internal/pkg/cmd/util"

	"github.com/spf13/cobra"
)

func NewDefaultsCmd() *cobra.Command {
	var ctxBuilder util.ContextBuilder

	cmd := &cobra.Command{
		Use:   "defaults",
		Short: "Print the default configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx, err := ctxBuilder.Build()
			if err != nil {
				return err
			}

			return defaults.Run(ctx)
		},
	}

	flags := cmd.Flags()

	ctxBuilder.
		WithCobraCommand(cmd).
		WithVerbose(util.NewVerboseFlag(flags))

	return cmd
}
