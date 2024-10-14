package api

import (
	"context"
	"time"
)

// Secret represents a secret.
type Secret struct {
	ID         string `json:"id,omitempty"`
	Value      string `json:"value,omitempty"`
	Passphrase string `json:"passphrase,omitempty"`
	Path       string `json:"path,omitempty"`
	TTL        string `json:"ttl,omitempty"`
	ExpiresAt  *Time  `json:"expiresAt,omitempty"`
}

// CreateSecretRequest represents a request to create a secret.
type CreateSecretRequest struct {
	Value      string `json:"value,omitempty"`
	Passphrase string `json:"passphrase,omitempty"`
	TTL        string `json:"ttl,omitempty"`
	ExpiresAt  *Time  `json:"expiresAt,omitempty"`
}

// Valid validates the CreateSecretRequest.
func (r CreateSecretRequest) Valid(ctx context.Context) map[string]string {
	errs := make(map[string]string)
	if len(r.Value) == 0 {
		errs["value"] = "value is required"
	}
	if len(r.TTL) > 0 {
		_, err := time.ParseDuration(r.TTL)
		if err != nil {
			errs["ttl"] = "ttl is invalid, expected format is 1h30m"
		}
	}
	return errs
}
