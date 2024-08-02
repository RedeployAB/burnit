package db

import "errors"

var (
	// ErrNilClient is returned when the client is nil.
	ErrNilClient = errors.New("client is nil")
	// ErrSecretNotFound is returned when the secret is not found.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrSecretsNotDeleted is returned when secrets are not deleted.
	ErrSecretsNotDeleted = errors.New("secrets not deleted")
)
