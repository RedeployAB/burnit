package db

import (
	"context"
	"fmt"
	"time"

	"github.com/RedeployAB/burnit/common/security"
	"github.com/RedeployAB/burnit/secretdb/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository represents a repository containing methods
// to interact with a database and collection.
type Repository interface {
	Find(id string) (*models.Secret, error)
	Insert(s *models.Secret) (*models.Secret, error)
	Delete(id string) (int64, error)
}

// SecretRepository handles interactions with the database
// and collection.
type SecretRepository struct {
	collection *mongo.Collection
	passphrase string
}

// NewSecretRepository creates and returns a SecretRepository
// object.
func NewSecretRepository(c *Client, passphrase string) *SecretRepository {
	return &SecretRepository{
		collection: c.Database("secretdb").Collection("secrets"),
		passphrase: passphrase,
	}
}

// Find queries the collection for an entry by ID.
func (r *SecretRepository) Find(id string) (*models.Secret, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, &QueryError{err.Error(), -1}
	}

	var s models.Secret
	bsonQ := bson.D{{Key: "_id", Value: oid}}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = r.collection.FindOne(ctx, bsonQ).Decode(&s); err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, &QueryError{err.Error(), 0}
		}
		return &models.Secret{}, err
	}

	s.Secret = decrypt(s.Secret, r.passphrase)

	return &s, nil
}

// Insert handles inserts into the database.
func (r *SecretRepository) Insert(s *models.Secret) (*models.Secret, error) {
	s.Secret = encrypt(s.Secret, r.passphrase)
	if len(s.Passphrase) > 0 {
		s.Passphrase = hash(s.Passphrase)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.collection.InsertOne(ctx, s)
	if err != nil {
		return &models.Secret{}, err
	}
	oid := res.InsertedID.(primitive.ObjectID)

	return &models.Secret{
		ID:         oid,
		Passphrase: s.Passphrase,
		CreatedAt:  s.CreatedAt,
		ExpiresAt:  s.ExpiresAt,
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

// DeleteExpired deletes all entries that has expiresAt
// less than current time (time of invocation).
func (r *SecretRepository) DeleteExpired() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.collection.DeleteMany(ctx, bson.D{{Key: "expiresAt", Value: bson.D{{Key: "$lt", Value: time.Now()}}}})
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

// QueryError wraps database related errors with
// Error, containing error message, and code
// that holds error code.
type QueryError struct {
	Message string
	Code    int
}

func (e *QueryError) Error() string {
	return fmt.Sprintf("query: %s, code: %d", e.Message, e.Code)
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

// hash hashes the incmong string with bcrypt.
func hash(s string) string {
	return security.Hash(s)
}
