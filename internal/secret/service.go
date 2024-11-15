package secret

import (
	"context"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
	"github.com/RedeployAB/burnit/internal/security"
	"github.com/google/uuid"
)

const (
	// defaultTTL is the default TTL of a secret.
	defaultTTL = 1 * time.Hour
	// defaultMinTTL is the default minimum TTL of a secret.
	defaultMinTTL = 1*time.Minute - 5*time.Second
	// defaultMaxTTL is the default maximum TTL of a secret.
	defaultMaxTTL = 168*time.Hour + 5*time.Second
	// defaultTimeout is the default timeout for service operations.
	defaultTimeout = 10 * time.Second
	// defaultCleanupInterval is the default interval for cleaning up expired secrets.
	defaultCleanupInterval = 30 * time.Second
)

const (
	// defaultPassphraseLength is the default length of a passphrase.
	defaultPassphraseLength = 32
)

const (
	// defaultMaxValueCharacters is the maximum number of characters (runes) in a secret.
	defaultMaxValueCharacters = 3500
)

// newUUID generates a new UUID.
var newUUID = func() string {
	return uuid.New().String()
}

// Service is the interface that provides methods for secret operations.
type Service interface {
	// Start the service and initialize its resources.
	Start() error
	// Close the service and its resources.
	Close() error
	// Generate a new secret.
	Generate(length int, specialCharacters bool) string
	// Get a secret.
	Get(id, passphrase string, options ...GetOption) (Secret, error)
	// Create a secret.
	Create(secret Secret) (Secret, error)
	// Delete a secret.
	Delete(id string, options ...DeleteOption) error
	// Delete expired secrets.
	DeleteExpired() error
}

// service provides handling operations for secrets and satisfies Service.
type service struct {
	secrets            db.SecretRepository
	timeout            time.Duration
	cleanupInterval    time.Duration
	maxValueCharacters int
	stopCh             chan struct{}
}

// ServiceOption is a function that sets options for the service.
type ServiceOption func(s *service)

// NewService creates a new secret service.
func NewService(secrets db.SecretRepository, options ...ServiceOption) (*service, error) {
	if secrets == nil {
		return nil, ErrNilRepository
	}

	svc := &service{
		secrets:            secrets,
		timeout:            defaultTimeout,
		cleanupInterval:    defaultCleanupInterval,
		maxValueCharacters: defaultMaxValueCharacters,
		stopCh:             make(chan struct{}),
	}

	for _, option := range options {
		option(svc)
	}

	return svc, nil
}

// Start the service and initialize its resources.
func (s *service) Start() error {
	for {
		select {
		case <-time.After(s.cleanupInterval):
			if err := s.DeleteExpired(); err != nil {
				return err
			}
		case <-s.stopCh:
			close(s.stopCh)
			return nil
		}
	}
}

// Close the service and its resources.
func (s *service) Close() error {
	s.stopCh <- struct{}{}
	return s.secrets.Close()
}

// Generate a new secret. The length of the secret is set by the provided
// length argument (with a max of 512 characters, a longer length will be trimmed to this value).
// If specialCharacters is set to true, the secret will contain special characters.
func (s service) Generate(length int, specialCharacters bool) string {
	return generate(length, specialCharacters)
}

// GetOptions contains options for getting a secret.
type GetOptions struct {
	NoDelete         bool
	PassphraseHashed bool
	context          context.Context
}

// GetOption is a function that sets options for getting a secret.
type GetOption func(o *GetOptions)

// Get a secret. The secret is deleted after it has been retrieved
// and successfully decrypted if the option to delete it is set.
func (s service) Get(id, passphrase string, options ...GetOption) (Secret, error) {
	opts := GetOptions{}
	for _, option := range options {
		option(&opts)
	}

	var ctx context.Context
	if opts.context != nil {
		ctx = opts.context
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.timeout)
		defer cancel()
	}

	dbSecret, err := s.secrets.Get(ctx, id)
	if err != nil {
		if errors.Is(err, dberrors.ErrSecretNotFound) {
			return Secret{}, ErrSecretNotFound
		}
		return Secret{}, err
	}

	if dbSecret.ExpiresAt.Before(now()) {
		if err := s.secrets.Delete(ctx, id); err != nil {
			return Secret{}, err
		}
		return Secret{}, ErrSecretNotFound
	}

	decrypted, err := decrypt(dbSecret.Value, passphrase, opts.PassphraseHashed)
	if err != nil {
		if errors.Is(err, security.ErrInvalidKey) {
			return Secret{}, ErrInvalidPassphrase
		}
		return Secret{}, err
	}

	secret := Secret{
		ID:    dbSecret.ID,
		Value: string(decrypted),
	}

	if opts.NoDelete {
		return secret, nil
	}

	if err := s.secrets.Delete(ctx, id); err != nil {
		return secret, err
	}

	return secret, nil
}

