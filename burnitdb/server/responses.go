package server

import (
	"time"

	"github.com/RedeployAB/burnit/burnitdb/secret"
)

// secretResponse defines the structure of an outgoing
// secret response.
type secretResponse struct {
	Secret secretBody `json:"secret"`
}

// secretBody defines the structure of an outgoing
// secret response body.
type secretBody struct {
	ID        string    `json:"id,omitempty"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

// newSecretResponse creates a secret response from a Secret (DTO).
func newSecretResponse(s *secret.Secret) *secretResponse {
	return &secretResponse{
		Secret: secretBody{
			ID:        s.ID,
			Value:     s.Secret,
			CreatedAt: s.CreatedAt,
			ExpiresAt: s.ExpiresAt,
		},
	}
}
