package config

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage the wtf-go configuration",
	}

	cmd.AddCommand(NewDefaultsCmd())
	cmd.AddCommand(NewDumpCmd())

	return cmd
}
