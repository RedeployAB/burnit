package api

import (
	"context"
	"time"
)

// Secret represents a secret.
type Secret struct {
	ID        string        `json:"id,omitempty"`
	Value     string        `json:"value,omitempty"`
	TTL       time.Duration `json:"ttl,omitempty"`
	ExpiresAt *time.Time    `json:"expiresAt,omitempty"`
}

// CreateSecretRequest represents a request to create a secret.
type CreateSecretRequest struct {
	Value      string        `json:"value,omitempty"`
	Passphrase string        `json:"passphrase,omitempty"`
	TTL        time.Duration `json:"ttl,omitempty"`
}

// Valid validates the CreateSecretRequest.
func (r CreateSecretRequest) Valid(ctx context.Context) map[string]string {
	errs := make(map[string]string)
	if len(r.Value) == 0 {
		errs["value"] = "value is required"
	}
	return errs
}
