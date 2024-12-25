package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
)

// secretStore is a Redis implementation of a SecretStore.
type secretStore struct {
	client Client
}

// SecretStoreOptions is the options for the SecretStore.
type SecretStoreOptions struct{}

// SecretStoreOption is a function that sets options for the SecretStore.
type SecretStoreOption func(o *SecretStoreOptions)

// NewSecretStore creates and configures a new SecretStore.
func NewSecretStore(client Client, options ...SecretStoreOption) (*secretStore, error) {
	if client == nil {
		return nil, ErrNilClient
	}

	opts := SecretStoreOptions{}
	for _, option := range options {
		option(&opts)
	}

	return &secretStore{
		client: client,
	}, nil
}

// Get a secret by its ID.
func (r secretStore) Get(ctx context.Context, id string) (db.Secret, error) {
	data, err := r.client.Get(ctx, id)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return db.Secret{}, dberrors.ErrSecretNotFound
		}
		return db.Secret{}, err
	}
	secret, err := secretFromJSON(data)
	if err != nil {
		return db.Secret{}, err
	}
	return secret, nil
}

// Create a secret.
func (r secretStore) Create(ctx context.Context, secret db.Secret) (db.Secret, error) {
	if err := r.client.Set(ctx, secret.ID, secretToJSON(secret), time.Until(secret.ExpiresAt)); err != nil {
		return db.Secret{}, err
	}
	return r.Get(ctx, secret.ID)
}

// Delete a secret by its ID.
func (r secretStore) Delete(ctx context.Context, id string) error {
	if err := r.client.Delete(ctx, id); err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return dberrors.ErrSecretNotFound
		}
		return err
	}
	return nil
}

// DeleteExpired deletes all expired secrets. This is a no-op for Redis
// since Redis handles expiration automatically.
func (r secretStore) DeleteExpired(ctx context.Context) error {
	return nil
}

// Close the store and its underlying connections.
func (r secretStore) Close() error {
	return r.client.Close()
}

// secretToJSON converts a secret to JSON.
func secretToJSON(s db.Secret) []byte {
	b, _ := json.Marshal(s)
	return b
}

// secretFromJSON converts JSON to a secret.
func secretFromJSON(b []byte) (db.Secret, error) {
	var s db.Secret
	if err := json.Unmarshal(b, &s); err != nil {
		return db.Secret{}, err
	}
	return s, nil
}
