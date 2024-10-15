package config

import (
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

// SetupServices configures and sets up the services.
func SetupServices(config Services) (*services, error) {
	repo, err := setupSecretRepository(&config.Database)
	if err != nil {
		return nil, fmt.Errorf("failed to setup secret repository: %w", err)
	}

	secrets, err := secret.NewService(
		repo,
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
	var err error
	config.Driver, err = databaseDriver(config)
	if err != nil {
		return nil, err
	}

	switch config.Driver {
	case databaseDriverMongo:
		repo, err = setupMongoSecretRepository(config)
	case databaseDriverPostgres, databaseDriverMSSQL, databaseDriverSQLite:
		repo, err = setupSQLSecretRepository(config, config.Driver)
	case databaseDriverRedis:
		repo, err = setupRedisSecretRepository(config)
	default:
		return nil, fmt.Errorf("unsupported database driver")
	}

	if err != nil {
		return nil, err
	}

	return repo, nil
}

// setupMongoSecretRepository sets up the MongoDB secret repository.
func setupMongoSecretRepository(config *Database) (*mongo.SecretRepository, error) {
	var enableTLS bool
	if config.Redis.EnableTLS != nil {
		enableTLS = *config.Redis.EnableTLS
	}

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

// setupSQLSecretRepository sets up the SQL secret repository.
func setupSQLSecretRepository(config *Database, driver string) (*sql.SecretRepository, error) {
	drv := sql.Driver(driver)
	if drv == sql.DriverMSSQL && config.Database == defaultDatabaseName {
		config.Database = strings.ToUpper(config.Database[:1]) + config.Database[1:]
	}

	var inMemory bool
	if config.SQLite.InMemory != nil {
		inMemory = *config.SQLite.InMemory
	}

	db, err := sql.Open(func(o *sql.Options) {
		o.Driver = drv
		o.DSN = config.URI
		o.Address = config.Address
		o.Database = config.Database
		o.Username = config.Username
		o.Password = config.Password
		o.Postgres.SSLMode = sql.PostgresSSLMode(config.Postgres.SSLMode)
		o.MSSQL.Encrypt = sql.MSSQLEncrypt(config.MSSQL.Encrypt)
		o.SQLite.File = config.SQLite.File
		o.SQLite.InMemory = inMemory
		o.ConnectTimeout = config.ConnectTimeout
	})
	if err != nil {
		return nil, err
	}

	return sql.NewSecretRepository(db, func(o *sql.SecretRepositoryOptions) {
		o.Timeout = config.Timeout
	})
}

// setupRedisSecretRepository sets up the Redis secret repository.
func setupRedisSecretRepository(config *Database) (*redis.SecretRepository, error) {
	var enableTLS bool
	if config.Redis.EnableTLS != nil {
		enableTLS = *config.Redis.EnableTLS
	}

	client, err := redis.NewClient(func(o *redis.ClientOptions) {
		o.URI = config.URI
		o.Address = config.Address
		o.Username = config.Username
		o.Password = config.Password
		o.ConnectTimeout = config.ConnectTimeout
		o.EnableTLS = enableTLS
		o.DialTimeout = config.Redis.DialTimeout
		o.MaxRetries = config.Redis.MaxRetries
		o.MinRetryBackoff = config.Redis.MinRetryBackoff
		o.MaxRetryBackoff = config.Redis.MaxRetryBackoff
	})
	if err != nil {
		return nil, err
	}

	return redis.NewSecretRepository(client, func(o *redis.SecretRepositoryOptions) {})
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
