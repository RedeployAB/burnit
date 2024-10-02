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
	var enableTLS bool
	if config.EnableTLS != nil {
		enableTLS = *config.EnableTLS
	}

	var repo db.SecretRepository
	driver, err := databaseDriver(config)
	if err != nil {
		return nil, err
	}

	switch driver {
	case databaseDriverMongo:
		repo, err = setupMongoSecretRepository(config, enableTLS)
	case databaseDriverPostgres:
		repo, err = setupSQLRepository(config, driver, enableTLS)
	default:
		return nil, fmt.Errorf("unsupported database driver")
	}

	if err != nil {
		return nil, err
	}

	return repo, nil
}

// setupMongoSecretRepository sets up the MongoDB secret repository.
func setupMongoSecretRepository(config *Database, enableTLS bool) (db.SecretRepository, error) {
	client, err := mongo.NewClient(func(o *mongo.ClientOptions) {
		o.URI = config.URI
		o.Hosts = []string{config.Address}
		o.Username = config.Username
		o.Password = config.Password
		o.ConnectTimeout = config.ConnectTimeout
		o.EnableTLS = enableTLS
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
func setupSQLRepository(config *Database, driver string, enableTLS bool) (db.SecretRepository, error) {
	db, err := sql.Open(sql.Driver(driver), config.URI)
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
	case databaseDriverMongo, databaseDriverPostgres:
		return true
	default:
		return false
	}
}

var (
	databaseDriverMongo    = "mongodb"
	databaseDriverPostgres = "postgres"

	databasePorts = map[string]string{
		databaseDriverMongo:    "27017",
		databaseDriverPostgres: "5432",
	}
)
