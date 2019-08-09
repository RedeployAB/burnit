package api

import "time"

// Response represents a standard response.
type Response struct {
	Data string `json:"data"`
}

// ErrorResponse represents an error response.
type ErrorResponse struct {
	Error string `json:"error"`
}

// SecretResponse represents a secret type response.
type SecretResponse struct {
	ID        string    `json:"id,omitempty"`
	Secret    string    `json:"secret,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
