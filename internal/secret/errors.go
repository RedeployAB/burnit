package secret

import "errors"

var (
	// ErrNilRepository is returned no repository is provided.
	ErrNilRepository = errors.New("nil repository")
	// ErrSecretNotFound is returned when a secret is not found.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrSecretInvalid is returned when the secret is invalid.
	ErrSecretInvalid = errors.New("secret invalid")
	// ErrInvalidPassphrase is returned when the passphrase is invalid.
	ErrInvalidPassphrase = errors.New("invalid passphrase")
	// ErrInvalidExpirationTime is returned when the expiration time is invalid.
	ErrInvalidExpirationTime = errors.New("invalid expiration time")
	// ErrSecretTooManyBytes is returned when the secret has too many bytes.
	ErrSecretTooManyBytes = errors.New("secret has too many bytes")
	// ErrSecretTooManyCharacters is returned when the secret has too many characters.
	ErrSecretTooManyCharacters = errors.New("secret has too many characters")
	// ErrPassphraseNotBase64 is returned when the passphrase is not base64 encoded.
	ErrPassphraseNotBase64 = errors.New("passphrase not base64 encoded")
)
