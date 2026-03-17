package api

import "time"

type MiscVersionResponse struct {
	Go       string    `json:"go"`
	Modified bool      `json:"modified"`
	Platform string    `json:"platform"`
	Revision string    `json:"revision,omitempty"`
	Time     time.Time `json:"time,omitzero"`
}

type MiscHealthResponse struct {
	Status string `json:"status"`
}
