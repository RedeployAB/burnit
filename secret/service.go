package secret

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/RedeployAB/burnit/db"
	"github.com/RedeployAB/burnit/security"
)

const (
	// defaultTimeout is the default timeout for service operations.
	defaultTimeout = 10 * time.Second
)

const (
	// stdCharacters is the standard letters used for generating a secret.
	stdCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUWXYZ"
	// spcCharacters is the special characters used for generating a secret.
	spcCharacters = "_-!?=()&%"
	// maxLength is the maximum length of a secret.
	maxLength = 512
)

// Service is the interface that provides methods for secret operations.
type Service interface {
	Generate() string
	Get(id, passphrase string) (Secret, error)
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
	if length > maxLength {
		length = maxLength
	}

	var strb strings.Builder
	strb.WriteString(stdCharacters)
	if specialCharacters {
		strb.WriteString(spcCharacters)
	}
	bltrs := []byte(strb.String())
	b := make([]byte, length)
	for i := range b {
		b[i] = bltrs[rand.Intn(len(bltrs))]
	}

	return string(b)
}

// Get a secret. The secret is deleted after it has been retrieved
// and successfully decrypted.
func (s service) Get(id, passphrase string) (Secret, error) {
	dbSecret, err := s.secrets.Get(id)
	if err != nil {
		if errors.Is(err, db.ErrSecretNotFound) {
			return Secret{}, ErrSecretNotFound
		}
		return Secret{}, err
	}

	key := passphrase
	if len(key) == 0 {
		key = s.encryptionKey
	}

	decrypted, err := security.Decrypt([]byte(dbSecret.Value), []byte(key))
	if err != nil {
		return Secret{}, err
	}

	if err := s.secrets.Delete(id); err != nil {
		return Secret{}, err
	}

	return Secret{
		ID:    dbSecret.ID,
		Value: string(decrypted),
	}, nil
}
