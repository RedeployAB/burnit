package db

import (
	"context"
	"time"

	"github.com/RedeployAB/redeploy-secrets/secretdb/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Find queries the collection for an entry by ID.
func Find(id string, client *mongo.Client) (Secret, error) {

}

// Insert handles inserts into the database.
func Insert(s Secret, client *mongo.Client) (Secret, error) {

	sm := &models.Secret{
		Secret:    s.Secret,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	}

	collection := client.Database("secretdb").Collection("secrets")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	res, err := collection.InsertOne(ctx, sm)
	if err != nil {
		return Secret{}, err
	}
	oid := res.InsertedID.(primitive.ObjectID).Hex()

	return Secret{ID: oid, CreatedAt: sm.CreatedAt, ExpiresAt: sm.ExpiresAt}, nil
}

// Delete removes an entry from the collection by ID.
func Delete(id string, client *mongo.Client) error {

}

// Secret represents a secret to be inserted into the
// database collection.
type Secret struct {
	ID         string    `json:"id,omitempty"`
	Secret     string    `json:"secret,omitempty"`
	Passphrase string    `json:"passphrase,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
}
