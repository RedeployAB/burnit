package secret

import (
	"context"
	"encoding/base64"
	"errors"
	"time"

	"github.com/RedeployAB/burnit/db"
	dberrors "github.com/RedeployAB/burnit/db/errors"
	"github.com/RedeployAB/burnit/security"
	"github.com/google/uuid"
)

const (
	// defaultTTL is the default TTL of a secret.
	defaultTTL = 5 * time.Minute
	// defaultTimeout is the default timeout for service operations.
	defaultTimeout = 10 * time.Second
	// defaultCleanupInterval is the default interval for cleaning up expired secrets.
	defaultCleanupInterval = 30 * time.Second
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
	Get(id, passphrase string) (Secret, error)
	// Create a secret.
	Create(secret Secret) (Secret, error)
	// Delete a secret.
	Delete(id string) error
	// Delete expired secrets.
	DeleteExpired() error
}

// service provides handling operations for secrets and satisfies Service.
type service struct {
	secrets         db.SecretRepository
	encryptionKey   string
	timeout         time.Duration
	cleanupInterval time.Duration
	stopCh          chan struct{}
}

// ServiceOption is a function that sets options for the service.
type ServiceOption func(s *service)

// NewService creates a new secret service.
func NewService(secrets db.SecretRepository, options ...ServiceOption) (*service, error) {
	if secrets == nil {
		return nil, ErrNilRepository
	}

	svc := &service{
		secrets:         secrets,
		timeout:         defaultTimeout,
		cleanupInterval: defaultCleanupInterval,
		stopCh:          make(chan struct{}),
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

// Get a secret. The secret is deleted after it has been retrieved
// and successfully decrypted.
func (s service) Get(id, passphrase string) (Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

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
	return time.Now().UTC()
}
