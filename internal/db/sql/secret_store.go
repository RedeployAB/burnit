package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/RedeployAB/burnit/internal/db"
	dberrors "github.com/RedeployAB/burnit/internal/db/errors"
)

const (
	// defaultSecretStoreTable is the default table for the SecretStore.
	defaultSecretStoreTable = "secrets"
	// defaultSecretStoreTimeout is the default timeout for the SecretStore.
	defaultSecretStoreTimeout = 10 * time.Second
)

// secretStore is a SQL implementation of a SecretStore.
type secretStore struct {
	client  Client
	driver  Driver
	table   string
	queries secretQueries
	timeout time.Duration
}

// SecretStoreOptions is the options for the SecretStore.
type SecretStoreOptions struct {
	Table   string
	Timeout time.Duration
}

// SecretStoreOption is a function that sets options for the SecreStore.
type SecretStoreOption func(o *SecretStoreOptions)

// NewSecretStore returns a new SecretStore.
func NewSecretStore(client Client, options ...SecretStoreOption) (*secretStore, error) {
	if client == nil {
		return nil, ErrNilClient
	}

	opts := SecretStoreOptions{
		Table:   defaultSecretStoreTable,
		Timeout: defaultSecretStoreTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	driver := client.Driver()
	queries, err := createSecretQueries(driver, opts.Table)
	if err != nil {
		return nil, err
	}

	s := &secretStore{
		client:  client,
		driver:  driver,
		table:   opts.Table,
		queries: queries,
		timeout: opts.Timeout,
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	if err := s.createTableIfNotExists(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

// createTableIfNotExists creates the table if it does not exist.
func (s secretStore) createTableIfNotExists(ctx context.Context) error {
	var query string
	var args []any

	switch s.driver {
	case DriverPostgres:
		query = `
		CREATE TABLE IF NOT EXISTS %s (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			value TEXT NOT NULL,
			expires_at TIMESTAMPTZ NOT NULL
		)`
		args = append(args, s.table)
	case DriverMSSQL:
		table := firstToUpper(s.table)
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
		args = append(args, s.table)
	default:
		return fmt.Errorf("%w: %s", ErrDriverNotSupported, s.driver)
	}

	if _, err := s.client.Exec(ctx, fmt.Sprintf(query, args...)); err != nil {
		return err
	}
	return nil
}

// Get a secret by its ID.
func (s secretStore) Get(ctx context.Context, id string) (db.Secret, error) {
	var secret db.Secret
	if err := s.client.Query(ctx, s.queries.selectByID, id).Scan(&secret.ID, &secret.Value, &secret.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Secret{}, dberrors.ErrSecretNotFound
		}
		return db.Secret{}, err
	}
	return secret, nil
}

// Create a secret.
func (s secretStore) Create(ctx context.Context, secret db.Secret) (db.Secret, error) {
	tx, err := s.client.Transaction(ctx)
	if err != nil {
		return db.Secret{}, err
	}

	if _, err := tx.Exec(ctx, s.queries.insert, secret.ID, secret.Value, secret.ExpiresAt); err != nil {
		if err := tx.Rollback(); err != nil {
			return db.Secret{}, err
		}
		return db.Secret{}, err
	}

	if err := tx.Query(ctx, s.queries.selectByID, secret.ID).Scan(&secret.ID, &secret.Value, &secret.ExpiresAt); err != nil {
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
func (s secretStore) Delete(ctx context.Context, id string) error {
	result, err := s.client.Exec(ctx, s.queries.delete, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return dberrors.ErrSecretNotFound
	}

	return nil
}

// DeleteExpired deletes all expired secrets.
func (s secretStore) DeleteExpired(ctx context.Context) error {
	result, err := s.client.Exec(ctx, s.queries.deleteExpired)
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

// Close the store and its underlying connections.
func (s secretStore) Close() error {
	return s.client.Close()
}

// secretQueries contains queries used by the store.
type secretQueries struct {
	selectByID    string
	insert        string
	delete        string
	deleteExpired string
}

// createSecretQueries creates the queries used by the store.
func createSecretQueries(driver Driver, table string) (secretQueries, error) {
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
		return secretQueries{}, fmt.Errorf("%w: %s", ErrDriverNotSupported, driver)
	}

	return secretQueries{
		selectByID:    fmt.Sprintf("SELECT %s, %s, %s FROM %s WHERE %s = %s", columns[0], columns[1], columns[2], table, columns[0], placeholders[0]),
		insert:        fmt.Sprintf("INSERT INTO %s (%s, %s, %s) VALUES (%s, %s, %s)", table, columns[0], columns[1], columns[2], placeholders[0], placeholders[1], placeholders[2]),
		delete:        fmt.Sprintf("DELETE FROM %s WHERE %s = %s", table, columns[0], placeholders[0]),
		deleteExpired: fmt.Sprintf("DELETE FROM %s WHERE %s < %s", table, columns[2], now),
	}, nil
}
