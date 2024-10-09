package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/RedeployAB/burnit/db"
	"github.com/RedeployAB/burnit/db/mongo"
	"github.com/RedeployAB/burnit/db/sql"
	"github.com/RedeployAB/burnit/secret"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/microsoft/go-mssqldb"
	_ "modernc.org/sqlite"
)

// services contains the configured and setup services.
type services struct {
	Secrets secret.Service
}

// SetupServices configures and sets up the services.
func SetupServices(config Services) (*services, error) {
	repo, err := setupSecretRepository(&config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret repository: %w", err)
	}

	secrets, err := secret.NewService(
		repo,
		secret.WithEncryptionKey(config.Secret.EncryptionKey),
		secret.WithTimeout(config.Secret.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret service: %w", err)
	}

	return &services{
		Secrets: secrets,
	}, nil
}

// setupSecretRepository sets up the secret repository.
func setupSecretRepository(config *Database) (db.SecretRepository, error) {
	var repo db.SecretRepository
	driver, err := databaseDriver(config)
	if err != nil {
		return nil, err
	}

	switch driver {
	case databaseDriverMongo:
		repo, err = setupMongoSecretRepository(config)
	case databaseDriverPostgres, databaseDriverMSSQL, databaseDriverSQLite:
		repo, err = setupSQLRepository(config, driver)
	default:
		return nil, fmt.Errorf("unsupported database driver")
	}

	if err != nil {
		return nil, err
	}

	return repo, nil
}

// setupMongoSecretRepository sets up the MongoDB secret repository.
func setupMongoSecretRepository(config *Database) (db.SecretRepository, error) {
	client, err := mongo.NewClient(func(o *mongo.ClientOptions) {
		o.URI = config.URI
		o.Hosts = []string{config.Address}
		o.Username = config.Username
		o.Password = config.Password
		o.ConnectTimeout = config.ConnectTimeout
		o.EnableTLS = parseBool(config.TLS)
	})
	if err != nil {
		return nil, err
	}

	return mongo.NewSecretRepository(client, func(o *mongo.SecretRepositoryOptions) {
		o.Database = config.Database
		o.Timeout = config.Timeout
	})
}

// setupSQLRepository sets up the SQL secret repository.
func setupSQLRepository(config *Database, driver string) (db.SecretRepository, error) {
	drv := sql.Driver(driver)
	if drv == sql.DriverMSSQL && config.Database == defaultDatabaseName {
		config.Database = strings.ToUpper(config.Database[:1]) + config.Database[1:]
	}

	var inMemory bool
	if config.InMemory != nil {
		inMemory = *config.InMemory
	}

	db, err := sql.Open(func(o *sql.Options) {
		o.Driver = drv
		o.DSN = config.URI
		o.Address = config.Address
		o.Database = config.Database
		o.Username = config.Username
		o.Password = config.Password
		o.File = config.File
		o.InMemory = inMemory
		o.TLSMode = tlsMode(driver, config.TLS)
	})
	if err != nil {
		return nil, err
	}

	return sql.NewSecretRepository(db, func(o *sql.SecretRepositoryOptions) {
		o.Timeout = config.Timeout
	})
}

// databaseDriver returns the database driver.
func databaseDriver(db *Database) (string, error) {
	var driver string

	if len(db.Driver) > 0 && supportedDBDriver(db.Driver) {
		return db.Driver, nil
	}

	if len(db.URI) > 0 {
		driver = dbDriverFromURI(db.URI)
		if len(driver) > 0 && supportedDBDriver(driver) {
			return driver, nil
		}
	}

	if len(db.Address) > 0 {
		driver = dbDriverFromAddress(db.Address)
		if len(driver) > 0 && supportedDBDriver(driver) {
			return driver, nil
		}
	}

	if len(db.File) > 0 || *db.InMemory {
		return databaseDriverSQLite, nil
	}

	if len(driver) == 0 {
		return "", fmt.Errorf("could not determine database driver")
	}

	return driver, nil
}

// dbDriverFromURI returns the database driver from the URI.
func dbDriverFromURI(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return ""
	}
	return u.Scheme
}

// dbDriverFromAddress returns the database driver from the address.
func dbDriverFromAddress(addr string) string {
	if !strings.Contains(addr, ":") {
		return ""
	}
	port := addr[strings.Index(addr, ":")+1:]
	return databasePorts[port]
}

// supportedDBDriver returns true if the driver is supported.
func supportedDBDriver(driver string) bool {
	switch driver {
	case databaseDriverMongo, databaseDriverPostgres, databaseDriverMSSQL, databaseDriverSQLite:
		return true
	default:
		return false
	}
}

var (
	databaseDriverMongo    = "mongodb"
	databaseDriverPostgres = "postgres"
	databaseDriverMSSQL    = "sqlserver"
	databaseDriverSQLite   = "sqlite"

	databasePorts = map[string]string{
		databaseDriverMongo:    "27017",
		databaseDriverPostgres: "5432",
		databaseDriverMSSQL:    "1433",
	}
)

// tlsMode returns the TLS mode based on the driver and TLS.
func tlsMode(driver string, tls string) sql.TLSMode {
	switch driver {
	case databaseDriverMSSQL:
		if tls == "require" {
			return sql.TLSModeTrue
		}
		if tls == "disable" {
			return sql.TLSModeFalse
		}
	}
	return sql.TLSMode(tls)
}

// parseBool parses a string to a boolean.
func parseBool(b string) bool {
	return b == "true" || b == "1" || b == "require"
}
