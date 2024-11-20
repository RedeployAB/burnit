package errors

import "errors"

var (
	// ErrSecretNotFound is returned when the secret is not found.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrSecretNotDeleted is returned when the secret is not deleted.
	ErrSecretNotDeleted = errors.New("secret not deleted")
	// ErrSecretsNotDeleted is returned when secrets are not deleted.
	ErrSecretsNotDeleted = errors.New("secrets not deleted")
)
