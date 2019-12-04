package db

import (
	"context"
	"time"

	"github.com/RedeployAB/burnit/common/security"
	"github.com/RedeployAB/burnit/secretdb/internal/dto"
	"github.com/RedeployAB/burnit/secretdb/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository is an interface that represents
// basuc CRUD operations.
type Repository interface {
	Find(id string) (*dto.Secret, error)
	Insert(s *dto.Secret) (*dto.Secret, error)
	Delete(id string) (int64, error)
	DeleteExpired() (int64, error)
}

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
func (r *SecretRepository) Find(id string) (*dto.Secret, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var secretModel model.Secret
	err = r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&secretModel)
	if err != nil {
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
