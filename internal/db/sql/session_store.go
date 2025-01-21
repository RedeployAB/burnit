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
	// defaultSessionStoreTable is the default table for the SessionStore.
	defaultSessionStoreTable = "sessions"
	// defaultSessionStoreTimeout is the default timeout for the SessionStore.
	defaultSessionStoreTimeout = 10 * time.Second
)

// sessionStore is a SQL implementation of a SessionStore.
type sessionStore struct {
	client  Client
	driver  Driver
	table   string
	queries sessionQueries
	timeout time.Duration
}

// SessionStoreOptions is the options for the SessionStore.
type SessionStoreOptions struct {
	Table   string
	Timeout time.Duration
}

// SessionStoreOption is a function that sets options for the SessionStore.
type SessionStoreOption func(o *SessionStoreOptions)

// NewSessionStore returns a new SessionStore.
func NewSessionStore(client Client, options ...SessionStoreOption) (*sessionStore, error) {
	if client == nil {
		return nil, errors.New("nil client")
	}

	opts := SessionStoreOptions{
		Table:   defaultSessionStoreTable,
		Timeout: defaultSessionStoreTimeout,
	}
	for _, option := range options {
		option(&opts)
	}

	driver := client.Driver()
	queries, err := createSessionQueries(driver, opts.Table)
	if err != nil {
		return nil, err
	}

	s := &sessionStore{
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
func (s sessionStore) createTableIfNotExists(ctx context.Context) error {
	var query string
	var args []any

	switch s.driver {
	case DriverPostgres:
		query = "CREATE TABLE IF NOT EXISTS %s (id UUID PRIMARY KEY DEFAULT gen_random_uuid(), expires_at TIMESTAMPTZ NOT NULL, csrf_token VARCHAR(43), csrf_expires_at TIMESTAMPTZ NOT NULL)"
		args = append(args, s.table)
	case DriverMSSQL:
		table := firstToUpper(s.table)
		query = "IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='%s' and xtype='U') CREATE TABLE %s (ID VARCHAR(36) NOT NULL PRIMARY KEY, ExpiresAt DATETIMEOFFSET NOT NULL, CSRFToken VARCHAR(43), CSRFExpiresAt DATETIMEOFFSET NOT NULL)"
		args = append(args, table, table)
	case DriverSQLite:
		query = "CREATE TABLE IF NOT EXISTS %s (id TEXT NOT NULL PRIMARY KEY, expires_at DATETIME NOT NULL, csrf_token TEXT NOT NULL, csrf_expires_at DATETIME NOT NULL)"
		args = append(args, s.table)
	default:
		return fmt.Errorf("%w: %s", ErrDriverNotSupported, s.driver)
	}

	if _, err := s.client.Exec(ctx, fmt.Sprintf(query, args...)); err != nil {
		return err
	}
	return nil
}

// Get a session by its ID.
func (s sessionStore) Get(ctx context.Context, id string) (db.Session, error) {
	var session db.Session
	if err := s.client.QueryRow(ctx, s.queries.selectByID, id).Scan(&session.ID, &session.ExpiresAt, &session.CSRF.Token, &session.CSRF.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Session{}, dberrors.ErrSessionNotFound
		}
		return db.Session{}, err
	}
	return session, nil
}

// Get a session by its CSRF token.
func (s sessionStore) GetByCSRFToken(ctx context.Context, token string) (db.Session, error) {
	var session db.Session
	if err := s.client.QueryRow(ctx, s.queries.selectByCSRFToken, token).Scan(&session.ID, &session.ExpiresAt, &session.CSRF.Token, &session.CSRF.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Session{}, dberrors.ErrSessionNotFound
		}
		return db.Session{}, err
	}
	return session, nil
}

