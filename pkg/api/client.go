package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func NewBadStatusCodeError(body []byte, resp *http.Response, expectedStatusCode int, operationID string) error {
	return &HTTPRequestError{
		Body:         body,
		HTTPResponse: resp,
		Message: fmt.Sprintf(
			"client.%s: unexpected status code: %d, expected: %d",
			operationID,
			resp.StatusCode,
			expectedStatusCode,
		),
	}
}

func NewHTTPError(body []byte, resp *http.Response, message string, wrapped error) error {
	return &HTTPRequestError{
		Body:         body,
		HTTPResponse: resp,
		Message:      fmt.Sprintf("%s: %v", message, wrapped),
		Wrapped:      wrapped,
	}
}

type HTTPRequestError struct {
	Body         []byte
	HTTPResponse *http.Response
	Message      string
	Wrapped      error
}

func (e *HTTPRequestError) Error() string {
	return e.Message
}

func (e *HTTPRequestError) Unwrap() error {
	return e.Wrapped
}

type ClientInterface interface {
	MiscVersion(context.Context) (*MiscVersionResponse, error)
	MiscHealth(context.Context) (*MiscHealthResponse, error)
}

type RequestEditorFn func(ctx context.Context, req *http.Request) error

type HttpRequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	Endpoint       *url.URL
	Client         HttpRequestDoer
	RequestEditors []RequestEditorFn
}

type ClientOption func(*Client) error

func WithHTTPClient(client HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = client
		return nil
	}
}

func WithRequestEditor(editor RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, editor)
		return nil
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return WithRequestEditor(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("User-Agent", userAgent)
		return nil
	})
}

func NewClient(endpoint string, opts ...ClientOption) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("client: failed to parse endpoint: %w", err)
	}

	return NewClientWithURL(u, opts...)
}

func NewClientWithURL(endpoint *url.URL, opts ...ClientOption) (*Client, error) {
	if endpoint == nil {
		return nil, fmt.Errorf("client: endpoint cannot be nil")
	}

	c := Client{
		Endpoint: endpoint,
	}

	for _, o := range opts {
		if err := o(&c); err != nil {
			return nil, err
		}
	}

	// If no client is provided, use the default HTTP client.
	if c.Client == nil {
		c.Client = http.DefaultClient
	}

	return &c, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}

	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) MiscVersion(ctx context.Context, reqEditors ...RequestEditorFn) (*MiscVersionResponse, error) {
	return doJsonRequest[MiscVersionResponse](ctx, c, "MiscVersion", http.MethodGet, "misc/version", http.StatusOK, reqEditors...)
}

func (c *Client) MiscHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*MiscHealthResponse, error) {
	return doJsonRequest[MiscHealthResponse](ctx, c, "MiscHealth", http.MethodGet, "misc/health", http.StatusOK, reqEditors...)
}

func doJsonRequest[T any](ctx context.Context, c *Client, operationID, method, path string, statusCode int, reqEditors ...RequestEditorFn) (*T, error) {
	httpResp, body, err := doRequest(ctx, c, operationID, method, path, statusCode, reqEditors...) //nolint:bodyclose
	if err != nil {
		return nil, err
	}

	var resp T
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, NewHTTPError(body, httpResp, fmt.Sprintf("client.%s: failed to unmarshal response body", operationID), err)
	}

	return &resp, nil
}

func doRequest(ctx context.Context, c *Client, operationID, method, path string, statusCode int, reqEditors ...RequestEditorFn) (*http.Response, []byte, error) {
	u := c.Endpoint.JoinPath(path)

	httpReq, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("client.%s: failed to create request: %w", operationID, err)
	}

	if applyErr := c.applyEditors(ctx, httpReq, reqEditors); applyErr != nil {
		return nil, nil, fmt.Errorf("client.%s: failed to apply request editors: %w", operationID, applyErr)
	}

	httpResp, err := c.Client.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("client.%s: request failed: %w", operationID, err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return httpResp, nil, fmt.Errorf("client.%s: failed to read response body: %w", operationID, err)
	}

	if httpResp.StatusCode != statusCode {
		return httpResp, body, NewBadStatusCodeError(body, httpResp, statusCode, operationID)
	}

	return httpResp, body, nil
}
