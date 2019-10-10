package db

import (
	"context"
	"time"

	"github.com/RedeployAB/burnit/common/security"
	"github.com/RedeployAB/burnit/secretdb/models"
	"github.com/RedeployAB/burnit/secretdb/secret"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository is an interface to handle different
// repository types.
/* type Repository interface {
	Find() (*secret.Secret, error)
	Insert() (*secret.Secret, error)
	Delete() (int64, error)
} */

// SecretRepository handles actions with the database and
// collection.
type SecretRepository struct {
	collection *mongo.Collection
	passphrase string
}

// NewSecretRepository creates a new SecretRepository.
func NewSecretRepository(c *Connection, passphrase string) *SecretRepository {
	return &SecretRepository{
		collection: c.Database("secretdb").Collection("secrets"),
		passphrase: passphrase,
	}
}

// Find queries the collection for an entry by ID.
func (r *SecretRepository) Find(id string) (*secret.Secret, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &secret.Secret{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sm models.Secret
	err = r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&sm)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return &secret.Secret{}, err
	}

	return &secret.Secret{
		ID:         oid.Hex(),
		Secret:     decrypt(sm.Secret, r.passphrase),
		Passphrase: sm.Passphrase,
		CreatedAt:  sm.CreatedAt,
		ExpiresAt:  sm.ExpiresAt,
	}, nil
}

// Insert handles inserts into the database.
func (r *SecretRepository) Insert(s *secret.Secret) (*secret.Secret, error) {
	if s.TTL == 0 {
		s.TTL = 10080
	}

	sm := &models.Secret{
		Secret:    encrypt(s.Secret, r.passphrase),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.TTL)),
	}

	if len(s.Passphrase) > 0 {
		sm.Passphrase = s.Passphrase
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.collection.InsertOne(ctx, sm)
	if err != nil {
		return &secret.Secret{}, err
	}
	oid := res.InsertedID.(primitive.ObjectID)

	return &secret.Secret{
		ID:         oid.Hex(),
		Passphrase: sm.Passphrase,
		CreatedAt:  sm.CreatedAt,
		ExpiresAt:  sm.ExpiresAt,
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
