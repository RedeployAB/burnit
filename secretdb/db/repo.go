package db

import (
	"context"
	"time"

	"github.com/RedeployAB/redeploy-secrets/common/security"
	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
	"github.com/RedeployAB/redeploy-secrets/secretdb/models"
	"github.com/RedeployAB/redeploy-secrets/secretdb/secret"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SecretRepository handles actions with the database and
// collection.
type SecretRepository struct {
	DB     *DB
	Config config.Server
}

// Find queries the collection for an entry by ID.
func (r *SecretRepository) Find(id string) (*secret.Secret, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return &secret.Secret{}, err
	}

	collection := r.DB.Database("secretdb").Collection("secrets")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var sm models.Secret
	err = collection.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&sm)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return &secret.Secret{}, err
	}

	return &secret.Secret{
		ID:         oid.Hex(),
		Secret:     decrypt(sm.Secret, r.Config.Passphrase),
		Passphrase: sm.Passphrase,
		CreatedAt:  sm.CreatedAt,
		ExpiresAt:  sm.ExpiresAt,
	}, nil
}

// Insert handles inserts into the database.
func (r *SecretRepository) Insert(s *secret.Secret) (*secret.Secret, error) {

	collection := r.DB.Database("secretdb").Collection("secrets")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if s.TTL == 0 {
		s.TTL = 10080
	}

	sm := &models.Secret{
		Secret:    encrypt(s.Secret, r.Config.Passphrase),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(s.TTL)),
	}

	if len(s.Passphrase) > 0 {
		sm.Passphrase = s.Passphrase
	}

	res, err := collection.InsertOne(ctx, sm)
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
	collection := r.DB.Database("secretdb").Collection("secrets")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return -1, nil
	}

	res, err := collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// DeleteExpired deletes all entries that has expiresAt
// less than current time (time of invocation).
func (r *SecretRepository) DeleteExpired() (int64, error) {
	collection := r.DB.Database("secretdb").Collection("secrets")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := collection.DeleteMany(ctx, bson.D{{Key: "expiresAt", Value: bson.D{{Key: "$lt", Value: time.Now()}}}})
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
