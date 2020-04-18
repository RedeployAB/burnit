package models

import (
	"time"
)

// Secret represents a secret entry in database.
type Secret struct {
	ID         string    `bson:"_id,omitempty"`
	Secret     string    `bson:"secret"`
	Passphrase string    `bson:"passphrase,omitempty"`
	CreatedAt  time.Time `bson:"createdAt"`
	ExpiresAt  time.Time `bson:"expiresAt"`
}
