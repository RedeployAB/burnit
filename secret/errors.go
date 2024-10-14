package secret

import "errors"

var (
	// ErrNilRepository is returned no repository is provided.
	ErrNilRepository = errors.New("nil repository")
	// ErrSecretNotFound is returned when a secret is not found.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrInvalidPassphrase is returned when the passphrase is invalid.
	ErrInvalidPassphrase = errors.New("invalid passphrase")
	// ErrInvalidExpirationTime is returned when the expiration time is invalid.
	ErrInvalidExpirationTime = errors.New("invalid expiration time")
)
