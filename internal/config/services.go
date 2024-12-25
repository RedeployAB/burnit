package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/RedeployAB/burnit/internal/db"
	"github.com/RedeployAB/burnit/internal/db/mongo"
	"github.com/RedeployAB/burnit/internal/db/redis"
	"github.com/RedeployAB/burnit/internal/db/sql"
	"github.com/RedeployAB/burnit/internal/secret"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/microsoft/go-mssqldb"
	_ "modernc.org/sqlite"
)

// services contains the configured and setup services.
type services struct {
	Secrets secret.Service
}

// dbClient contains configured db client.
type dbClient struct {
	mongo mongo.Client
	sql   *sql.DB
	redis redis.Client
}

// SetupServices configures and sets up the services.
func SetupServices(config Services) (*services, error) {
	dbClient, err := setupDBClient(&config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to setup database client: %w", err)
	}

	store, err := setupSecretStore(dbClient, &config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret store: %w", err)
	}

	secrets, err := secret.NewService(
		store,
		secret.WithValueMaxCharacters(config.Secret.ValueMaxCharacters),
		secret.WithTimeout(config.Secret.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret service: %w", err)
	}

	return &services{
		Secrets: secrets,
	}, nil
}

// setupDBClient sets up the db client.
func setupDBClient(config *Database) (*dbClient, error) {
	var err error
	config.Driver, err = databaseDriver(config)
	if err != nil {
		return nil, err
	}

	var client dbClient
	switch config.Driver {
	case databaseDriverMongo:
		client.mongo, err = setupMongoClient(config)
	case databaseDriverPostgres, databaseDriverMSSQL, databaseDriverSQLite:
		client.sql, err = setupSQLClient(config)
	case databaseDriverRedis:
		client.redis, err = setupRedisClient(config)
	}
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// setupMongoClient sets up the mongo client.
func setupMongoClient(config *Database) (mongo.Client, error) {
	var enableTLS bool
	if config.Redis.EnableTLS != nil {
		enableTLS = *config.Redis.EnableTLS
	}

	return mongo.NewClient(func(o *mongo.ClientOptions) {
		o.URI = config.URI
		o.Hosts = []string{config.Address}
		o.Username = config.Username
		o.Password = config.Password
		o.ConnectTimeout = config.ConnectTimeout
		o.MaxOpenConnections = config.MaxOpenConnections
		o.EnableTLS = enableTLS
	})
}

// setupSQLClient sets up SQL client (connection).
func setupSQLClient(config *Database) (*sql.DB, error) {
	drv := sql.Driver(config.Driver)
	if drv == sql.DriverMSSQL && config.Database == defaultDatabaseName {
		config.Database = strings.ToUpper(config.Database[:1]) + config.Database[1:]
	}

	var inMemory bool
	if config.SQLite.InMemory != nil {
		inMemory = *config.SQLite.InMemory
	}

	return sql.Open(func(o *sql.Options) {
		o.Driver = drv
		o.DSN = config.URI
		o.Address = config.Address
		o.Database = config.Database
		o.Username = config.Username
		o.Password = config.Password
		o.ConnectTimeout = config.ConnectTimeout
		o.MaxOpenConnections = config.MaxOpenConnections
		o.MaxIdleConnections = config.MaxIdleConnections
		o.MaxConnectionLifetime = config.MaxConnectionLifetime
		o.Postgres.SSLMode = sql.PostgresSSLMode(config.Postgres.SSLMode)
		o.MSSQL.Encrypt = sql.MSSQLEncrypt(config.MSSQL.Encrypt)
		o.SQLite.File = config.SQLite.File
		o.SQLite.InMemory = inMemory
	})
}

// setupSecretStore sets up the secret store.
func setupSecretStore(clients *dbClient, config *Database) (db.SecretStore, error) {
	var store db.SecretStore
	var err error
	switch {
	case clients.mongo != nil:
		store, err = mongo.NewSecretStore(clients.mongo, func(o *mongo.SecretStoreOptions) {
			o.Timeout = config.Timeout
		})
	case clients.sql != nil:
		store, err = sql.NewSecretStore(clients.sql, func(o *sql.SecretStoreOptions) {
			o.Timeout = config.Timeout
		})
	case clients.redis != nil:
		store, err = redis.NewSecretStore(clients.redis)
	default:
		return nil, errors.New("no database clients configured")
	}

	if err != nil {
		return nil, err
	}

	return store, nil
}

// setupRedisClient sets up redis client.
func setupRedisClient(config *Database) (redis.Client, error) {
	var enableTLS bool
	if config.Redis.EnableTLS != nil {
		enableTLS = *config.Redis.EnableTLS
	}

	return redis.NewClient(func(o *redis.ClientOptions) {
		o.URI = config.URI
		o.Address = config.Address
		o.Username = config.Username
		o.Password = config.Password
		o.ConnectTimeout = config.ConnectTimeout
		o.DialTimeout = config.Redis.DialTimeout
		o.MaxRetries = config.Redis.MaxRetries
		o.MinRetryBackoff = config.Redis.MinRetryBackoff
		o.MaxRetryBackoff = config.Redis.MaxRetryBackoff
		o.MaxOpenConnections = config.MaxOpenConnections
		o.MaxIdleConnections = config.MaxIdleConnections
		o.MaxConnectionLifetime = config.MaxConnectionLifetime
		o.EnableTLS = enableTLS
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

	if len(db.SQLite.File) > 0 || db.SQLite.InMemory != nil && *db.SQLite.InMemory {
		return databaseDriverSQLite, nil
	}

	return driver, fmt.Errorf("could not determine database driver")
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
	case databaseDriverMongo, databaseDriverPostgres, databaseDriverMSSQL, databaseDriverSQLite, databaseDriverRedis:
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
	databaseDriverRedis    = "redis"

	databasePorts = map[string]string{
		databaseDriverMongo:    "27017",
		databaseDriverPostgres: "5432",
		databaseDriverMSSQL:    "1433",
		databaseDriverRedis:    "6379",
	}
)
