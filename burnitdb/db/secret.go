package db

import "time"

// Secret represents a secret entry in database.
type Secret struct {
	ID         string    `json:"id,omitempty" bson:"_id,omitempty"`
	Secret     string    `json:"secret" bson:"secret"`
	Passphrase string    `json:"passphrase,omitempty" bson:"passphrase,omitempty"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	ExpiresAt  time.Time `json:"expiresAt" bson:"expiresAt"`
}
