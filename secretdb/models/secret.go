package models

import "time"

// Secret represents a secret entry in database.
type Secret struct {
	Secret     string    `bson:"secret"`
	Passphrase string    `bson:"passphrase,omitempty"`
	CreatedAt  time.Time `bson:"created_at"`
	ExpiresAt  time.Time `bson:"expires_at"`
}
