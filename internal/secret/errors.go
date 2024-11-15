package secret

import "errors"

var (
	// ErrNilRepository is returned no repository is provided.
	ErrNilRepository = errors.New("nil repository")
	// ErrSecretNotFound is returned when a secret is not found.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrInvalidPassphrase is returned when the passphrase is invalid.
	ErrInvalidPassphrase = errors.New("invalid passphrase")
	// ErrValueInvalid is returned when the secret value is invalid.
	ErrValueInvalid = errors.New("value invalid")
	// ErrValueTooManyCharacters is returned when the secret value has too many characters.
	ErrValueTooManyCharacters = errors.New("value has too many characters")
	// ErrInvalidExpirationTime is returned when the expiration time is invalid.
	ErrInvalidExpirationTime = errors.New("invalid expiration time")
	// ErrPassphraseNotBase64 is returned when the passphrase is not base64 encoded.
	ErrPassphraseNotBase64 = errors.New("passphrase not base64 encoded")
)
