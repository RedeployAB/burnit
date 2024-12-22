package db

import (
	"context"
)

// SecretStore defines the methods needed for persisting
// and retrieving secrets.
type SecretStore interface {
	// Get a secret by its ID.
	Get(ctx context.Context, id string) (Secret, error)
	// Create a secret.
	Create(ctx context.Context, secret Secret) (Secret, error)
	// Delete a secret by its ID.
	Delete(ctx context.Context, id string) error
	// DeleteExpired deletes all expired secrets.
	DeleteExpired(ctx context.Context) error
	// Close the SecretStore and its underlying connections.
	Close() error
}

// SessionStore defines the methods needed for persisting
// and retrieving sessions.
type SessionStore interface {
	// Get a session by its ID.
	Get(ctx context.Context, id string) (Session, error)
	// Get a session by its CSRF token.
	GetByCSRFToken(ctx context.Context, token string) (Session, error)
	// Upsert a session. Create the session if it does not exist, otherwise
	// update it.
	Upsert(ctx context.Context, session Session) (Session, error)
	// Delete a session by its ID.
	Delete(ctx context.Context, id string) error
	// DeleteExpired deletes all expired sessions.
	DeleteExpired(ctx context.Context) error
	// Close the SecretStore and its underlying connections.
	Close() error
}
