package cmd

import (
	"fmt"
	"os"

	configpkg "github.com/lvjp/wtf-go/cmd/config"
	"github.com/lvjp/wtf-go/pkg/buildinfo"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "wtf-go",
	Short:         "What The Fuck in go",
	Long:          `wtf-go is just something in go. Just coding`,
	SilenceErrors: true,
	SilenceUsage:  true,
	Version:       buildinfo.Get().String(),
}

func init() {
	rootCmd.AddCommand(configpkg.New())

	rootCmd.AddCommand(NewHealthCheckCmd())
	rootCmd.AddCommand(NewServerCmd())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
