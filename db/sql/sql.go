package sql

import (
	"database/sql"
)

// Driver is the type for the database driver.
type Driver string

const (
	// DriverPostgres is the PostgreSQL driver.
	DriverPostgres Driver = "pgx"
	// DriverMSSQL is the Microsoft SQL Server driver.
	DriverMSSQL Driver = "sqlserver"
	// DriverMySQL is the MySQL driver.
	DriverMySQL Driver = "mysql"
	// DriverMariaDB is the MariaDB driver.
	DriverMariaDB Driver = "mariadb"
	// DriverSQLite is the SQLite driver.
	DriverSQLite Driver = "sqlite3"
)

// DB is a database handle representing a pool of zero or more underlying connections.
type DB struct {
	*sql.DB
	driver Driver
}

// Open a database specified by its database driver and data source name.
func Open(driver Driver, dsn string) (*DB, error) {
	if driver == "postgres" {
		driver = DriverPostgres
	}

	db, err := sql.Open(string(driver), dsn)
	if err != nil {
		return nil, err
	}

	return &DB{DB: db, driver: driver}, nil
}
