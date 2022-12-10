package db

import (
	"time"
)

// Secret represents a secret entry in database.
type Secret struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Value     string    `json:"value" bson:"value"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
}
