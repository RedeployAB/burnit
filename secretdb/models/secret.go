package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Secret represents a secret entry in database.
type Secret struct {
	ID         bson.ObjectId `bson:"_id,omitempty"`
	Secret     string
	Passphrase string
	CreatedAt  time.Time `bson:"created_at"`
	ExpiresAt  time.Time `bson:"expires_at"`
}
