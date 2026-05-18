package util

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestNewConfigFlag_coverage(t *testing.T) {
	NewConfigFlag(pflag.NewFlagSet(t.Name(), pflag.ContinueOnError))
}

func TestNewVerboseFlag_coverage(t *testing.T) {
	NewVerboseFlag(pflag.NewFlagSet(t.Name(), pflag.ContinueOnError))
}
