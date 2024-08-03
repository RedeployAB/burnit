package api

import "time"

// Secret represents a secret.
type Secret struct {
	ID    string        `json:"id,omitempty"`
	Value string        `json:"value,omitempty"`
	TTL   time.Duration `json:"ttl,omitempty"`
}

// CreateSecretRequest represents a request to create a secret.
type CreateSecretRequest struct {
	Value string        `json:"value,omitempty"`
	TTL   time.Duration `json:"ttl,omitempty"`
}

// Valid validates the CreateSecretRequest.
func (r CreateSecretRequest) Valid() map[string]string {
	errs := make(map[string]string)
	if r.Value == "" {
		errs["value"] = "value is required"
	}
	return errs
}
