package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/RedeployAB/burnit/db"
)

const (
	// defaultSecretRepositoryTable is the default table for the SecretRepository.
	defaultSecretRepositoryTable = "secrets"
	// defaultSecretRepositoryTimeout is the default timeout for the SecretRepository.
	defaultSecretRepositoryTimeout = 10 * time.Second
)

// SecretRepository is a PostgreSQL implementation of a SecretRepository.
type SecretRepository struct {
	db      *sql.DB
	driver  Driver
	table   string
	queries queries
	timeout time.Duration
}

// SecretRepositoryOptions is the options for the SecretRepository.
type SecretRepositoryOptions struct {
	Table   string
	Timeout time.Duration
}

// SecretRepositoryOption is a function that sets options for the SecretRepository.
type SecretRepositoryOption func(o *SecretRepositoryOptions)

// NewSecretRepository returns a new SecretRepository.
func NewSecretRepository(db *DB, options ...SecretRepositoryOption) (*SecretRepository, error) {
	if db == nil {
		return nil, ErrNilDB
	}

	opts := SecretRepositoryOptions{
		Table:   defaultSecretRepositoryTable,
		Timeout: defaultSecretRepositoryTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	queries, err := createQueries(db.driver, opts.Table)
	if err != nil {
		return nil, err
	}

	r := &SecretRepository{
		db:      db.DB,
		driver:  db.driver,
		table:   opts.Table,
		queries: queries,
		timeout: opts.Timeout,
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	if err := r.createTableIfNotExists(ctx); err != nil {
		return nil, err
	}

	return r, nil
}

// createTableIfNotExists creates the table if it does not exist.
func (r SecretRepository) createTableIfNotExists(ctx context.Context) error {
	var query string

	var args []any
	switch r.driver {
	case DriverPostgres:
		query = `
		CREATE TABLE IF NOT EXISTS %s (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			value TEXT NOT NULL,
			expires_at TIMESTAMPTZ NOT NULL
		)`
		args = append(args, r.table)
	case DriverMSSQL:
		query = `
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U')
		CREATE TABLE %s (
			id VARCHAR(36) NOT NULL PRIMARY KEY,
			value NVARCHAR(MAX) NOT NULL,
			expires_at DATETIMEOFFSET NOT NULL
		)`
		args = append(args, r.table, r.table)
	case DriverSQLite:
		query = `
		CREATE TABLE IF NOT EXISTS %s (
		    id TEXT NOT NULL PRIMARY KEY,
			value TEXT NOT NULL,
			expires_at DATETIME NOT NULL
		)`
		args = append(args, r.table)
	default:
		return fmt.Errorf("unsupported driver: %s", r.driver)
	}

	if _, err := r.db.ExecContext(ctx, fmt.Sprintf(query, args...)); err != nil {
		return err
	}
	return nil
}

// Get a secret by its ID.
func (r SecretRepository) Get(ctx context.Context, id string) (db.Secret, error) {
	var secret db.Secret
	if err := r.db.QueryRowContext(ctx, r.queries.selectByID, id).Scan(&secret.ID, &secret.Value, &secret.ExpiresAt); err != nil {
		return db.Secret{}, err
	}
	return secret, nil
}

// Create a secret.
func (r SecretRepository) Create(ctx context.Context, secret db.Secret) (db.Secret, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return db.Secret{}, err
	}

	if _, err := tx.ExecContext(ctx, r.queries.insert, secret.ID, secret.Value, secret.ExpiresAt); err != nil {
		if err := tx.Rollback(); err != nil {
			return db.Secret{}, err
		}
		return db.Secret{}, err
	}

	if err := tx.QueryRowContext(ctx, r.queries.selectByID, secret.ID).Scan(&secret.ID, &secret.Value, &secret.ExpiresAt); err != nil {
		if err := tx.Rollback(); err != nil {
			return db.Secret{}, err
		}
		return db.Secret{}, err
	}

	if err := tx.Commit(); err != nil {
		return db.Secret{}, err
	}

	return secret, nil
}

// Delete a secret by its ID.
func (r SecretRepository) Delete(ctx context.Context, id string) error {
	if _, err := r.db.ExecContext(ctx, r.queries.delete, id); err != nil {
		return err
	}
	return nil
}

// DeleteExpired deletes all expired secrets.
func (r SecretRepository) DeleteExpired(ctx context.Context) error {
	if _, err := r.db.ExecContext(ctx, r.queries.deleteExpired); err != nil {
		return err
	}
	return nil
}

// Close the repository and its underlying connections.
func (r SecretRepository) Close() error {
	return r.db.Close()
}

// querie contains queries used by the repository.
type queries struct {
	selectByID    string
	insert        string
	delete        string
	deleteExpired string
}

// createQueries creates the queries used by the repository.
func createQueries(driver Driver, table string) (queries, error) {
	var placeholders []string
	var now string
	switch driver {
	case DriverPostgres:
		placeholders = []string{"$1", "$2", "$3"}
		now = "NOW() AT TIME ZONE 'UTC'"
	case DriverMSSQL:
		placeholders = []string{"@p1", "@p2", "@p3"}
		now = "GETUTCDATE()"
	case DriverSQLite:
		placeholders = []string{"?1", "?2", "?3"}
		now = "DATETIME('now')"
	default:
		return queries{}, fmt.Errorf("unsupported driver: %s", driver)
	}

	return queries{
		selectByID:    fmt.Sprintf("SELECT id, value, expires_at FROM %s WHERE id = %s", table, placeholders[0]),
		insert:        fmt.Sprintf("INSERT INTO %s (id, value, expires_at) VALUES (%s, %s, %s)", table, placeholders[0], placeholders[1], placeholders[2]),
		delete:        fmt.Sprintf("DELETE FROM %s WHERE id = %s", table, placeholders[0]),
		deleteExpired: fmt.Sprintf("DELETE FROM %s WHERE expires_at < %s", table, now),
	}, nil
}
