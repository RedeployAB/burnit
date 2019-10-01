package api

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// secret represents the data part of the response body.
type secret struct {
	ID        primitive.ObjectID `json:"id,omitempty"`
	Secret    string             `json:"secret,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty"`
	ExpiresAt time.Time          `json:"expiresAt,omitempty"`
}
