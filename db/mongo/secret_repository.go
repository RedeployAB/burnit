package mongo

import (
	"context"
	"errors"
	"time"

	dberrors "github.com/RedeployAB/burnit/db/errors"
	"github.com/RedeployAB/burnit/db/models"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// defaultSecretRepositoryDatabase is the default database for the secret repository.
	defaultSecretRepositoryDatabase = "burnit"
	// defaultSecretRepositoryCollection is the default collection for the secret repository.
	defaultSecretRepositoryCollection = "secrets"
	// defaultSecretRepositoryTimeout is the default timeout for the secret repository.
	defaultSecretRepositoryTimeout = 10 * time.Second
)

// SecretRepository is a MongoDB implementation of a SecretRepository.
type SecretRepository struct {
	client     Client
	collection string
	timeout    time.Duration
}

// SecretRepositoryOptions is the options for the SecretRepository.
type SecretRepositoryOptions struct {
	Database   string
	Collection string
	Timeout    time.Duration
}

// SecretRepositoryOption is a function that sets options for the SecretRepository.
type SecretRepositoryOption func(o *SecretRepositoryOptions)

// NewSecretRepository creates and configures a new SecretRepository.
func NewSecretRepository(client Client, options ...SecretRepositoryOption) (*SecretRepository, error) {
	if client == nil {
		return nil, ErrNilClient
	}

	opts := SecretRepositoryOptions{
		Database:   defaultSecretRepositoryDatabase,
		Collection: defaultSecretRepositoryCollection,
		Timeout:    defaultSecretRepositoryTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	return &SecretRepository{
		client:     client.Database(opts.Database),
		collection: opts.Collection,
		timeout:    opts.Timeout,
	}, nil
}

// Get a secret by its ID.
func (r SecretRepository) Get(ctx context.Context, id string) (models.Secret, error) {
	res, err := r.client.Collection(r.collection).FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		if errors.Is(err, ErrNoDocuments) {
			return models.Secret{}, dberrors.ErrSecretNotFound
		}
		return models.Secret{}, err
	}

	var secret models.Secret
	if err := res.Decode(&secret); err != nil {
		return models.Secret{}, err
	}
	return secret, nil
}

// Create a new secret.
func (r SecretRepository) Create(ctx context.Context, secret models.Secret) (models.Secret, error) {
	id, err := r.client.Collection(r.collection).InsertOne(ctx, secret)
	if err != nil {
		return models.Secret{}, err
	}
	return r.Get(ctx, id)
}

// Delete a secret by its ID.
func (r SecretRepository) Delete(ctx context.Context, id string) error {
	secret, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	return r.client.Collection(r.collection).DeleteOne(ctx, bson.D{{Key: "_id", Value: secret.ID}})
}

// DeleteExpired deletes all expired secrets.
func (r SecretRepository) DeleteExpired(ctx context.Context) error {
	filter := bson.D{{Key: "expiresAt", Value: bson.D{{Key: "$lt", Value: now()}}}}
	err := r.client.Collection(r.collection).DeleteMany(ctx, filter)
	if err != nil {
		if errors.Is(err, ErrDocumentsNotDeleted) {
			return dberrors.ErrSecretsNotDeleted
		}
		return err
	}
	return nil
}

// Close the repository and its underlying connections.
func (r SecretRepository) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.client.Disconnect(ctx)
}

// now is a function that returns the current time.
var now = func() time.Time {
	return time.Now()
}
