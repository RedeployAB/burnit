package config

import (
	"net/url"
	"strings"

	"github.com/RedeployAB/burnit/internal/db/mongo"
	"github.com/RedeployAB/burnit/internal/db/redis"
	"github.com/RedeployAB/burnit/internal/db/sql"
)

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

// dbClient contains configured db client.
type dbClient struct {
	mongo mongo.Client
	sql   *sql.DB
	redis redis.Client
}

// databaseDriver returns the database driver. Returns empty string if the driver could not be determined.
func databaseDriver(db *Database) string {
	var driver string
	if len(db.Driver) > 0 && supportedDBDriver(db.Driver) {
		return db.Driver
	}

	if len(db.URI) > 0 {
		driver = dbDriverFromURI(db.URI)
		if len(driver) > 0 && supportedDBDriver(driver) {
			return driver
		}
	}

	if len(db.Address) > 0 {
		driver = dbDriverFromAddress(db.Address)
		if len(driver) > 0 && supportedDBDriver(driver) {
			return driver
		}
	}

	if len(db.SQLite.File) > 0 || db.SQLite.InMemory != nil && *db.SQLite.InMemory {
		return databaseDriverSQLite
	}

	return driver
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

// setupDBClient sets up the db client.
func setupDBClient(config *Database) (*dbClient, error) {
	var client dbClient
	var err error
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
