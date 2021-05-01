package db

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
}

// NewSecretRepository creates and returns a SecretRepository
// object.
func NewSecretRepository(c Client, opts *SecretRepositoryOptions) *SecretRepository {
	var db Database
	switch c := c.(type) {
	case *mongoClient:
		db = c
	case *redisClient:
		db = c
	}

	return &SecretRepository{
		db:      db,
		options: opts,
	}
}

// SecretRepositoryOptions provides additional options
// for the repository. It contains: Driver.
type SecretRepositoryOptions struct {
	Driver string
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
