package secret

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/db"
	"github.com/RedeployAB/burnit/security"
	"github.com/google/uuid"
)

const (
	// defaultTimeout is the default timeout for service operations.
	defaultTimeout = 10 * time.Second
	// defaultTTL is the default TTL of a secret.
	defaultTTL = 5 * time.Minute
)

// Service is the interface that provides methods for secret operations.
type Service interface {
	// Generate a new secret.
	Generate(length int, specialCharacters bool) string
	// Get a secret.
	Get(id, passphrase string) (Secret, error)
	// Create a secret.
	Create(secret Secret) (Secret, error)
	// Delete a secret.
	Delete(id string) error
	// Delete expired secrets.
	DeleteExpired() error
	// Close the service and its resources.
	Close() error
}

// service provides handling operations for secrets and satisfies Service.
type service struct {
	secrets       db.SecretRepository
	encryptionKey string
	timeout       time.Duration
}

// ServiceOption is a function that sets options for the service.
type ServiceOption func(s *service)

// NewService creates a new secret service.
func NewService(secrets db.SecretRepository, options ...ServiceOption) (*service, error) {
	if secrets == nil {
		return nil, ErrNilRepository
	}

	svc := &service{
		secrets: secrets,
		timeout: defaultTimeout,
	}

	for _, option := range options {
		option(svc)
	}

	return svc, nil
}

// Generate a new secret. The length of the secret is set by the provided
// length argument (with a max of 512 characters, a longer length will be trimmed to this value).
// If specialCharacters is set to true, the secret will contain special characters.
func (s service) Generate(length int, specialCharacters bool) string {
	return Generate(length, specialCharacters)
}

// Get a secret. The secret is deleted after it has been retrieved
// and successfully decrypted.
func (s service) Get(id, passphrase string) (Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	dbSecret, err := s.secrets.Get(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrSecretNotFound) {
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

	key := passphrase
	if len(key) == 0 {
		key = s.encryptionKey
	}

	decrypted, err := decrypt(dbSecret.Value, key)
	if err != nil {
		if errors.Is(err, security.ErrInvalidKey) {
			return Secret{}, ErrInvalidPassphrase
		}
		return Secret{}, err
	}

	if err := s.secrets.Delete(ctx, id); err != nil {
		return Secret{}, err
	}

	return Secret{
		ID:    dbSecret.ID,
		Value: string(decrypted),
	}, nil
}

// Create a secret.
func (s service) Create(secret Secret) (Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if secret.TTL == 0 {
		secret.TTL = defaultTTL
	}

	key := secret.Passphrase
	if len(key) == 0 {
		key = s.encryptionKey
	}

	encrypted, err := encrypt(secret.Value, key)
	if err != nil {
		return Secret{}, err
	}

	dbSecret, err := s.secrets.Create(ctx, db.Secret{
		ID:        newUUID(),
		Value:     encrypted,
		ExpiresAt: now().Add(secret.TTL),
	})
	if err != nil {
		return Secret{}, err
	}

	return Secret{
		ID:        dbSecret.ID,
		TTL:       secret.TTL,
		ExpiresAt: dbSecret.ExpiresAt,
	}, nil
}

// Delete a secret.
func (s service) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	return s.secrets.Delete(ctx, id)
}

// DeleteExpired deletes all expired secrets.
func (s service) DeleteExpired() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if err := s.secrets.DeleteExpired(ctx); err != nil && !errors.Is(err, db.ErrSecretsNotDeleted) {
		return err
	}

	return s.secrets.DeleteExpired(ctx)
}

// Close the service and its resources.
func (s service) Close() error {
	return s.secrets.Close()
}

// newUUID generates a new UUID.
var newUUID = func() string {
	return uuid.New().String()
}

// encrypt a value using a key and returns the encrypted value
// as a base64 encoded string.
func encrypt(value, key string) (string, error) {
	encrypted, err := security.Encrypt([]byte(value), security.ToMD5([]byte(key)))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// decrypt a value using a key and returns the decrypted value
// as a string.
func decrypt(value, key string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	decrypted, err := security.Decrypt(decoded, security.ToMD5([]byte(key)))
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

// now returns the current time.
var now = func() time.Time {
	return time.Now()
}
