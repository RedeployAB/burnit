package config

import (
	"bytes"
	"flag"
	"strconv"
	"time"
)

// flags contains the flags.
type flags struct {
	configPath                   string
	host                         string
	port                         int
	tlsCertFile                  string
	tlsKeyFile                   string
	encryptionKey                string
	corsOrigin                   string
	timeout                      time.Duration
	databaseDriver               string
	databaseURI                  string
	databaseAddr                 string
	database                     string
	databaseUser                 string
	databasePass                 string
	databaseTimeout              time.Duration
	databaseConnectTimeout       time.Duration
	databaseTLS                  string
	databaseSQLiteFile           string
	databaseSQLiteInMemory       *bool
	databaseRedisDialTimeout     time.Duration
	databaseRedisMaxRetries      int
	databaseRedisMinRetryBackoff time.Duration
	databaseRedisMaxRetryBackoff time.Duration
}

// parseFlags parses the flags.
func parseFlags(args []string) (flags, string, error) {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	var buf bytes.Buffer
	fs.SetOutput(&buf)

	var f flags
	var databaseSQLiteInMemory boolFlag

	fs.StringVar(&f.configPath, "config-path", "", "Optional. Path to a configuration file. Defaults to: "+defaultConfigPath+".")
	fs.StringVar(&f.host, "host", "", "Optional. Host to listen on. Defaults to: "+defaultListenHost+".")
	fs.IntVar(&f.port, "port", 0, "Optional. Port to listen on. Defaults to: "+strconv.Itoa(defaultListenPort)+".")
	fs.StringVar(&f.tlsCertFile, "tls-cert-file", "", "Optional. TLS certificate file.")
	fs.StringVar(&f.tlsKeyFile, "tls-key-file", "", "Optional. TLS key file.")
	fs.StringVar(&f.encryptionKey, "encryption-key", "", "Optional. Default encryption key for the secrets service.")
	fs.StringVar(&f.corsOrigin, "cors-origin", "", "Optional. CORS origin.")
	fs.DurationVar(&f.timeout, "timeout", 0, "Optional. Default timeout for the service. Defaults to: "+defaultSecretServiceTimeout.String()+".")
	fs.StringVar(&f.databaseDriver, "database-driver", "", "Optional. Database driver.")
	fs.StringVar(&f.databaseURI, "database-uri", "", "Optional. URI for the database.")
	fs.StringVar(&f.databaseAddr, "database-address", "", "Optional. Address for the database.")
	fs.StringVar(&f.database, "database", "", "Optional. Database name.")
	fs.StringVar(&f.databaseUser, "database-user", "", "Optional. Database username.")
	fs.StringVar(&f.databasePass, "database-password", "", "Optional. Database password.")
	fs.DurationVar(&f.databaseTimeout, "database-timeout", 0, "Optional. Timeout for the database. Defaults to: "+defaultDatabaseTimeout.String()+".")
	fs.DurationVar(&f.databaseConnectTimeout, "database-connect-timeout", 0, "Optional. Connect timeout for the database. Defaults to: "+defaultDatabaseConnectTimeout.String()+".")
	fs.StringVar(&f.databaseTLS, "database-tls", "", "Optional. Enable and set TLS mode for the database.")
	fs.StringVar(&f.databaseSQLiteFile, "database-sqlite-file", "", "Optional. Path to the database file for SQLite.")
	fs.Var(&databaseSQLiteInMemory, "database-sqlite-in-memory", "Optional. Use an in-memory database for SQLite. Defaults to: false.")
	fs.DurationVar(&f.databaseRedisDialTimeout, "database-redis-dial-timeout", 0, "Optional. Dial timeout for the Redis client.")
	fs.IntVar(&f.databaseRedisMaxRetries, "database-redis-max-retries", 0, "Optional. Maximum number of retries for the Redis client.")
	fs.DurationVar(&f.databaseRedisMinRetryBackoff, "database-redis-min-retry-backoff", 0, "Optional. Minimum retry backoff for the Redis client.")
	fs.DurationVar(&f.databaseRedisMaxRetryBackoff, "database-redis-max-retry-backoff", 0, "Optional. Maximum retry backoff for the Redis client.")

	if err := fs.Parse(args); err != nil {
		return f, buf.String(), err
	}

	if databaseSQLiteInMemory.isSet {
		f.databaseSQLiteInMemory = &databaseSQLiteInMemory.value
	}

	return f, buf.String(), nil
}

// boolFlag is a flag for bool values that keeps track if it was set.
type boolFlag struct {
	value bool
	isSet bool
}

// Set sets the value of the boolFlag.
func (f *boolFlag) Set(value string) error {
	var v bool
	if value == "true" {
		v = true
	}
	f.value = v
	f.isSet = true
	return nil
}

// String returns the string representation of the boolFlag.
func (f *boolFlag) String() string {
	if f.value {
		return "true"
	}
	return "false"
}
