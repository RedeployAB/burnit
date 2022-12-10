package db

import (
	"context"
	"time"
)

// SecretRepository handles interactions with the database
// and collection.
type SecretRepository struct {
	client  Client
	timeout time.Duration
}

// SecretRepositoryOptions provides additional options
// for the repository. It contains: Driver.
type SecretRepositoryOptions struct {
	Timeout time.Duration
}

// NewSecretRepository creates and returns a SecretRepository
// object.
func NewSecretRepository(c Client, opts *SecretRepositoryOptions) *SecretRepository {
	if opts == nil {
		opts = &SecretRepositoryOptions{}
	}

	if opts.Timeout == 0 {
		opts.Timeout = time.Second * 5
	}

	return &SecretRepository{
		client:  c,
		timeout: opts.Timeout,
	}
}

// Client returns the underlying client.
func (r *SecretRepository) Client() Client {
	return r.client
}

// Get a secret by ID.
func (r *SecretRepository) Get(id string) (*Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	s, err := r.client.Find(ctx, id)
	if err != nil || s == nil {
		return s, err
	}
	return s, nil
}

// Create a secret.
func (r *SecretRepository) Create(s *Secret) (*Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	s, err := r.client.Insert(ctx, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Delete a secret.
func (r *SecretRepository) Delete(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	deleted, err := r.client.Delete(ctx, id)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

// DeleteExpired deletes all expired secrets.
func (r *SecretRepository) DeleteExpired() (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	deleted, err := r.client.DeleteMany(ctx)
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
