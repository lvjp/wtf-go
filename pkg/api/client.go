package api

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	"github.com/lvjp/wtf-go/pkg/chain-of-responsibility"
)

type HTTPRequestDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	opts Options
}

func NewClient(endpoint string, optsFn ...OptionsFn) *Client {
	var opts Options
	opts.SetDefaults()

	// endpoint is not validated, but it will be used as-is in the http.NewRequest(),
	// which will return an error if the endpoint is invalid.
	opts.Endpoint = endpoint

	for _, o := range optsFn {
		o(&opts)
	}

	// Copy options to ensure immutability of the client after creation.
	return &Client{opts: opts.Copy()}
}

func (c *Client) doRequest(ctx context.Context, operationID, method, path string, caller []OptionsFn, internal ...Middleware) error {
	opts := c.opts.Copy()

	for _, o := range caller {
		o(&opts)
	}

	opts.Middlewares = append(opts.Middlewares, internal...)

	httpReq, err := http.NewRequestWithContext(ctx, method, opts.Endpoint, nil)
	if err != nil {
		return fmt.Errorf("client.%s: failed to create request: %w", operationID, err)
	}

	httpReq.URL = httpReq.URL.JoinPath(path)

	wc := &CallContext{
		Request:     httpReq,
		OperationID: operationID,
	}

	h := chain.NewChain(httpHandler(opts.HTTPClient), opts.Middlewares...)
	return h.Handle(ctx, wc)
}

type Handler = chain.Handler[*CallContext]
type HandlerFunc = chain.HandlerFunc[*CallContext]
type Middleware = chain.Middleware[*CallContext]
type MiddlewareFunc = chain.MiddlewareFunc[*CallContext]

type Options struct {
	Endpoint    string
	HTTPClient  HTTPRequestDoer
	Middlewares []Middleware
}

func (o Options) Copy() Options {
	to := o
	to.Middlewares = slices.Clone(o.Middlewares)

	return to
}

func (o *Options) SetDefaults() {
	if o.HTTPClient == nil {
		o.HTTPClient = http.DefaultClient
	}
}

// OptionsFn configures a Client. It intentionally returns no error: anything
// that can fail (loading a certificate, validating a parameter) must be
// resolved before constructing the OptionsFn, not inside it. This keeps
// NewClient infallible and pushes error handling to the call site.
type OptionsFn func(*Options)

func WithHTTPClient(client HTTPRequestDoer) OptionsFn {
	return func(o *Options) {
		o.HTTPClient = client
	}
}

func WithMiddleware(middleware Middleware) OptionsFn {
	return func(o *Options) {
		o.Middlewares = append(o.Middlewares, middleware)
	}
}

func WithMiddlewareFunc(middlewareFunc MiddlewareFunc) OptionsFn {
	return func(o *Options) {
		o.Middlewares = append(o.Middlewares, middlewareFunc)
	}
}

func WithUserAgent(userAgent string) OptionsFn {
	return WithMiddleware(MiddlewareFunc(func(ctx context.Context, cc *CallContext, next Handler) error {
		cc.Request.Header.Set("User-Agent", userAgent)
		return next.Handle(ctx, cc)
	}))
}
