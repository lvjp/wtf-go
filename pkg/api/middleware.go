package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CallContext struct {
	OperationID string

	Request  *http.Request
	Response *http.Response
	Body     []byte
}

func httpHandler(httpClient HTTPRequestDoer) Handler {
	return HandlerFunc(func(_ context.Context, cc *CallContext) error {
		resp, err := httpClient.Do(cc.Request)
		if err != nil {
			return fmt.Errorf("client.%s: request failed: %w", cc.OperationID, err)
		}
		defer resp.Body.Close()

		cc.Response = resp
		cc.Body, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("client.%s: failed to read response body: %w", cc.OperationID, err)
		}

		return nil
	})
}

func checkStatusCode(expectedStatusCode int) Middleware {
	return MiddlewareFunc(func(ctx context.Context, cc *CallContext, next Handler) error {
		if err := next.Handle(ctx, cc); err != nil {
			return err
		}

		if cc.Response.StatusCode != expectedStatusCode {
			msg := fmt.Sprintf(
				"client.%s: unexpected status code: %d, expected: %d",
				cc.OperationID,
				cc.Response.StatusCode,
				expectedStatusCode,
			)
			return NewHTTPError(cc.Body, cc.Response, msg, nil)
		}

		return nil
	})
}

func deserializeJSON[T any](dest *T) Middleware {
	return MiddlewareFunc(func(ctx context.Context, cc *CallContext, next Handler) error {
		if err := next.Handle(ctx, cc); err != nil {
			return err
		}

		if err := json.Unmarshal(cc.Body, dest); err != nil {
			return NewHTTPError(cc.Body, cc.Response, fmt.Sprintf("client.%s: failed to unmarshal response body", cc.OperationID), err)
		}

		return nil
	})
}
