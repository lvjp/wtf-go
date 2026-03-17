package healthcheck

import (
	"fmt"

	"github.com/lvjp/wtf-go/internal/pkg/cmd/util"
	"github.com/lvjp/wtf-go/pkg/api"
	"github.com/lvjp/wtf-go/pkg/buildinfo"
)

func Run(ctx *util.Context) error {
	endpoint := fmt.Sprintf("http://%s/api/v0", *ctx.Config.Server.ListenAddress)
	fmt.Fprintln(ctx.Output, "Endpoint:", endpoint)

	c, err := api.NewClient(endpoint, api.WithUserAgent("wtf-go/"+buildinfo.Get().Revision))
	if err != nil {
		return fmt.Errorf("client creation: %v", err)
	}

	resp, err := c.MiscHealth(ctx)
	if err != nil {
		return fmt.Errorf("health check failed: %v", err)
	}

	fmt.Fprintln(ctx.Output, "Status:", resp.Status)

	return nil
}
