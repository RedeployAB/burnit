package db

import "time"

// Secret represents a secret entry in the database.
type Secret struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Value     string    `json:"value" bson:"value"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
}
