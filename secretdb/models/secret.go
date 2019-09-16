package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Secret represents a secret entry in database.
type Secret struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Secret     string             `bson:"secret" json:"secret,omitempty"`
	Passphrase string             `bson:"passphrase,omitempty" json:"passphrase,omitempty"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at,omitempty"`
	ExpiresAt  time.Time          `bson:"expires_at" json:"expires_at,omitempty"`
}
