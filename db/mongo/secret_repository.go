package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/db"
	dberrors "github.com/RedeployAB/burnit/db/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// defaultSecretRepositoryDatabase is the default database for the secret repository.
	defaultSecretRepositoryDatabase = "burnit"
	// defaultSecretRepositoryCollection is the default collection for the secret repository.
	defaultSecretRepositoryCollection = "secrets"
	// defaultSettingsCollection is the default collection for the settings.
	defaultSettingsCollection = "settings"
	// defaultSecretRepositoryTimeout is the default timeout for the secret repository.
	defaultSecretRepositoryTimeout = 10 * time.Second
)

// SecretRepository is a MongoDB implementation of a SecretRepository.
type SecretRepository struct {
	client             Client
	collection         string
	settingsCollection string
	timeout            time.Duration
}

// SecretRepositoryOptions is the options for the SecretRepository.
type SecretRepositoryOptions struct {
	Database           string
	Collection         string
	SettingsCollection string
	Timeout            time.Duration
}

// SecretRepositoryOption is a function that sets options for the SecretRepository.
type SecretRepositoryOption func(o *SecretRepositoryOptions)

// NewSecretRepository creates and configures a new SecretRepository.
func NewSecretRepository(client Client, options ...SecretRepositoryOption) (*SecretRepository, error) {
	if client == nil {
		return nil, ErrNilClient
	}

	opts := SecretRepositoryOptions{
		Database:           defaultSecretRepositoryDatabase,
		Collection:         defaultSecretRepositoryCollection,
		SettingsCollection: defaultSettingsCollection,
		Timeout:            defaultSecretRepositoryTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	if len(opts.Database) == 0 {
		return nil, ErrDatabaseNotSet
	}

	return &SecretRepository{
		client:             client.Database(opts.Database),
		collection:         opts.Collection,
		settingsCollection: opts.SettingsCollection,
		timeout:            opts.Timeout,
	}, nil
}

// Get a secret by its ID.
func (r SecretRepository) Get(ctx context.Context, id string) (db.Secret, error) {
	res, err := r.client.Collection(r.collection).FindOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		if errors.Is(err, ErrNoDocuments) {
			return db.Secret{}, dberrors.ErrSecretNotFound
		}
		return db.Secret{}, err
	}

	var secret db.Secret
	if err := res.Decode(&secret); err != nil {
		return db.Secret{}, err
	}
	return secret, nil
}

// Create a new secret.
func (r SecretRepository) Create(ctx context.Context, secret db.Secret) (db.Secret, error) {
	id, err := r.client.Collection(r.collection).InsertOne(ctx, secret)
	if err != nil {
		return db.Secret{}, err
	}
	return r.Get(ctx, id)
}

// Delete a secret by its ID.
func (r SecretRepository) Delete(ctx context.Context, id string) error {
	secret, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	err = r.client.Collection(r.collection).DeleteOne(ctx, bson.D{{Key: "_id", Value: secret.ID}})
	if err != nil {
		if errors.Is(err, ErrDocumentNotDeleted) {
			return dberrors.ErrSecretNotDeleted
		}
	}
	return nil
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
	return time.Now().UTC()
}
