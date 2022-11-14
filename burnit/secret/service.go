package secret

import (
	"github.com/RedeployAB/burnit/burnit/db"
)

const maxLength = 5000

// Service provides handling operations for secrets.
type Service interface {
	Get(id, passphrase string) (*Secret, error)
	Create(secret *Secret) (*Secret, error)
	Delete(id string) (int64, error)
	DeleteExpired() (int64, error)
	Generate(l int, sc bool) *Secret
}

// service is the layer that handles translation
// of incoming JSON payload to Secret, to Secret
// compatible with database.
type service struct {
	secrets db.Repository
	options Options
}

// Options for Service.
type Options struct {
	EncryptionKey string
}

// NewService creates a new service for handling secrets.
func NewService(secrets db.Repository, opts Options) Service {
	return &service{secrets: secrets, options: opts}
}

// Gets a secret from the repository by ID.
func (svc service) Get(id, passphrase string) (*Secret, error) {
	model, err := svc.secrets.Find(id)
	if err != nil || model == nil {
		return nil, err
	}

	var encryptionKey string
	if len(passphrase) > 0 {
		encryptionKey = passphrase
	} else {
		encryptionKey = svc.options.EncryptionKey
	}

	return toSecret(model, encryptionKey), nil
}

// Create a secret in the repository.
func (svc service) Create(s *Secret) (*Secret, error) {
	var encryptionKey string
	if len(s.Passphrase) > 0 {
		encryptionKey = s.Passphrase
	} else {
		encryptionKey = svc.options.EncryptionKey
	}
	model, err := svc.secrets.Insert(toModel(s, encryptionKey))
	if err != nil {
		return nil, err
	}
	return toSecret(model, svc.options.EncryptionKey), nil
}

// Delete a secret from the repository.
func (svc service) Delete(id string) (int64, error) {
	deleted, err := svc.secrets.Delete(id)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// Delete expired deletes all entries that has expiresAt
// less than current  time (time of invocation).
func (svc service) DeleteExpired() (int64, error) {
	deleted, err := svc.secrets.DeleteExpired()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// Generate a secret with specified length (l) and special characters (sc).
func (svc service) Generate(l int, sc bool) *Secret {
	value := Generate(l, sc)
	return &Secret{Value: value}
}