// Upsert a session. Create the session if it does not exist, otherwise
// update it.
func (s sessionStore) Upsert(ctx context.Context, session db.Session) (db.Session, error) {
	tx, err := s.client.Transaction(ctx)
	if err != nil {
		return db.Session{}, nil
	}

	if _, err := tx.Exec(ctx, s.queries.upsert, session.ID, session.ExpiresAt, session.CSRF.Token, session.CSRF.ExpiresAt); err != nil {
		if err := tx.Rollback(); err != nil {
			return db.Session{}, err
		}
		return db.Session{}, err
	}

	if err := tx.QueryRow(ctx, s.queries.selectByID, session.ID).Scan(&session.ID, &session.ExpiresAt, &session.CSRF.Token, &session.CSRF.ExpiresAt); err != nil {
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
func (s sessionStore) Delete(ctx context.Context, id string) error {
	result, err := s.client.Exec(ctx, s.queries.delete, id)
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

// DeleteByCSRFToken deletes a session by its CSRF token.
func (s sessionStore) DeleteByCSRFToken(ctx context.Context, token string) error {
	result, err := s.client.Exec(ctx, s.queries.deleteByCSRFToken, token)
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
func (s sessionStore) DeleteExpired(ctx context.Context) error {
	result, err := s.client.Exec(ctx, s.queries.deleteExpired)
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

// Close the store and its underlying connections.
func (s sessionStore) Close() error {
	return s.client.Close()
}

// sessionQueries contains queries used by the Store.
type sessionQueries struct {
	selectByID        string
	selectByCSRFToken string
	upsert            string
	delete            string
	deleteByCSRFToken string
	deleteExpired     string
}

// createSessionQueries creates the queries used by the Store.
func createSessionQueries(driver Driver, table string) (sessionQueries, error) {
	var columns, placeholders []string
	var now, upsert string
	switch driver {
	case DriverPostgres:
		columns = []string{"id", "expires_at", "csrf_token", "csrf_expires_at"}
		placeholders = []string{"$1", "$2", "$3", "$4"}
		now = "NOW() AT TIME ZONE 'UTC'"
		upsert = "INSERT INTO %s (id, expires_at, csrf_token, csrf_expires_at) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET expires_at = EXCLUDED.expires_at, csrf_token = EXCLUDED.csrf_token, csrf_expires_at = EXCLUDED.csrf_expires_at"
	case DriverMSSQL:
		table = firstToUpper(table)
		columns = []string{"ID", "ExpiresAt", "CSRFToken", "CSRFExpiresAt"}
		placeholders = []string{"@p1", "@p2", "@p3", "@p4"}
		now = "GETUTCDATE()"
		upsert = "MERGE INTO %s AS target USING (VALUES (@p1, @p2, @p3, @p4)) AS source (ID, ExpiresAt, CSRFToken, CSRFExpiresAt) ON target.ID = source.ID WHEN MATCHED THEN UPDATE SET target.ExpiresAt = source.ExpiresAt, target.CSRFToken = source.CSRFToken, target.CSRFExpiresAt = source.CSRFExpiresAt WHEN NOT MATCHED THEN INSERT (ID, ExpiresAt, CSRFToken, CSRFExpiresAt) VALUES (source.ID, source.ExpiresAt, source.CSRFToken, source.CSRFExpiresAt);"
	case DriverSQLite:
		columns = []string{"id", "expires_at", "csrf_token", "csrf_expires_at"}
		placeholders = []string{"?1", "?2", "?3", "?4"}
		now = "DATETIME('now')"
		upsert = "INSERT INTO %s (id, expires_at, csrf_token, csrf_expires_at) VALUES (?1, ?2, ?3, ?4) ON CONFLICT(id) DO UPDATE SET expires_at = excluded.expires_at, csrf_token = excluded.csrf_token, csrf_expires_at = excluded.csrf_expires_at"
	default:
		return sessionQueries{}, fmt.Errorf("%w: %s", ErrDriverNotSupported, driver)
	}

	return sessionQueries{
		selectByID:        fmt.Sprintf("SELECT %s, %s, %s, %s FROM %s WHERE %s = %s", columns[0], columns[1], columns[2], columns[3], table, columns[0], placeholders[0]),
		selectByCSRFToken: fmt.Sprintf("SELECT %s, %s, %s, %s FROM %s WHERE %s = %s", columns[0], columns[1], columns[2], columns[3], table, columns[2], placeholders[0]),
		upsert:            fmt.Sprintf(upsert, table),
		delete:            fmt.Sprintf("DELETE FROM %s WHERE %s = %s", table, columns[0], placeholders[0]),
		deleteByCSRFToken: fmt.Sprintf("DELETE FROM %s WHERE %s = %s", table, columns[2], placeholders[0]),
		deleteExpired:     fmt.Sprintf("DELETE FROM %s WHERE %s < %s", table, columns[1], now),
	}, nil
}
