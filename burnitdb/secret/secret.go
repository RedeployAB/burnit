package secret

import (
	"encoding/json"
	"io"
	"time"

	"github.com/RedeployAB/burnit/burnitdb/db"
	"github.com/RedeployAB/burnit/common/security"
)

// Secret is to be used as the middle-hand between
// incoming JSON payload and the data model.
type Secret struct {
	ID         string
	Value      string
	Passphrase string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	TTL        int
}

// Service provides handling operations for secrets.
type Service interface {
	Get(id, passphrase string) (*Secret, error)
	Create(secret *Secret) (*Secret, error)
	Delete(id string) (int64, error)
	DeleteExpired() (int64, error)
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

// NewService creates a new service for handling secrets.
func NewService(secrets db.Repository, opts Options) Service {
	return &service{secrets: secrets, options: opts}
}

// NewFromJSON creates an incoming JSON payload
// and creates a Secret from it.
func NewFromJSON(b io.ReadCloser) (*Secret, error) {
	var s *Secret
	if err := json.NewDecoder(b).Decode(&s); err != nil {
		return nil, err
	}

	if s.TTL != 0 {
		s.ExpiresAt = time.Now().Add(time.Minute * time.Duration(s.TTL))
	} else {
		s.ExpiresAt = time.Now().AddDate(0, 0, 7)
	}

	return s, nil
}

// toModel transforms a Secret to the
// data model variant of Secret.
func toModel(s *Secret, passphrase string) *db.Secret {
	var createdAt, expiresAt time.Time
	if s.CreatedAt.IsZero() {
		createdAt = time.Now()
	} else {
		createdAt = s.CreatedAt
	}

	if s.ExpiresAt.IsZero() {
		expiresAt = time.Now().Add(time.Minute * time.Duration(10080))
	} else {
		expiresAt = s.ExpiresAt
	}

	val, err := security.Encrypt([]byte(s.Value), passphrase)
	if err != nil {
		return nil
	}

	return &db.Secret{
		Value:     string(val),
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}

// toSecret transforms the data model variant
// of Secret to a Secret.
func toSecret(s *db.Secret, passphrase string) *Secret {
	var val string
	if len(s.Value) > 0 {
		decrypted, err := security.Decrypt([]byte(s.Value), passphrase)
		if err != nil {
			return &Secret{}
		}
		val = string(decrypted)
	}

	return &Secret{
		ID:        s.ID,
		Value:     val,
		CreatedAt: s.CreatedAt,
		ExpiresAt: s.ExpiresAt,
	}
}
