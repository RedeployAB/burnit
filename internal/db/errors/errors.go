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

var (
	// ErrSessionNotFound is returned when the session is not found.
	ErrSessionNotFound = errors.New("session not found")
	// ErrSessionNotDeleted is returned when the session is not deleted.
	ErrSessionNotDeleted = errors.New("session not deleted")
	// ErrSessionsNotDeleted is returned when sessions are not deleted.
	ErrSessionsNotDeleted = errors.New("sessions not deleted")
)
