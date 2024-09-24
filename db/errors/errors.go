package errors

import "errors"

var (
	// ErrSecretNotFound is returned when the secret is not found.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrSecretsNotDeleted is returned when secrets are not deleted.
	ErrSecretsNotDeleted = errors.New("secrets not deleted")
)

var (
	// ErrSettingsNotFound is returned when the settings are not found.
	ErrSettingsNotFound = errors.New("settings not found")
)
