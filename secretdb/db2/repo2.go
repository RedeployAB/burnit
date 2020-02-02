package db2

import (
	"context"
	"time"

	"github.com/RedeployAB/burnit/common/security"
	"github.com/RedeployAB/burnit/secretdb/internal/model"

	"github.com/RedeployAB/burnit/secretdb/internal/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository represents a repository containing methods
// to interact with a database and collection.
type Repository interface {
	Find(id string) (*dto.Secret, error)
	Delete(id string) (int64, error)
	/* 	Insert(m *interface{}) (interface{}, error)
	   	Delete(id string) (int64, error) */
}

// SecretRepository handles interactions with the database
// and collection.
type SecretRepository struct {
	collection *mongo.Collection
	passphrase string
}

// Find queries the collection for an entry by ID.
func (r *SecretRepository) Find(id string) (*dto.Secret, error) {
	// Create an ObjectID from incoming string to
	// use in search query.
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// Revisit this scenario.
		return nil, nil
	}

	// Make timeout configurable.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var secretModel model.Secret
	bsonQ := bson.D{{Key: "_id", Value: oid}}
	if err = r.collection.FindOne(ctx, bsonQ).Decode(&secretModel); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return &dto.Secret{}, err
	}

	return &dto.Secret{
		ID:         oid.Hex(),
		Secret:     decrypt(secretModel.Secret, r.passphrase),
		Passphrase: secretModel.Passphrase,
		CreatedAt:  secretModel.CreatedAt,
		ExpiresAt:  secretModel.ExpiresAt,
	}, nil
}

// encrypt encrypts the field Secret.
func encrypt(plaintext, passphrase string) string {
	encrypted, err := security.Encrypt([]byte(plaintext), passphrase)
	if err != nil {
		panic(err)
	}
	return string(encrypted)
}

// decrypt decrypts the field Secret.
func decrypt(ciphertext, passphrase string) string {
	decrypted, err := security.Decrypt([]byte(ciphertext), passphrase)
	if err != nil {
		panic(err)
	}
	return string(decrypted)
}
