package app

import (
	"time"

	"github.com/RedeployAB/burnit/burnitdb/internal/dto"
)

// responseBody represents a secret type response body.
type responseBody struct {
	Data data `json:"data"`
}

// data represents the data part of the response body.
type data struct {
	ID        string    `json:"id,omitempty"`
	Secret    string    `json:"secret,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	ExpiresAt time.Time `json:"expiresAt,omitempty"`
}

// response creates a response from a Secret (DTO).
func response(s *dto.Secret) *responseBody {
	return &responseBody{
		Data: data{
			ID:        s.ID,
			Secret:    s.Secret,
			CreatedAt: s.CreatedAt,
			ExpiresAt: s.ExpiresAt,
		},
	}
}
