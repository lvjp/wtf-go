package util

import "github.com/spf13/pflag"

func NewConfigFlag(flags *pflag.FlagSet) *string {
	// ANCHOR: default_config_path
	return flags.String("config", "", "Path to the configuration file")
	// ANCHOR_END: default_config_path
}

func NewVerboseFlag(flags *pflag.FlagSet) *bool {
	return flags.Bool("verbose", false, "Enable verbose logging (debug level)")
}
