package server

import (
	"encoding/json"
	"time"

	"github.com/RedeployAB/burnit/burnit/secret"
)

// secretResponse defines the structure of an outgoing
// secret response.
type secretResponse struct {
	ID        string     `json:"id,omitempty"`
	Value     string     `json:"value,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

// JSON marshals the secret to JSON.
func (s *secretResponse) JSON() []byte {
	b, _ := json.Marshal(&s)
	return b
}

// newSecretResponse creates a secret response from a Secret (DTO).
func newSecretResponse(s *secret.Secret) *secretResponse {
	var createdAt, expiresAt *time.Time
	if !s.CreatedAt.IsZero() {
		createdAt = &s.CreatedAt
	} else {
		createdAt = nil
	}
	if !s.ExpiresAt.IsZero() {
		expiresAt = &s.ExpiresAt
	} else {
		expiresAt = nil
	}

	return &secretResponse{
		ID:        s.ID,
		Value:     s.Value,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}
