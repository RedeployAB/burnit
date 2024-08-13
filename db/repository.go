package db

import (
	"context"

	"github.com/RedeployAB/burnit/db/models"
)

// SecretRepository defines the methods needed for persisting
// and retrieving secrets.
type SecretRepository interface {
	// Get a secret by its ID.
	Get(ctx context.Context, id string) (models.Secret, error)
	// Create a new secret.
	Create(ctx context.Context, secret models.Secret) (models.Secret, error)
	// Delete a secret by its ID.
	Delete(ctx context.Context, id string) error
	// DeleteExpired deletes all expired secrets.
	DeleteExpired(ctx context.Context) error
	// GetSettings gets the settings.
	GetSettings(ctx context.Context) (models.Settings, error)
	// CreateSettings creates settings.
	CreateSettings(ctx context.Context, settings models.Settings) (models.Settings, error)
	// UpdateSettings updates the settings.
	UpdateSettings(ctx context.Context, settings models.Settings) (models.Settings, error)
	// Close the repository and its underlying connections.
	Close() error
}
