package cmd

import (
	"git.sr.ht/~lvjp/wtf-go/internal/app/cmd/charmbracelet"
	"git.sr.ht/~lvjp/wtf-go/internal/pkg/cmd/util"

	"github.com/spf13/cobra"
)

func NewCharmBracelet(factory *util.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "charmbracelet",
		Short: "A simple TUI example using the Charmbracelet framework",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := factory.NewContext(cmd)
			return charmbracelet.Run(ctx)
		},
	}
}
