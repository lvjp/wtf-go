package api

import (
	"context"
	"net/http"
)

func (c *Client) AuthTokenCreate(ctx context.Context, subject string, opts ...OptionsFn) (*AuthTokenResponse, error) {
	var resp AuthTokenResponse

	err := c.doRequest(
		ctx,
		"AuthTokenCreate", http.MethodPost, "auth/token",
		opts,
		deserializeJSON(&resp),
		checkStatusCode(http.StatusCreated),
		serializeJSON(AuthTokenCreateRequest{Subject: subject}),
	)
	if err != nil {
		return nil, err
	}

	c.opts.AuthToken = resp.ID
	return &resp, nil
}

func (c *Client) AuthTokenRevoke(ctx context.Context, opts ...OptionsFn) error {
	err := c.doRequest(
		ctx,
		"AuthTokenRevoke", http.MethodDelete, "auth/token",
		opts,
		checkStatusCode(http.StatusNoContent),
	)
	if err != nil {
		return err
	}

	c.opts.AuthToken = ""
	return nil
}

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
