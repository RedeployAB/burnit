package api

import "time"

// ResponseBody represents a standard response.
type ResponseBody struct {
	Data string `json:"data"`
}

// ErrorResponseBody represents an error response.
type ErrorResponseBody struct {
	Error string `json:"error"`
}

// SecretResponseBody represents a secret type response.
type SecretResponseBody struct {
	Data secretData `json:"data"`
}

type secretData struct {
	ID        string    `json:"id,omitempty"`
	Secret    string    `json:"secret,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
