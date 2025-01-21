package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	// defaultSecretStoreDatabase is the default database for the SecretStore.
	defaultSecretStoreDatabase = "burnit"
	// defaultSecretStoreCollection is the default collection for the SecretStore.
	defaultSecretStoreCollection = "secrets"
	// defaultSecretStoreTimeout is the default timeout for the SecretStore.
	defaultSecretStoreTimeout = 10 * time.Second
)

// secretStore is a MongoDB implementation of a SecretStore.
type secretStore struct {
	client       Client
	collection   string
	createSecret createSecretFunc
	timeout      time.Duration
}

// SecretStoreOptions is the options for the SecretStore.
type SecretStoreOptions struct {
	Database   string
	Collection string
	Timeout    time.Duration
}

// SecretStoreOption is a function that sets options for the SecetStore.
type SecretStoreOption func(o *SecretStoreOptions)

// NewSecretStore creates and configures a new SecretStore.
func NewSecretStore(client Client, options ...SecretStoreOption) (*secretStore, error) {
	if client == nil {
		return nil, errors.New("nil client")
	}

	opts := SecretStoreOptions{
		Database:   defaultSecretStoreDatabase,
		Collection: defaultSecretStoreCollection,
		Timeout:    defaultSecretStoreTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	if len(opts.Database) == 0 {
		return nil, errors.New("database not set")
	}

	store := &secretStore{
		client:     client.Database(opts.Database),
		collection: opts.Collection,
		timeout:    opts.Timeout,
	}

	if client.ReplicaSetEnabled() {
		setCreateSecretWithTransaction(store)
	} else {
		setCreateSecret(store)
	}

	return store, nil
}

// Get a secret by its ID.
func (s secretStore) Get(ctx context.Context, id string) (db.Secret, error) {
	res, err := s.client.Collection(s.collection).FindOne(ctx, bson.D{{Key: "_id", Value: id}})
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

// Create a secret.
func (s secretStore) Create(ctx context.Context, secret db.Secret) (db.Secret, error) {
	return s.createSecret(ctx, secret)
}

// Delete a secret by its ID.
func (s secretStore) Delete(ctx context.Context, id string) error {
	if err := s.client.Collection(s.collection).DeleteOne(ctx, bson.D{{Key: "_id", Value: id}}); err != nil {
		if errors.Is(err, ErrDocumentNotDeleted) {
			return dberrors.ErrSecretNotFound
		}
		return err
	}
	return nil
}

// DeleteExpired deletes all expired secrets.
func (s secretStore) DeleteExpired(ctx context.Context) error {
	filter := bson.D{{Key: "expiresAt", Value: bson.D{{Key: "$lt", Value: now()}}}}
	err := s.client.Collection(s.collection).DeleteMany(ctx, filter)
	if err != nil {
		if errors.Is(err, ErrDocumentsNotDeleted) {
			return dberrors.ErrSecretsNotDeleted
		}
		return err
	}
	return nil
}

// Close the store and its underlying connections.
func (s secretStore) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.client.Disconnect(ctx)
}

// createSecretFunc is a function that creates a secret.
type createSecretFunc func(ctx context.Context, secret db.Secret) (db.Secret, error)

// setCreateSecret sets a createSecretFunc to the store.
func setCreateSecret(store *secretStore) {
	store.createSecret = func(ctx context.Context, secret db.Secret) (db.Secret, error) {
		id, err := store.client.Collection(store.collection).InsertOne(ctx, secret)
		if err != nil {
			return db.Secret{}, err
		}
		return store.Get(ctx, id)
	}
}

// setCreateSecretWithTransaction sets a createSecretFunc for use with transactions
// to the store.
func setCreateSecretWithTransaction(store *secretStore) {
	store.createSecret = func(ctx context.Context, secret db.Secret) (db.Secret, error) {
		result, err := store.client.WithTransaction(ctx, func(ctx context.Context, client Client) (any, error) {
			id, err := store.client.Collection(store.collection).InsertOne(ctx, secret)
			if err != nil {
				return nil, err
			}

			res, err := store.client.Collection(store.collection).FindOne(ctx, bson.D{{Key: "_id", Value: id}})
			if err != nil {
				if errors.Is(err, ErrNoDocuments) {
					return db.Secret{}, err
				}
				return nil, err
			}

			var secret db.Secret
			if err := res.Decode(&secret); err != nil {
				return db.Secret{}, err
			}
			return secret, nil
		})
		if err != nil {
			return db.Secret{}, err
		}

		secret, ok := result.(db.Secret)
		if !ok {
			return db.Secret{}, errors.New("invalid document for secret")
		}
		return secret, nil
	}
}
