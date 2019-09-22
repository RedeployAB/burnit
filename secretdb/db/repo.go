package db

import (
	"context"
	"time"

	"github.com/RedeployAB/redeploy-secrets/common/security"
	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
	"github.com/RedeployAB/redeploy-secrets/secretdb/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Find queries the collection for an entry by ID.
func Find(id string, client *mongo.Client) (models.Secret, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Secret{}, err
	}

	collection := client.Database("secretdb").Collection("secrets")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	var s models.Secret
	err = collection.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&s)
	if err != nil {
		return models.Secret{}, err
	}
	s.Secret = decrypt(s.Secret, config.Config.Passphrase)

	return s, nil
}

// Insert handles inserts into the database.
func Insert(s models.Secret, client *mongo.Client) (models.Secret, error) {
	collection := client.Database("secretdb").Collection("secrets")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	sm := &models.Secret{
		Secret:    encrypt(s.Secret, config.Config.Passphrase),
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	}

	if len(s.Passphrase) > 0 {
		sm.Passphrase = security.Hash(s.Passphrase)
	}

	res, err := collection.InsertOne(ctx, sm)
	if err != nil {
		return models.Secret{}, err
	}
	oid := res.InsertedID.(primitive.ObjectID)

	return models.Secret{ID: oid, CreatedAt: sm.CreatedAt, ExpiresAt: sm.ExpiresAt}, nil
}

// Delete removes an entry from the collection by ID.
func Delete(id string, client *mongo.Client) (int64, error) {
	collection := client.Database("secretdb").Collection("secrets")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}

	res, err := collection.DeleteOne(ctx, bson.D{{"_id", oid}})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

// Wrapper around internal.Encrypt to ease usage.
func encrypt(plaintext, passphrase string) string {
	encrypted, err := security.Encrypt([]byte(plaintext), passphrase)
	if err != nil {
		panic(err)
	}
	return string(encrypted)
}

// Wrapper around internal.Decrupt to ease usage.
func decrypt(ciphertext, passphrase string) string {
	decrypted, err := security.Decrypt([]byte(ciphertext), passphrase)
	if err != nil {
		panic(err)
	}
	return string(decrypted)
}
