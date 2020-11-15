package db

import (
	"github.com/RedeployAB/burnit/common/security"
)

// Repository defined the methods needed for interact
// with a database and collection.
type Repository interface {
	Find(id string) (*Secret, error)
	Insert(s *Secret) (*Secret, error)
	Delete(id string) (int64, error)
	DeleteExpired() (int64, error)
}

// SecretRepository handles interactions with the database
// and collection.
type SecretRepository struct {
	db      Database
	options *SecretRepositoryOptions
	hash    func(s string) string
}

// NewSecretRepository creates and returns a SecretRepository
// object.
func NewSecretRepository(c Client, opts *SecretRepositoryOptions) *SecretRepository {
	var hash func(s string) string
	switch opts.HashMethod {
	case "md5":
		hash = security.ToMD5
	case "bcrypt":
		hash = security.Bcrypt
	}

	var db Database
	switch c.(type) {
	case *mongoClient:
		db = c.(*mongoClient)
	case *redisClient:
		db = c.(*redisClient)
	}

	return &SecretRepository{
		db:      db,
		options: opts,
		hash:    hash,
	}
}

// SecretRepositoryOptions provides additional options
// for the repository. It contains: hashMethod.
type SecretRepositoryOptions struct {
	Driver     string
	HashMethod string
}

// Find queries the collection for an entry by ID.
func (r *SecretRepository) Find(id string) (*Secret, error) {
	s, err := r.db.FindOne(id)
	if err != nil || s == nil {
		return s, err
	}
	return s, nil
}

// Insert handles inserts into the database.
func (r *SecretRepository) Insert(s *Secret) (*Secret, error) {
	if len(s.Passphrase) > 0 {
		s.Passphrase = r.hash(s.Passphrase)
	}

	s, err := r.db.InsertOne(s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Delete removes an entry from the collection by ID.
func (r *SecretRepository) Delete(id string) (int64, error) {
	deleted, err := r.db.DeleteOne(id)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// DeleteExpired deletes all entries that has expiresAt
// less than current time (time of invocation).
func (r *SecretRepository) DeleteExpired() (int64, error) {
	deleted, err := r.db.DeleteMany()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
