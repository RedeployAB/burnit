package sql

import (
	"database/sql"
	"fmt"
	"net/url"
)

const (
	// defaultDatabase = "burnit" is the default name for the database.
	defaultDatabase = "burnit"
	// defaultDatabaseFile is the default file for the SQLite database.
	defaultDatabaseFile = "burnit.db"
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

// TLSMode is the available settings for SQL encryption and SSL/TLS
// settings accross drivers. All modes are not supported by all drivers.
type TLSMode string

const (
	// TLSModeDisable disables encryption for PostgreSQL.
	TLSModeDisable TLSMode = "disable"
	// TLSModePrefer sets encryption to prefer for PostgreSQL.
	TLSModePrefer TLSMode = "prefer"
	// TLSModeRequire sets encryption to require for PostgreSQL.
	TLSModeRequire TLSMode = "require"
	// TLSModeVerifyCA sets encryption to verify-ca for PostgreSQL.
	TLSModeVerifyCA TLSMode = "verify-ca"
	// TLSModeVerifyFull sets encryption to verify-full for PostgreSQL.
	TLSModeVerifyFull TLSMode = "verify-full"
	// TLSModeTrue sets encryption to true for MSSQL.
	TLSModeTrue TLSMode = "true"
	// TLSModeFalse sets encryption to false for MSSQL.
	TLSModeFalse TLSMode = "false"
	// TLSModeStrict sets encryption to strict for MSSQL.
	TLSModeStrict TLSMode = "strict"
)

// DB is a database handle representing a pool of zero or more underlying connections.
type DB struct {
	*sql.DB
	driver Driver
}

// Options contains the options for the database.
type Options struct {
	Driver   Driver
	DSN      string
	Address  string
	Database string
	Username string
	Password string
	TLSMode  TLSMode
	File     string
	InMemory bool
}

// Option is a function that sets options for the database.
type Option func(o *Options)

// Open a database specified by its database driver and data source name.
func Open(options ...Option) (*DB, error) {
	opts := Options{
		Database: defaultDatabase,
		File:     defaultDatabaseFile,
	}
	for _, option := range options {
		option(&opts)
	}

	driver, err := checkDriver(opts.Driver)
	if err != nil {
		return nil, err
	}

	var dsn string
	if len(opts.DSN) > 0 {
		dsn = opts.DSN
	} else {
		if driver != DriverSQLite {
			dsn = buildDSN(driver, &opts)
		} else {
			dsn = databaseFile(opts.File, opts.InMemory)
		}
	}

	db, err := sql.Open(string(driver), dsn)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db, driver: driver}, nil
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
	return "", fmt.Errorf("unsupported database driver: %s", driver)
}

// buildDSN builds the data source name for the database connection.
func buildDSN(driver Driver, options *Options) string {
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
		if len(options.TLSMode) > 0 {
			u.RawQuery = "sslmode=" + string(options.TLSMode)
		}
	case DriverMSSQL:
		u.RawQuery = "database=" + database
		if len(options.TLSMode) > 0 {
			u.RawQuery += "&encrypt=" + string(options.TLSMode)
		}
	}

	return u.String()
}

// databaseFile returns the database file for SQLite.
func databaseFile(file string, inMemory bool) string {
	if inMemory {
		return ":memory:"
	}
	if len(file) > 0 {
		return "file:" + file
	}
	return "file:" + defaultDatabaseFile
}
