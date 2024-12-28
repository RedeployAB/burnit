package redis

import (
	"context"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
)

const (
	// secretPrefix is the key prefix for secrets.
	secretPrefix = "secret:"
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
func (s secretStore) Get(ctx context.Context, id string) (db.Secret, error) {
	data, err := s.client.HGet(ctx, secretPrefix+id)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return db.Secret{}, dberrors.ErrSecretNotFound
		}
		return db.Secret{}, err
	}
	return secretFromMap(data)
}

// Create a secret.
func (s secretStore) Create(ctx context.Context, secret db.Secret) (db.Secret, error) {
	result, err := s.client.WithTransaction(ctx, func(tx Tx) {
		tx.HSet(ctx, secretPrefix+secret.ID, secretToMap(&secret))
		tx.Expire(ctx, secretPrefix+secret.ID, time.Until(secret.ExpiresAt))
		tx.HGet(ctx, secretPrefix+secret.ID)
	})
	if err != nil {
		return db.Secret{}, nil
	}

	data := result.LastMap()
	if data == nil {
		return db.Secret{}, ErrKeyNotFound
	}
	return secretFromMap(data)
}

// Delete a secret by its ID.
func (s secretStore) Delete(ctx context.Context, id string) error {
	if err := s.client.Delete(ctx, secretPrefix+id); err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return dberrors.ErrSecretNotFound
		}
		return err
	}
	return nil
}

// DeleteExpired deletes all expired secrets. This is a no-op for Redis
// since Redis handles expiration automatically.
func (s secretStore) DeleteExpired(ctx context.Context) error {
	return nil
}

// Close the store and its underlying connections.
func (s secretStore) Close() error {
	return s.client.Close()
}

// secretToMap creates a map from the provided secret.
func secretToMap(secret *db.Secret) map[string]any {
	return map[string]any{
		"id":         secret.ID,
		"value":      secret.Value,
		"expires_at": secret.ExpiresAt,
	}
}

// secretFromMap creates a db.Secret from the provided map.
func secretFromMap(secret map[string]string) (db.Secret, error) {
	expiresAt, err := time.Parse(time.RFC3339, secret["expires_at"])
	if err != nil {
		return db.Secret{}, err
	}
	return db.Secret{
		ID:        secret["id"],
		Value:     secret["value"],
		ExpiresAt: expiresAt,
	}, nil
}
