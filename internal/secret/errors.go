package secret

import "errors"

var (
	// ErrNilRepository is returned no repository is provided.
	ErrNilRepository = errors.New("nil repository")
	// ErrSecretNotFound is returned when a secret is not found.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrInvalidPassphrase is returned when the passphrase is invalid for a secret.
	ErrInvalidPassphrase = errors.New("invalid passphrase")
	// ErrValueInvalid is returned when the secret value is invalid.
	ErrValueInvalid = errors.New("value invalid")
	// ErrValueTooManyCharacters is returned when the secret value has too many characters.
	ErrValueTooManyCharacters = errors.New("value has too many characters")
	// ErrInvalidExpirationTime is returned when the expiration time is invalid.
	ErrInvalidExpirationTime = errors.New("invalid expiration time")
	// ErrPassphraseNotBase64 is returned when the passphrase is not base64 encoded.
	ErrPassphraseNotBase64 = errors.New("passphrase not base64 encoded")
	// ErrPassphraseInvalid is returned when the passphrase input is invalid.
	ErrPassphraseInvalid = errors.New("passphrase invalid")
	// ErrPassphraseTooManyCharacters is returned when the passphrase has too many characters.
	ErrPassphraseTooManyCharacters = errors.New("passphrase has too many characters")
	// ErrPassphraseTooFewCharacters is returned when the passphrase has too few characters.
	ErrPassphraseTooFewCharacters = errors.New("passphrase has too few characters")
)
