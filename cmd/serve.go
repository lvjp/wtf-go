package cmd

import (
	"github.com/lvjp/wtf-go/internal/app/cmd/serve"
	"github.com/lvjp/wtf-go/internal/pkg/cmd/util"

	"github.com/spf13/cobra"
)

func NewServerCmd(factory *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Serve the wtf-go API",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := factory.NewContext(cmd)
			return serve.Run(ctx)
		},
	}
}
