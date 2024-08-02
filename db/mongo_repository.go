package db

import (
	"context"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/db/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// defaultMongoSecretRepositoryDatabase is the default database for the secret repository.
	defaultMongoSecretRepositoryDatabase = "burnit"
	// defaultMongoSecretRepositoryCollection is the default collection for the secret repository.
	defaultMongoSecretRepositoryCollection = "secrets"
	// defaultMongoSecretRepositoryTimeout is the default timeout for the secret repository.
	defaultMongoSecretRepositoryTimeout = 10 * time.Second
)

// MongoSecretRepository is a MongoDB implementation of a SecretRepository.
type MongoSecretRepository struct {
	client     mongo.Client
	collection string
	timeout    time.Duration
}

// MongoSecretRepositoryOptions is the options for the SecretRepository.
type MongoSecretRepositoryOptions struct {
	Database   string
	Collection string
	Timeout    time.Duration
}

// MongoSecretRepositoryOption is a function that sets options for the SecretRepository.
type MongoSecretRepositoryOption func(o *MongoSecretRepositoryOptions)

// NewMongoSecretRepository creates and configures a new SecretRepository.
func NewMongoSecretRepository(client mongo.Client, options ...MongoSecretRepositoryOption) (*MongoSecretRepository, error) {
	if client == nil {
		return nil, ErrNilClient
	}

	opts := MongoSecretRepositoryOptions{
		Database:   defaultMongoSecretRepositoryDatabase,
		Collection: defaultMongoSecretRepositoryCollection,
		Timeout:    defaultMongoSecretRepositoryTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	return &MongoSecretRepository{
		client:     client.Database(opts.Database),
		collection: opts.Collection,
	}, nil
}

// Get a secret by its ID.
func (r MongoSecretRepository) Get(ctx context.Context, id string) (Secret, error) {
	res, err := r.client.Collection(r.collection).FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return Secret{}, ErrSecretNotFound
		}
		return Secret{}, err
	}

	var secret Secret
	if err := res.Decode(&secret); err != nil {
		return Secret{}, err
	}
	return secret, nil
}

// Create a new secret.
func (r MongoSecretRepository) Create(ctx context.Context, secret Secret) (Secret, error) {
	id, err := r.client.Collection(r.collection).InsertOne(ctx, secret)
	if err != nil {
		return Secret{}, err
	}
	return r.Get(ctx, id)
}

// Delete a secret by its ID.
func (r MongoSecretRepository) Delete(ctx context.Context, id string) error {
	secret, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.client.Collection(r.collection).DeleteOne(ctx, bson.D{{Key: "_id", Value: secret.ID}})
}

// DeleteExpired deletes all expired secrets.
func (r MongoSecretRepository) DeleteExpired(ctx context.Context) error {
	filter := bson.D{{Key: "expiresAt", Value: bson.D{{Key: "$lt", Value: now()}}}}
	err := r.client.Collection(r.collection).DeleteMany(ctx, filter)
	if err != nil {
		if errors.Is(err, mongo.ErrDocumentsNotDeleted) {
			return ErrSecretsNotDeleted
		}
		return err
	}
	return nil
}

// Close the repository and its underlying connections.
func (r MongoSecretRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Disconnect(ctx)
}

// now is a function that returns the current time.
var now = func() time.Time {
	return time.Now()
}
