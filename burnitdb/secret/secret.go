package secret

import (
	"encoding/json"
	"io"
	"time"

	"github.com/RedeployAB/burnit/burnitdb/db"
)

// Secret is to be used as the middle-hand between
// incoming JSON payload and the data model.
type Secret struct {
	ID         string
	Secret     string
	Passphrase string
	CreatedAt  time.Time
	ExpiresAt  time.Time
	TTL        int
}

// Service provides handling operations for secrets.
type Service interface {
	Get(id string) (*Secret, error)
	Create(secret *Secret) (*Secret, error)
	Delete(id string) (int64, error)
	DeleteExpired() (int64, error)
}

// service is the layer that handles translation
// of incoming JSON payload to Secret, to Secret
// compatible with database.
type service struct {
	secrets db.Repository
}

// Gets a secret from the repository by ID.
func (s service) Get(id string) (*Secret, error) {
	model, err := s.secrets.Find(id)
	if err != nil || model == nil {
		return nil, err
	}
	return toSecret(model), nil
}

// Create a secret in the repository.
func (s service) Create(secret *Secret) (*Secret, error) {
	model, err := s.secrets.Insert(toModel(secret))
	if err != nil {
		return nil, err
	}
	return toSecret(model), nil
}

// Delete a secret from the repository.
func (s service) Delete(id string) (int64, error) {
	deleted, err := s.secrets.Delete(id)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// Delete expired deletes all entries that has expiresAt
// less than current  time (time of invocation).
func (s service) DeleteExpired() (int64, error) {
	deleted, err := s.secrets.DeleteExpired()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// NewService creates a new service for handling secrets.
func NewService(secrets db.Repository) Service {
	return &service{secrets: secrets}
}

// NewFromJSON creates an incoming JSON payload
// and creates a Secret from it.
func NewFromJSON(b io.ReadCloser) (*Secret, error) {
	var secret *Secret
	if err := json.NewDecoder(b).Decode(&secret); err != nil {
		return nil, err
	}

	if secret.TTL != 0 {
		secret.ExpiresAt = time.Now().Add(time.Minute * time.Duration(secret.TTL))
	} else {
		secret.ExpiresAt = time.Now().AddDate(0, 0, 7)
	}

	return secret, nil
}

// toModel transforms a Secret to the
// data model variant of Secret.
func toModel(secret *Secret) *db.Secret {
	var createdAt, expiresAt time.Time
	if secret.CreatedAt.IsZero() {
		createdAt = time.Now()
	} else {
		createdAt = secret.CreatedAt
	}

	if secret.ExpiresAt.IsZero() {
		expiresAt = time.Now().Add(time.Minute * time.Duration(10080))
	} else {
		expiresAt = secret.ExpiresAt
	}

	return &db.Secret{
		Secret:     secret.Secret,
		CreatedAt:  createdAt,
		ExpiresAt:  expiresAt,
		Passphrase: secret.Passphrase,
	}
}

// toSecret transforms the data model variant
// of Secret to a Secret.
func toSecret(s *db.Secret) *Secret {
	return &Secret{
		ID:         s.ID,
		Secret:     s.Secret,
		Passphrase: s.Passphrase,
		CreatedAt:  s.CreatedAt,
		ExpiresAt:  s.ExpiresAt,
	}
}
