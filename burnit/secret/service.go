package secret

import (
	"context"
	"fmt"
	"time"

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
	Start() error
	Stop() error
}

// Repository defined the methods needed for interact
// with a database and collection.
type Repository interface {
	// Get a secret.
	Get(id string) (*db.Secret, error)
	// Create a secret.
	Create(s *db.Secret) (*db.Secret, error)
	// Delete a secret.
	Delete(id string) (int64, error)
	// DeleteExpired deletes expired secrets.
	DeleteExpired() (int64, error)
	// Client returns the underlying database client.
	Client() db.Client
}

// service is the layer that handles translation
// of incoming JSON payload to Secret, to Secret
// compatible with database.
type service struct {
	secrets       Repository
	timeout       time.Duration
	encryptionKey string
}

// Options for Service.
type ServiceOptions struct {
	Secrets       Repository
	EncryptionKey string
	Timeout       time.Duration
}

// NewService creates a new service for handling secrets.
func NewService(opts *ServiceOptions) *service {
	if opts == nil {
		opts = &ServiceOptions{}
	}

	if opts.Timeout == 0 {
		opts.Timeout = time.Second * 30
	}

	return &service{
		secrets:       opts.Secrets,
		timeout:       opts.Timeout,
		encryptionKey: opts.EncryptionKey,
	}
}

// Gets a secret from the repository by ID.
func (svc service) Get(id, passphrase string) (*Secret, error) {
	model, err := svc.secrets.Get(id)
	if err != nil || model == nil {
		return nil, err
	}

	var encryptionKey string
	if len(passphrase) > 0 {
		encryptionKey = passphrase
	} else {
		encryptionKey = svc.encryptionKey
	}

	return toSecret(model, encryptionKey), nil
}

// Create a secret in the repository.
func (svc service) Create(s *Secret) (*Secret, error) {
	var encryptionKey string
	if len(s.Passphrase) > 0 {
		encryptionKey = s.Passphrase
	} else {
		encryptionKey = svc.encryptionKey
	}
	model, err := svc.secrets.Create(toModel(s, encryptionKey))
	if err != nil {
		return nil, err
	}
	return toSecret(model, svc.encryptionKey), nil
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

// Start the service and connect to the repository and database.
func (svc service) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), svc.timeout)
	defer cancel()

	if err := svc.secrets.Client().Connect(ctx); err != nil {
		return fmt.Errorf("connecting to database: %v", err)
	}
	return nil
}

// Stop the service and disconnect from the repository and database.
func (svc service) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), svc.timeout)
	defer cancel()

	return svc.secrets.Client().Disconnect(ctx)
}
