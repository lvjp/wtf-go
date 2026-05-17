package api

import (
	"fmt"
	"net/http"
)

func NewHTTPError(body []byte, resp *http.Response, message string, wrapped error) error {
	if wrapped != nil {
		message = fmt.Sprintf("%s: %v", message, wrapped)
	}

	return &HTTPRequestError{
		Body:         body,
		HTTPResponse: resp,
		Message:      message,
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
