package db

import (
	"context"
)

// SecretRepository defines the methods needed for persisting
// and retrieving secrets.
type SecretRepository interface {
	// Get a secret by its ID.
	Get(ctx context.Context, id string) (Secret, error)
	// Create a new secret.
	Create(ctx context.Context, secret Secret) (Secret, error)
	// Delete a secret by its ID.
	Delete(ctx context.Context, id string) error
	// DeleteExpired deletes all expired secrets.
	DeleteExpired(ctx context.Context) error
	// Close the repository and its underlying connections.
	Close() error
}
