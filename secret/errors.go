package secret

import "errors"

var (
	// ErrNilRepository is returned no repository is provided.
	ErrNilRepository = errors.New("nil repository")
	// ErrSecretNotFound is returned when a secret is not found.
	ErrSecretNotFound = errors.New("secret not found")
)
