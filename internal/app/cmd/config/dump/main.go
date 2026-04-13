package dump

import (
	"fmt"

	"github.com/lvjp/wtf-go/internal/pkg/cmd/util"

	"github.com/goccy/go-yaml"
)

func Run(ctx *util.Context) error {
	enc := yaml.NewEncoder(ctx.Output)

	if err := enc.Encode(ctx.Config); err != nil {
		return fmt.Errorf("cannot encode configuration: %v", err)
	}

	return nil
}
