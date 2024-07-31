package db

// SecretRepository defines the methods needed for persisting
// and retrieving secrets.
type SecretRepository interface {
	// Get a secret by its ID.
	Get(id string) (Secret, error)
	// Create a new secret.
	Create(secret Secret) (Secret, error)
	// Delete a secret by its ID.
	Delete(id string) error
	// DeleteExpired deletes all expired secrets.
	DeleteExpired() error
}
