package api

import (
	"context"
	"net/http"
)

func (c *Client) MiscVersion(ctx context.Context, opts ...OptionsFn) (*MiscVersionResponse, error) {
	var resp MiscVersionResponse

	err := c.doRequest(
		ctx,
		"MiscVersion", http.MethodGet, "misc/version",
		opts,
		deserializeJSON(&resp),
		checkStatusCode(http.StatusOK),
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) MiscHealth(ctx context.Context, opts ...OptionsFn) (*MiscHealthResponse, error) {
	var resp MiscHealthResponse

	err := c.doRequest(
		ctx,
		"MiscHealth", http.MethodGet, "misc/health",
		opts,
		deserializeJSON(&resp),
		checkStatusCode(http.StatusOK),
	)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}
