package server

import (
	"time"

	"github.com/RedeployAB/burnit/burnitdb/secret"
)

// secretResponse defines the structure of an outgoing
// secret response.
type secretResponse struct {
	ID        string    `json:"id,omitempty"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

// newSecretResponse creates a secret response from a Secret (DTO).
func newSecretResponse(s *secret.Secret) *secretResponse {
	return &secretResponse{
		ID:        s.ID,
		Value:     s.Value,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}
