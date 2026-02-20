package cmd

import (
	"fmt"
	"os"

	"git.sr.ht/~lvjp/wtf-go/internal/pkg/cmd/util"
	"git.sr.ht/~lvjp/wtf-go/pkg/buildinfo"

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
	flags := rootCmd.PersistentFlags()

	// ANCHOR: default_config_path
	config := flags.String("config", "/etc/wtf-go/config.yaml", "Path to the configuration file")
	// ANCHOR_END: default_config_path

	verbose := flags.Bool("verbose", false, "Enable verbose logging (debug level)")

	factory := util.NewFactory(config, verbose)
	rootCmd.AddCommand(NewCharmBracelet(factory))
	rootCmd.AddCommand(NewHealthCheckCmd(factory))
	rootCmd.AddCommand(NewServerCmd(factory))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
