package db

import (
	"context"
)

// SecretRepository defines the methods needed for persisting
// and retrieving secrets.
type SecretRepository interface {
	// Get a secret by its ID.
	Get(ctx context.Context, id string) (Secret, error)
	// Create a secret.
	Create(ctx context.Context, secret Secret) (Secret, error)
	// Delete a secret by its ID.
	Delete(ctx context.Context, id string) error
	// DeleteExpired deletes all expired secrets.
	DeleteExpired(ctx context.Context) error
	// Close the repository and its underlying connections.
	Close() error
}

// SessionRepository defines the methods needed for persisting
// and retrieving sessions.
type SessionRepository interface {
	// Get a session by its ID.
	Get(ctx context.Context, id string) (Session, error)
	// Create a session.
	Create(ctx context.Context, session Session) (Session, error)
	// Delete a session by its ID.
	Delete(ctx context.Context, id string) error
	// DeleteExpired deletes all expired sessions.
	DeleteExpired(ctx context.Context) error
	// Close the repository and its underlying connections.
	Close() error
}
