package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

const (
	// defaultDatabase = "burnit" is the default name for the database.
	defaultDatabase = "burnit"
	// defaultDatabaseFile is the default file for the SQLite database.
	defaultDatabaseFile = "burnit.db"
	// defaultConnectTimeout is the default timeout for the database connection.
	defaultConnectTimeout = 10 * time.Second
)

// Driver is the type for the database driver.
type Driver string

const (
	// DriverPostgres is the PostgreSQL driver.
	DriverPostgres Driver = "pgx"
	// DriverMSSQL is the Microsoft SQL Server driver.
	DriverMSSQL Driver = "sqlserver"
	// DriverSQLite is the SQLite driver.
	DriverSQLite Driver = "sqlite"
)

// Scheme returns the scheme for the database driver, used for
// the data source name.
func (d Driver) Scheme() string {
	switch d {
	case DriverPostgres:
		return "postgres"
	case DriverMSSQL:
		return "sqlserver"
	case DriverSQLite:
		return "sqlite"
	}
	return ""
}

// PostgresSSLMode is the available settings for PostgreSQL SSL/TLS settings.
type PostgresSSLMode string

const (
	// PostgresSSLModeDisable disables encryption for PostgreSQL.
	PostgresSSLModeDisable PostgresSSLMode = "disable"
	// PostgresSSLModePrefer sets encryption to prefer for PostgreSQL.
	PostgresSSLModePrefer PostgresSSLMode = "prefer"
	// PostgresSSLModeRequire sets encryption to require for PostgreSQL.
	PostgresSSLModeRequire PostgresSSLMode = "require"
	// PostgresSSLModeVerifyCA sets encryption to verify-ca for PostgreSQL.
	PostgresSSLModeVerifyCA PostgresSSLMode = "verify-ca"
	// PostgresSSLModeVerifyFull sets encryption to verify-full for PostgreSQL.
	PostgresSSLModeVerifyFull PostgresSSLMode = "verify-full"
)

// MSSQLEncrypt is the available settings for MSSQL encryption.
type MSSQLEncrypt string

const (
	// MSSQLEncryptTrue sets encryption to true for MSSQL.
	MSSQLEncryptTrue MSSQLEncrypt = "true"
	// MSSQLEncryptFalse sets encryption to false for MSSQL.
	MSSQLEncryptFalse MSSQLEncrypt = "false"
	// MSSQLEncryptStrict sets encryption to strict for MSSQL.
	MSSQLEncryptStrict MSSQLEncrypt = "strict"
)

// Row is a result row.
type Row interface {
	Scan(dest ...any) error
}

// Result is the result of a query.
type Result interface {
	RowsAffected() (int64, error)
}

// Client is the interface for the database client.
type Client interface {
	QueryRow(ctx context.Context, query string, args ...any) Row
	Exec(ctx context.Context, query string, args ...any) (Result, error)
	Transaction(ctx context.Context) (Tx, error)
	Driver() Driver
	Close() error
}

// Tx is the interface for the database transaction.
type Tx interface {
	QueryRow(ctx context.Context, query string, args ...any) Row
	Exec(ctx context.Context, query string, args ...any) (Result, error)
	Commit() error
	Rollback() error
}

// client is the database client.
type client struct {
	db     *sql.DB
	driver Driver
}

// ClientOptions contains the options for the client.
type ClientOptions struct {
	Driver                Driver
	DSN                   string
	Address               string
	Database              string
	Username              string
	Password              string
	ConnectTimeout        time.Duration
	MaxOpenConnections    int
	MaxIdleConnections    int
	MaxConnectionLifetime time.Duration
	Postgres              PostgresOptions
	MSSQL                 MSSQLOptions
	SQLite                SQLiteOptions
}

// ClientOption is a function that sets options for the client.
type ClientOption func(o *ClientOptions)

// NewClient creates a new database client.
func NewClient(options ...ClientOption) (*client, error) {
	opts := ClientOptions{
		Database:       defaultDatabase,
		ConnectTimeout: defaultConnectTimeout,
		SQLite: SQLiteOptions{
			File: defaultDatabaseFile,
		},
	}
	for _, option := range options {
		option(&opts)
	}

	driver, err := checkDriver(opts.Driver)
	if err != nil {
		return nil, err
	}

	if len(opts.Postgres.SSLMode) > 0 && !validPostgresSSLMode(opts.Postgres.SSLMode) {
		return nil, fmt.Errorf("invalid postgres sslmode: %s", opts.Postgres.SSLMode)
	}
	if len(opts.MSSQL.Encrypt) > 0 && !validMSSQLEncrypt(opts.MSSQL.Encrypt) {
		return nil, fmt.Errorf("invalid mssql encrypt setting: %s", opts.MSSQL.Encrypt)
	}

	dsn := buildDSN(driver, &opts)

	db, err := sql.Open(string(driver), dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	if opts.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(opts.MaxOpenConnections)
	}
	if opts.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(opts.MaxIdleConnections)
	}
	if opts.MaxConnectionLifetime > 0 {
		db.SetConnMaxLifetime(opts.MaxConnectionLifetime)
	}

	return &client{db: db, driver: driver}, nil
}

