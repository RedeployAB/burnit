package inmem

import (
	"context"
	"sync"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
)

// secretStore is an in-memory store for secrets.
type secretStore struct {
	secrets map[string]db.Secret
	mu      sync.RWMutex
}

// NewSecretStore creates a new in-memory secret store.
func NewSecretStore() *secretStore {
	s := &secretStore{
		secrets: make(map[string]db.Secret),
		mu:      sync.RWMutex{},
	}
	return s
}

// Get a secret by its ID.
func (s *secretStore) Get(ctx context.Context, id string) (db.Secret, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	secret, ok := s.secrets[id]
	if !ok {
		return db.Secret{}, dberrors.ErrSecretNotFound
	}

	return secret, nil
}

// Create a secret.
func (s *secretStore) Create(ctx context.Context, secret db.Secret) (db.Secret, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.secrets[secret.ID] = db.Secret{
		ID:        secret.ID,
		Value:     secret.Value,
		ExpiresAt: secret.ExpiresAt,
	}

	return s.secrets[secret.ID], nil
}

// Delete a secret by its ID.
func (s *secretStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.secrets[id]; !ok {
		return dberrors.ErrSecretNotFound
	}

	delete(s.secrets, id)

	return nil
}

// DeleteExpired deletes all expired secrets.
// Note: The current implementation is very inefficient.
func (s *secretStore) DeleteExpired(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, secret := range s.secrets {
		if secret.ExpiresAt.Before(now()) {
			delete(s.secrets, id)
		}
	}

	return nil
}

// Close the store and its underlying connections.
func (s *secretStore) Close() error {
	return nil
}
