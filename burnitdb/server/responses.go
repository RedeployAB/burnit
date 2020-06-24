package server

import (
	"time"

	"github.com/RedeployAB/burnit/burnitdb/internal/dto"
)

// responseBody represents a secret type response body.
type responseBody struct {
	Secret secret `json:"secret"`
}

// secret represents the data part of the response body.
type secret struct {
	ID        string    `json:"id,omitempty"`
	Value     string    `json:"value,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

// response creates a response from a Secret (DTO).
func response(s *dto.Secret) *responseBody {
	return &responseBody{
		Secret: secret{
			ID:        s.ID,
			Value:     s.Secret,
			CreatedAt: s.CreatedAt,
			ExpiresAt: s.ExpiresAt,
		},
	}
}