// QueryRow executes a query that is expected to return at most one row.
func (c client) QueryRow(ctx context.Context, query string, args ...any) Row {
	return c.db.QueryRowContext(ctx, query, args...)
}

// Exec executes a query without returning any rows.
func (c client) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	return c.db.ExecContext(ctx, query, args...)
}

// Driver returns the database driver.
func (c client) Driver() Driver {
	return c.driver
}

// Close the database and release any open resources.
func (c client) Close() error {
	err := c.db.Close()
	if err == nil || errors.Is(err, sql.ErrConnDone) {
		return nil
	}
	return err
}

// tx is a transaction. It wraps a *sql.Tx.
type tx struct {
	*sql.Tx
}

// QueryRow executes a query that is expected to return at most one row.
func (t tx) QueryRow(ctx context.Context, query string, args ...any) Row {
	return t.Tx.QueryRowContext(ctx, query, args...)
}

// Exec executes a query without returning any rows.
func (t tx) Exec(ctx context.Context, query string, args ...any) (Result, error) {
	return t.Tx.ExecContext(ctx, query, args...)
}

// Commit commits the transaction.
func (t tx) Commit() error {
	return t.Tx.Commit()
}

// Rollback rolls back the transaction.
func (t tx) Rollback() error {
	return t.Tx.Rollback()
}

// Tx starts a new transaction.
func (c client) Transaction(ctx context.Context) (Tx, error) {
	t, err := c.db.BeginTx(ctx, nil)
	return tx{t}, err
}

// PostgresOptions contains the options for PostgreSQL.
type PostgresOptions struct {
	SSLMode PostgresSSLMode
}

// MSSQLOptions contains the options for MSSQL.
type MSSQLOptions struct {
	Encrypt MSSQLEncrypt
}

// SQLiteOptions contains the options for SQLite.
type SQLiteOptions struct {
	File     string
	InMemory bool
}

// checkDriver checks if the database driver is supported and returns the driver
// with the correct name for the imported driver packages.
func checkDriver(driver Driver) (Driver, error) {
	switch driver {
	case DriverPostgres, "postgres":
		return DriverPostgres, nil
	case DriverMSSQL, "mssql":
		return DriverMSSQL, nil
	case DriverSQLite:
		return DriverSQLite, nil
	}
	return "", fmt.Errorf("%w: %s", ErrDriverNotSupported, driver)
}

// buildDSN builds the data source name for the database connection.
func buildDSN(driver Driver, options *ClientOptions) string {
	if options == nil {
		return ""
	}

	if len(options.DSN) > 0 {
		return options.DSN
	}

	if driver == DriverSQLite {
		return databaseFileDSN(options.SQLite.File, options.SQLite.InMemory)
	}

	u := url.URL{
		Scheme: driver.Scheme(),
		Host:   options.Address,
	}

	if len(options.Username) > 0 {
		u.User = url.UserPassword(options.Username, options.Password)
	}

	database := options.Database
	if len(database) == 0 {
		database = defaultDatabase
	}

	switch driver {
	case DriverPostgres:
		u.Path = database
		if len(options.Postgres.SSLMode) > 0 {
			u.RawQuery = "sslmode=" + string(options.Postgres.SSLMode)
		}
	case DriverMSSQL:
		if database == defaultDatabase {
			database = firstToUpper(database)
		}
		u.RawQuery = "database=" + database
		if len(options.MSSQL.Encrypt) > 0 {
			u.RawQuery += "&encrypt=" + string(options.MSSQL.Encrypt)
		}
	}

	return u.String()
}

// databaseFileDSN returns the database file for SQLite.
func databaseFileDSN(file string, inMemory bool) string {
	if inMemory {
		return ":memory:"
	}
	if len(file) > 0 {
		return "file:" + file
	}
	return "file:" + defaultDatabaseFile
}

// firstToUpper returns the string with the first letter in uppercase.
func firstToUpper(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

// validPostgresSSLMode checks if the PostgreSQL SSL mode is valid.
func validPostgresSSLMode(mode PostgresSSLMode) bool {
	switch mode {
	case PostgresSSLModeDisable, PostgresSSLModePrefer, PostgresSSLModeRequire, PostgresSSLModeVerifyCA, PostgresSSLModeVerifyFull:
		return true
	}
	return false
}

// validMSSQLEncrypt checks if the MSSQL encrypt setting is valid.
func validMSSQLEncrypt(encrypt MSSQLEncrypt) bool {
	switch encrypt {
	case MSSQLEncryptTrue, MSSQLEncryptFalse, MSSQLEncryptStrict:
		return true
	}
	return false
}