// Create a secret.
func (s service) Create(secret Secret) (Secret, error) {
	if err := validValue(secret.Value, s.maxValueCharacters); err != nil {
		return Secret{}, err
	}

	expiresAt, err := expirationTime(secret.TTL, secret.ExpiresAt)
	if err != nil {
		return Secret{}, err
	}

	passphrase := secret.Passphrase
	if len(passphrase) == 0 {
		passphrase = generate(defaultPassphraseLength, true)
	}

	encrypted, err := encrypt(secret.Value, passphrase)
	if err != nil {
		return Secret{}, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	dbSecret, err := s.secrets.Create(ctx, db.Secret{
		ID:        newUUID(),
		Value:     encrypted,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return Secret{}, err
	}

	return Secret{
		ID:         dbSecret.ID,
		Passphrase: passphrase,
		TTL:        time.Until(dbSecret.ExpiresAt).Round(time.Minute),
		ExpiresAt:  dbSecret.ExpiresAt,
	}, nil
}

// DeleteOptions contains options for deleting a secret.
type DeleteOptions struct {
	Passphrase       string
	VerifyPassphrase bool
	PassphraseHashed bool
}

// DeleteOption is a function that sets options for deleting a secret.
type DeleteOption func(o *DeleteOptions)

// Delete a secret.
func (s service) Delete(id string, options ...DeleteOption) error {
	opts := DeleteOptions{}
	for _, option := range options {
		option(&opts)
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if opts.VerifyPassphrase {
		_, err := s.Get(id, opts.Passphrase, func(o *GetOptions) {
			o.PassphraseHashed = opts.PassphraseHashed
			o.context = ctx
		})
		if err != nil {
			return err
		}
		return nil
	}

	err := s.secrets.Delete(ctx, id)
	if err == nil {
		return nil
	}

	if errors.Is(err, dberrors.ErrSecretNotFound) || errors.Is(err, dberrors.ErrSecretNotDeleted) {
		return ErrSecretNotFound
	}
	return err
}

// DeleteExpired deletes all expired secrets.
func (s service) DeleteExpired() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	err := s.secrets.DeleteExpired(ctx)
	if err == nil {
		return nil
	}
	if errors.Is(err, dberrors.ErrSecretsNotDeleted) {
		return nil
	}

	return err
}

// expirationTime returns the expiration time of a secret. It
// validates the provided duration and expiration time and returns
// the expiration time based on the provided values.
func expirationTime(ttl time.Duration, expiresAt time.Time) (time.Time, error) {
	current := now()
	n := current
	if !expiresAt.IsZero() {
		n = expiresAt
		if n.Before(current) {
			return time.Time{}, fmt.Errorf("%w: must be in the future", ErrInvalidExpirationTime)
		}
	} else if ttl > 0 {
		n = n.Add(ttl)
	} else {
		return n.Add(defaultTTL), nil
	}

	if n.Before(current.Add(defaultMinTTL)) || n.After(current.Add(defaultMaxTTL)) {
		return time.Time{}, fmt.Errorf("%w: must be between 1 minutes and 7 days", ErrInvalidExpirationTime)
	}

	return n, nil
}

// encrypt a value using a key and returns the encrypted value
// as a base64 encoded string and the hash as a base64 raw url encoded string.
func encrypt(value, key string) (string, error) {
	encrypted, err := security.Encrypt([]byte(value), security.SHA256([]byte(key)))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// decrypt a value using a key and returns the decrypted value
// as a string.
func decrypt(value, key string, hashed bool) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	var hash []byte
	if !hashed {
		hash = security.SHA256([]byte(key))
	} else {
		hash = []byte(key)
	}

	decrypted, err := security.Decrypt(decoded, hash)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// now returns the current time.
var now = func() time.Time {
	return time.Now().UTC()
}

// decoderFunc is a function that decodes a string and
// returns a byte slice and an error.
type decoderFunc func(string) ([]byte, error)

// validValue validates a secret value and returns an error if the value is invalid.
func validValue(value string, maxValueCharacters int) error {
	if len(value) == 0 {
		return fmt.Errorf("%w: must not be empty", ErrValueInvalid)
	}

	if utf8.RuneCountInString(value) > maxValueCharacters {
		return fmt.Errorf("%w: max characters are %d", ErrValueTooManyCharacters, maxValueCharacters)
	}

	v := []byte(value)

	decoderFuncs := []decoderFunc{
		base32.StdEncoding.DecodeString,
		base32.HexEncoding.DecodeString,
		hex.DecodeString,
		security.DecodeBase64,
	}

	for _, decode := range decoderFuncs {
		if decoded, err := decode(value); err == nil {
			v = decoded
			break
		}
	}

	if !validString(v) {
		return fmt.Errorf("%w: must be a valid non-empty UTF-8 encoded string", ErrValueInvalid)
	}

	return nil
}

// validString returns true if the byte slice is a string.
func validString(b []byte) bool {
	if nullByte(b) {
		return false
	}
	return http.DetectContentType(b) == "text/plain; charset=utf-8" && string(b) != "\x00"
}

// nullByte returns true if the byte slice is a null byte.
func nullByte(b []byte) bool {
	return len(b) == 1 && b[0] == 0
}
