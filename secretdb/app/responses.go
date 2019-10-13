package app

import (
	"time"
)

// secretResponse Bodyrepresents a secret type response body.
type secretResponse struct {
	Data secretData `json:"data"`
}

// secretData represents the data part of the response body.
type secretData struct {
	ID        string    `json:"id,omitempty"`
	Secret    string    `json:"secret,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}
