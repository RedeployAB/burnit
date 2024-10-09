package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/RedeployAB/burnit/db"
	dberrors "github.com/RedeployAB/burnit/db/errors"
)

const (
	// defaultSecretRepositoryTable is the default table for the SecretRepository.
	defaultSecretRepositoryTable = "secrets"
	// defaultSecretRepositoryTimeout is the default timeout for the SecretRepository.
	defaultSecretRepositoryTimeout = 10 * time.Second
)

// SecretRepository is a SQL implementation of a SecretRepository.
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
		table := firstToUpper(r.table)
		query = `
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U')
		CREATE TABLE %s (
			ID VARCHAR(36) NOT NULL PRIMARY KEY,
			Value NVARCHAR(MAX) NOT NULL,
			ExpiresAt DATETIMEOFFSET NOT NULL
		)`
		args = append(args, table, table)
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
		if errors.Is(err, sql.ErrNoRows) {
			return db.Secret{}, dberrors.ErrSecretNotFound
		}
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
	_, err := r.Get(ctx, id)
	if err != nil {
		return err
	}

	result, err := r.db.ExecContext(ctx, r.queries.delete, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return dberrors.ErrSecretsNotDeleted
	}

	return nil
}

// DeleteExpired deletes all expired secrets.
func (r SecretRepository) DeleteExpired(ctx context.Context) error {
	result, err := r.db.ExecContext(ctx, r.queries.deleteExpired)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return dberrors.ErrSecretsNotDeleted
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
	var columns, placeholders []string
	var now string
	switch driver {
	case DriverPostgres:
		columns = []string{"id", "value", "expires_at"}
		placeholders = []string{"$1", "$2", "$3"}
		now = "NOW() AT TIME ZONE 'UTC'"
	case DriverMSSQL:
		table = firstToUpper(table)
		columns = []string{"ID", "Value", "ExpiresAt"}
		placeholders = []string{"@p1", "@p2", "@p3"}
		now = "GETUTCDATE()"
	case DriverSQLite:
		columns = []string{"id", "value", "expires_at"}
		placeholders = []string{"?1", "?2", "?3"}
		now = "DATETIME('now')"
	default:
		return queries{}, fmt.Errorf("unsupported driver: %s", driver)
	}

	return queries{
		selectByID:    fmt.Sprintf("SELECT %s, %s, %s FROM %s WHERE %s = %s", columns[0], columns[1], columns[2], table, columns[0], placeholders[0]),
		insert:        fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES (%s, %s, %s)", table, columns[0], columns[1], columns[2], placeholders[0], placeholders[1], placeholders[2]),
		delete:        fmt.Sprintf("DELETE FROM %s WHERE %s = %s", table, columns[0], placeholders[0]),
		deleteExpired: fmt.Sprintf("DELETE FROM %s WHERE %s < %s", table, columns[2], now),
	}, nil
}
