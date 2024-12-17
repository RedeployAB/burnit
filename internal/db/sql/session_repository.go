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
	// defaultSessionRepositoryTable is the default table for the SessionRepository.
	defaultSessionRepositoryTable = "sessions"
	// defaultSessionRepositoryTimeout is the default timeout for the SessionRepository.
	defaultSessionRepositoryTimeout = 10 * time.Second
)

// SessionRepository is a SQL implementation of a SessionRepository.
type SessionRepository struct {
	db      *sql.DB
	driver  Driver
	table   string
	queries queries
	timeout time.Duration
}

// SessionRepositoryOptions is the options for the SessionRepository.
type SessionRepositoryOptions struct {
	Table   string
	Timeout time.Duration
}

// SessionRepositoryOption is a function that sets options for the SessionRepository.
type SessionRepositoryOption func(o *SessionRepositoryOptions)

// NewSessionRepository returns a new SessionRepository.
func NewSessionRepository(db *DB, options ...SessionRepositoryOption) (*SessionRepository, error) {
	if db == nil {
		return nil, ErrNilDB
	}

	opts := SessionRepositoryOptions{
		Table:   defaultSessionRepositoryTable,
		Timeout: defaultSessionRepositoryTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	queries, err := createSessionQueries(db.driver, opts.Table)
	if err != nil {
		return nil, err
	}

	r := &SessionRepository{
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
func (r SessionRepository) createTableIfNotExists(ctx context.Context) error {
	var query string
	var args []any

	switch r.driver {
	case DriverPostgres:
		query = `
		CREATE TABLE IF NOT EXISTS %s (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			expires_at TIMESTAMPTZ NOT NULL,
			csrf_token VARCHAR(43),
			csrf_expires_at TIMESTAMPTZ NOT NULL
		)`
		args = append(args, r.table)
	case DriverMSSQL:
		table := firstToUpper(r.table)
		query = `
		IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U')
		CREATE TABLE %s (
			ID VARCHAR(36) NOT NULL PRIMARY KEY,
			ExpiresAt DATETIMEOFFSET NOT NULL,
			CSRFToken VARCHAR(43),
			CSRFExpiresAt DATETIMEOFFSET NOT NULL
		)`
		args = append(args, table, table)
	case DriverSQLite:
		query = `
		CREATE TABLE IF NOT EXISTS %s (
			id TEXT NOT NULL PRIMARY KEY,
			expires_at DATETIME NOT NULL,
			csrf_token TEXT NOT NULL,
			csrf_expires_at DATETIME NOT NULL
		)`
		args = append(args, r.table)
	default:
		return fmt.Errorf("%w: %s", ErrDriverNotSupported, r.driver)
	}

	if _, err := r.db.ExecContext(ctx, fmt.Sprintf(query, args...)); err != nil {
		return err
	}
	return nil
}

// Get a session by its ID.
func (r SessionRepository) Get(ctx context.Context, id string) (db.Session, error) {
	var session db.Session
	if err := r.db.QueryRowContext(ctx, r.queries.selectByID, id).Scan(&session.ID, &session.ExpiresAt, &session.CSRF.Token, &session.CSRF.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Session{}, dberrors.ErrSessionNotFound
		}
		return db.Session{}, err
	}
	return session, nil
}

// Create a session.
func (r SessionRepository) Create(ctx context.Context, session db.Session) (db.Session, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return db.Session{}, err
	}

	if _, err := tx.ExecContext(ctx, r.queries.insert, session.ID, session.ExpiresAt, session.CSRF.Token, session.CSRF.ExpiresAt); err != nil {
		if err := tx.Rollback(); err != nil {
			return db.Session{}, err
		}
		return db.Session{}, err
	}

	if err := tx.QueryRowContext(ctx, r.queries.selectByID, session.ID).Scan(&session.ID, &session.ExpiresAt, &session.CSRF.Token, &session.CSRF.ExpiresAt); err != nil {
		if err := tx.Rollback(); err != nil {
			return db.Session{}, err
		}
		return db.Session{}, err
	}

	if err := tx.Commit(); err != nil {
		return db.Session{}, err
	}

	return session, nil
}

// Delete a session by its ID.
func (r SessionRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, r.queries.delete, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return dberrors.ErrSessionNotFound
	}

	return nil
}

// DeleteExpired deletes all expired sessions.
func (r SessionRepository) DeleteExpired(ctx context.Context) error {
	result, err := r.db.ExecContext(ctx, r.queries.deleteExpired)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return dberrors.ErrSessionsNotDeleted
	}

	return nil
}

// Close the repository and its underlying connections.
func (r SessionRepository) Close() error {
	return r.db.Close()
}

// createSessionQueries creates the queries used by the repository.
func createSessionQueries(driver Driver, table string) (queries, error) {
	var columns, placeholders []string
	var now string
	switch driver {
	case DriverPostgres:
		columns = []string{"id", "expires_at", "csrf_token", "csrf_expires_at"}
		placeholders = []string{"$1", "$2", "$3", "$4"}
		now = "NOW() AT TIME ZONE 'UTC'"
	case DriverMSSQL:
		table = firstToUpper(table)
		columns = []string{"ID", "ExpiresAt", "CSRFToken", "CSRFExpiresAt"}
		placeholders = []string{"@p1", "@p2", "@p3", "@p4"}
		now = "GETUTCDATE()"
	case DriverSQLite:
		columns = []string{"id", "expires_at", "csrf_token", "csrf_expires_at"}
		placeholders = []string{"?1", "?2", "?3", "?4"}
		now = "DATETIME('now')"
	default:
		return queries{}, fmt.Errorf("%w: %s", ErrDriverNotSupported, driver)
	}

	return queries{
		selectByID:    fmt.Sprintf("SELECT %s, %s, %s, %s FROM %s WHERE %s = %s", columns[0], columns[1], columns[2], columns[3], table, columns[0], placeholders[0]),
		insert:        fmt.Sprintf("INSERT INTO %s (%s, %s, %s, %s) VALUES (%s, %s, %s, %s)", table, columns[0], columns[1], columns[2], columns[3], placeholders[0], placeholders[1], placeholders[2], placeholders[3]),
		delete:        fmt.Sprintf("DELETE FROM %s WHERE %s = %s", table, columns[0], placeholders[0]),
		deleteExpired: fmt.Sprintf("DELETE FROM %s WHERE %s < %s", table, columns[1], now),
	}, nil
}
