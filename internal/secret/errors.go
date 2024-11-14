package secret

import "errors"

var (
	// ErrNilRepository is returned no repository is provided.
	ErrNilRepository = errors.New("nil repository")
	// ErrSecretNotFound is returned when a secret is not found.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrInvalidSecretValue is returned when the secret value is invalid.
	ErrInvalidSecretValue = errors.New("invalid secret value")
	// ErrInvalidPassphrase is returned when the passphrase is invalid.
	ErrInvalidPassphrase = errors.New("invalid passphrase")
	// ErrInvalidExpirationTime is returned when the expiration time is invalid.
	ErrInvalidExpirationTime = errors.New("invalid expiration time")
	// ErrSecretValueTooManyBytes is returned when the secret value has too many bytes.
	ErrSecretValueTooManyBytes = errors.New("secret value has too many bytes")
	// ErrSecretValueTooManyCharacters is returned when the secret value has too many characters.
	ErrSecretValueTooManyCharacters = errors.New("secret value has too many characters")
	// ErrPassphraseNotBase64 is returned when the passphrase is not base64 encoded.
	ErrPassphraseNotBase64 = errors.New("passphrase not base64 encoded")
)
