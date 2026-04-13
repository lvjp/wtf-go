package defaults

import (
	"fmt"

	"github.com/lvjp/wtf-go/internal/app/config"
	"github.com/lvjp/wtf-go/internal/pkg/cmd/util"

	"github.com/goccy/go-yaml"
)

func Run(ctx *util.Context) error {
	config, path, err := config.New(config.WithDefaults())
	if err != nil {
		return fmt.Errorf("config builder: %v", err)
	}

	if path != "" {
		return fmt.Errorf("unexpected config file used: %s", path)
	}

	enc := yaml.NewEncoder(ctx.Output)

	if err := enc.Encode(config); err != nil {
		return fmt.Errorf("cannot encode configuration: %v", err)
	}

	return nil
}
