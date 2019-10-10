package app

import (
	"time"
)

// secretResponse Bodyrepresents a secret type response body.
type secretResponseBody struct {
	Data secretResponse `json:"data"`
}

// secretResponse represents the data part of the response body.
type secretResponse struct {
	ID        string    `json:"id,omitempty"`
	Secret    string    `json:"secret,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}
