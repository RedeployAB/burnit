package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Secret represents a secret entry in database.
type Secret struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Secret     string             `bson:"secret"`
	Passphrase string             `bson:"passphrase,omitempty"`
	CreatedAt  time.Time          `bson:"createdAt"`
	ExpiresAt  time.Time          `bson:"expiresAt"`
}
