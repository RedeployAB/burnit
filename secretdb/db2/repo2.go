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
	Insert(s *dto.Secret) (*dto.Secret, error)
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

// Insert handles inserts into the database.
func (r *SecretRepository) Insert(s *dto.Secret) (*dto.Secret, error) {
	if s.TTL == 0 {
		s.TTL = 10080
	}

	secretModel := &model.Secret{
		Secret:    encrypt(s.Secret, r.passphrase),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.TTL)),
	}

	if len(s.Passphrase) > 0 {
		secretModel.Passphrase = s.Passphrase
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.collection.InsertOne(ctx, secretModel)
	if err != nil {
		return &dto.Secret{}, err
	}
	oid := res.InsertedID.(primitive.ObjectID)

	return &dto.Secret{
		ID:         oid.Hex(),
		Passphrase: secretModel.Passphrase,
		CreatedAt:  secretModel.CreatedAt,
		ExpiresAt:  secretModel.ExpiresAt,
	}, nil
}

// Delete removes an entry from the collection by ID.
func (r *SecretRepository) Delete(id string) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// NewSecretRepository creates and returns a SecretRepository
// object.
func NewSecretRepository(c *mongo.Client, passphrase string) *SecretRepository {
	return &SecretRepository{
		collection: c.Database("secretdb").Collection("secrets"),
		passphrase: passphrase,
	}
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
