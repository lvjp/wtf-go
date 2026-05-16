package api

import "time"

type AuthTokenCreateRequest struct {
	Subject string `json:"subject"`
}

type AuthTokenResponse struct {
	ID       string    `json:"id"`
	NotAfter time.Time `json:"not_after"`
}

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
