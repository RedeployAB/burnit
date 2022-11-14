package server

import (
	"encoding/json"
	"time"

	"github.com/RedeployAB/burnit/burnit/secrets"
)

// secret defines the structure of an outgoing
// secret response.
type secret struct {
	ID        string     `json:"id,omitempty"`
	Value     string     `json:"value,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

// JSON marshals the secret to JSON.
func (s *secret) JSON() []byte {
	b, _ := json.Marshal(&s)
	return b
}

// newSecret creates a secret response from a Secret (DTO).
func newSecret(s *secrets.Secret) *secret {
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

	return &secret{
		ID:        s.ID,
		Value:     s.Value,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}
