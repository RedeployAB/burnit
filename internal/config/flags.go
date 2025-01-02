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
	corsOrigin                   string
	rateLimiterRate              float64
	rateLimiterBurst             int
	rateLimiterCleanupInterval   time.Duration
	rateLimiterTTL               time.Duration
	timeout                      time.Duration
	databaseDriver               string
	databaseURI                  string
	databaseAddr                 string
	database                     string
	databaseUser                 string
	databasePass                 string
	databaseTimeout              time.Duration
	databaseConnectTimeout       time.Duration
	databaseMongoEnableTLS       *bool
	databasePostgresSSLMode      string
	databaseMSSQLEncrypt         string
	databaseSQLiteFile           string
	databaseSQLiteInMemory       *bool
	databaseRedisDialTimeout     time.Duration
	databaseRedisMaxRetries      int
	databaseRedisMinRetryBackoff time.Duration
	databaseRedisMaxRetryBackoff time.Duration
	databaseRedisEnableTLS       *bool
	// UI flags.
	runtimeRender *bool
	// Session database flags.
	sessionDatabaseDriver               string
	sessionDatabaseURI                  string
	sessionDatabaseAddr                 string
	sessionDatabase                     string
	sessionDatabaseUser                 string
	sessionDatabasePass                 string
	sessionDatabaseTimeout              time.Duration
	sessionDatabaseConnectTimeout       time.Duration
	sessionDatabaseMongoEnableTLS       *bool
	sessionDatabasePostgresSSLMode      string
	sessionDatabaseMSSQLEncrypt         string
	sessionDatabaseSQLiteFile           string
	sessionDatabaseSQLiteInMemory       *bool
	sessionDatabaseRedisDialTimeout     time.Duration
	sessionDatabaseRedisMaxRetries      int
	sessionDatabaseRedisMinRetryBackoff time.Duration
	sessionDatabaseRedisMaxRetryBackoff time.Duration
	sessionDatabaseRedisEnableTLS       *bool
	// Local development flag.
	localDevelopment *bool
}

// parseFlags parses the flags.
func parseFlags(args []string) (flags, string, error) {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	var buf bytes.Buffer
	fs.SetOutput(&buf)

	var (
		f                             flags
		runtimeRender                 boolFlag
		databaseMongoEnableTLS        boolFlag
		databaseSQLiteInMemory        boolFlag
		databaseRedisEnableTLS        boolFlag
		sessionDatabaseMongoEnableTLS boolFlag
		sessionDatabaseSQLiteInMemory boolFlag
		sessionDatabaseRedisEnableTLS boolFlag
		localDevelopment              boolFlag
	)

	fs.StringVar(&f.configPath, "config-path", "", "Optional. Path to a configuration file. Defaults to: "+defaultConfigPath+".")
	fs.StringVar(&f.host, "host", "", "Optional. Host to listen on. Defaults to: "+defaultListenHost+".")
	fs.IntVar(&f.port, "port", 0, "Optional. Port to listen on. Defaults to: "+strconv.Itoa(defaultListenPort)+".")
	fs.StringVar(&f.tlsCertFile, "tls-cert-file", "", "Optional. TLS certificate file.")
	fs.StringVar(&f.tlsKeyFile, "tls-key-file", "", "Optional. TLS key file.")
	fs.StringVar(&f.corsOrigin, "cors-origin", "", "Optional. CORS origin.")
	fs.Float64Var(&f.rateLimiterRate, "rate-limiter-rate", 0, "Optional. The average number of requests per second.")
	fs.IntVar(&f.rateLimiterBurst, "rate-limiter-burst", 0, "Optional. The maximum burst of requests.")
	fs.DurationVar(&f.rateLimiterCleanupInterval, "rate-limiter-cleanup-interval", 0, "Optional. The interval at which to clean up stale rate limiters.")
	fs.DurationVar(&f.rateLimiterTTL, "rate-limiter-ttl", 0, "Optional. The time-to-live for rate limiters.")
	fs.DurationVar(&f.timeout, "timeout", 0, "Optional. Default timeout for the service. Defaults to: "+defaultSecretServiceTimeout.String()+".")
	// Database flags.
	fs.StringVar(&f.databaseDriver, "database-driver", "", "Optional. Database driver.")
	fs.StringVar(&f.databaseURI, "database-uri", "", "Optional. URI for the database.")
	fs.StringVar(&f.databaseAddr, "database-address", "", "Optional. Address for the database.")
	fs.StringVar(&f.database, "database", "", "Optional. Database name.")
	fs.StringVar(&f.databaseUser, "database-user", "", "Optional. Database username.")
	fs.StringVar(&f.databasePass, "database-password", "", "Optional. Database password.")
	fs.DurationVar(&f.databaseTimeout, "database-timeout", 0, "Optional. Timeout for the database. Defaults to: "+defaultDatabaseTimeout.String()+".")
	fs.DurationVar(&f.databaseConnectTimeout, "database-connect-timeout", 0, "Optional. Connect timeout for the database. Defaults to: "+defaultDatabaseConnectTimeout.String()+".")
	fs.Var(&databaseMongoEnableTLS, "database-mongo-enable-tls", "Optional. Enable TLS for MongoDB. Defaults to: true.")
	fs.StringVar(&f.databasePostgresSSLMode, "database-postgres-ssl-mode", "", "Optional. SSL mode for PostgreSQL. Defaults to: require.")
	fs.StringVar(&f.databaseMSSQLEncrypt, "database-mssql-encrypt", "", "Optional. Encrypt for MSSQL. Defaults to: true.")
	fs.StringVar(&f.databaseSQLiteFile, "database-sqlite-file", "", "Optional. Path to the database file for SQLite.")
	fs.Var(&databaseSQLiteInMemory, "database-sqlite-in-memory", "Optional. Use an in-memory database for SQLite. Defaults to: false.")
	fs.DurationVar(&f.databaseRedisDialTimeout, "database-redis-dial-timeout", 0, "Optional. Dial timeout for the Redis client.")
	fs.IntVar(&f.databaseRedisMaxRetries, "database-redis-max-retries", 0, "Optional. Maximum number of retries for the Redis client.")
	fs.DurationVar(&f.databaseRedisMinRetryBackoff, "database-redis-min-retry-backoff", 0, "Optional. Minimum retry backoff for the Redis client.")
	fs.DurationVar(&f.databaseRedisMaxRetryBackoff, "database-redis-max-retry-backoff", 0, "Optional. Maximum retry backoff for the Redis client.")
	fs.Var(&databaseRedisEnableTLS, "database-redis-enable-tls", "Optional. Enable TLS for the Redis client. Defaults to: true.")
	// UI flags.
	fs.Var(&runtimeRender, "runtime-render", "Optional. Enable runtime rendering of the UI.")
	// Session database flags.
	fs.StringVar(&f.sessionDatabaseDriver, "session-database-driver", "", "Optional. Database driver.")
	fs.StringVar(&f.sessionDatabaseURI, "session-database-uri", "", "Optional. URI for the database.")
	fs.StringVar(&f.sessionDatabaseAddr, "session-database-address", "", "Optional. Address for the database.")
	fs.StringVar(&f.sessionDatabase, "session-database", "", "Optional. Database name.")
	fs.StringVar(&f.sessionDatabaseUser, "session-database-user", "", "Optional. Database username.")
	fs.StringVar(&f.sessionDatabasePass, "session-database-password", "", "Optional. Database password.")
	fs.DurationVar(&f.sessionDatabaseTimeout, "session-database-timeout", 0, "Optional. Timeout for the database. Defaults to: "+defaultDatabaseTimeout.String()+".")
	fs.DurationVar(&f.sessionDatabaseConnectTimeout, "session-database-connect-timeout", 0, "Optional. Connect timeout for the database. Defaults to: "+defaultDatabaseConnectTimeout.String()+".")
	fs.Var(&sessionDatabaseMongoEnableTLS, "session-database-mongo-enable-tls", "Optional. Enable TLS for MongoDB. Defaults to: true.")
	fs.StringVar(&f.sessionDatabasePostgresSSLMode, "session-database-postgres-ssl-mode", "", "Optional. SSL mode for PostgreSQL. Defaults to: require.")
	fs.StringVar(&f.sessionDatabaseMSSQLEncrypt, "session-database-mssql-encrypt", "", "Optional. Encrypt for MSSQL. Defaults to: true.")
	fs.StringVar(&f.sessionDatabaseSQLiteFile, "session-database-sqlite-file", "", "Optional. Path to the database file for SQLite.")
	fs.Var(&sessionDatabaseSQLiteInMemory, "session-database-sqlite-in-memory", "Optional. Use an in-memory database for SQLite. Defaults to: false.")
	fs.DurationVar(&f.sessionDatabaseRedisDialTimeout, "session-database-redis-dial-timeout", 0, "Optional. Dial timeout for the Redis client.")
	fs.IntVar(&f.sessionDatabaseRedisMaxRetries, "session-database-redis-max-retries", 0, "Optional. Maximum number of retries for the Redis client.")
	fs.DurationVar(&f.sessionDatabaseRedisMinRetryBackoff, "session-database-redis-min-retry-backoff", 0, "Optional. Minimum retry backoff for the Redis client.")
	fs.DurationVar(&f.sessionDatabaseRedisMaxRetryBackoff, "session-database-redis-max-retry-backoff", 0, "Optional. Maximum retry backoff for the Redis client.")
	fs.Var(&sessionDatabaseRedisEnableTLS, "session-database-redis-enable-tls", "Optional. Enable TLS for the Redis client. Defaults to: true.")

	// Local development flag.
	fs.Var(&localDevelopment, "local-development", "Optional. Enable local development mode.")

	if err := fs.Parse(args); err != nil {
		return f, buf.String(), err
	}

	if databaseMongoEnableTLS.isSet {
		f.databaseMongoEnableTLS = &databaseMongoEnableTLS.value
	}
	if databaseSQLiteInMemory.isSet {
		f.databaseSQLiteInMemory = &databaseSQLiteInMemory.value
	}
	if databaseRedisEnableTLS.isSet {
		f.databaseRedisEnableTLS = &databaseRedisEnableTLS.value
	}
	if runtimeRender.isSet {
		f.runtimeRender = &runtimeRender.value
	}
	if sessionDatabaseMongoEnableTLS.isSet {
		f.sessionDatabaseMongoEnableTLS = &sessionDatabaseMongoEnableTLS.value
	}
	if sessionDatabaseSQLiteInMemory.isSet {
		f.sessionDatabaseSQLiteInMemory = &sessionDatabaseSQLiteInMemory.value
	}
	if sessionDatabaseRedisEnableTLS.isSet {
		f.sessionDatabaseRedisEnableTLS = &sessionDatabaseRedisEnableTLS.value
	}
	if localDevelopment.isSet {
		f.localDevelopment = &localDevelopment.value
	}

	return f, buf.String(), nil
}

// configurationFromFlags reads the configuration from the flags.
func configurationFromFlags(flags *flags) (Configuration, error) {
	return Configuration{
		Server: Server{
			Host: flags.host,
			Port: flags.port,
			TLS: TLS{
				CertFile: flags.tlsCertFile,
				KeyFile:  flags.tlsKeyFile,
			},
			CORS: CORS{
				Origin: flags.corsOrigin,
			},
			RateLimiter: RateLimiter{
				Rate:            flags.rateLimiterRate,
				Burst:           flags.rateLimiterBurst,
				CleanupInterval: flags.rateLimiterCleanupInterval,
				TTL:             flags.rateLimiterTTL,
			},
		},
		Services: Services{
			Secret: Secret{
				Timeout: flags.timeout,
			},
			Database: Database{
				Driver:         flags.databaseDriver,
				URI:            flags.databaseURI,
				Address:        flags.databaseAddr,
				Database:       flags.database,
				Username:       flags.databaseUser,
				Password:       flags.databasePass,
				Timeout:        flags.databaseTimeout,
				ConnectTimeout: flags.databaseConnectTimeout,
				Mongo: Mongo{
					EnableTLS: flags.databaseMongoEnableTLS,
				},
				Postgres: Postgres{
					SSLMode: flags.databasePostgresSSLMode,
				},
				MSSQL: MSSQL{
					Encrypt: flags.databaseMSSQLEncrypt,
				},
				SQLite: SQLite{
					File:     flags.databaseSQLiteFile,
					InMemory: flags.databaseSQLiteInMemory,
				},
				Redis: Redis{
					DialTimeout:     flags.databaseRedisDialTimeout,
					MaxRetries:      flags.databaseRedisMaxRetries,
					MinRetryBackoff: flags.databaseRedisMinRetryBackoff,
					MaxRetryBackoff: flags.databaseRedisMaxRetryBackoff,
					EnableTLS:       flags.databaseRedisEnableTLS,
				},
			},
		},
		UI: UI{
			RuntimeRender: flags.runtimeRender,
			Services: UIServices{
				Session: Session{
					Database: SessionDatabase{
						Driver:         flags.sessionDatabaseDriver,
						URI:            flags.sessionDatabaseURI,
						Address:        flags.sessionDatabaseAddr,
						Database:       flags.sessionDatabase,
						Username:       flags.sessionDatabaseUser,
						Password:       flags.sessionDatabasePass,
						Timeout:        flags.sessionDatabaseTimeout,
						ConnectTimeout: flags.sessionDatabaseConnectTimeout,
						Mongo: SessionMongo{
							EnableTLS: flags.sessionDatabaseMongoEnableTLS,
						},
						Postgres: SessionPostgres{
							SSLMode: flags.sessionDatabasePostgresSSLMode,
						},
						MSSQL: SessionMSSQL{
							Encrypt: flags.sessionDatabaseMSSQLEncrypt,
						},
						SQLite: SessionSQLite{
							File:     flags.sessionDatabaseSQLiteFile,
							InMemory: flags.sessionDatabaseSQLiteInMemory,
						},
						Redis: SessionRedis{
							DialTimeout:     flags.sessionDatabaseRedisDialTimeout,
							MaxRetries:      flags.sessionDatabaseRedisMaxRetries,
							MinRetryBackoff: flags.sessionDatabaseRedisMinRetryBackoff,
							MaxRetryBackoff: flags.sessionDatabaseRedisMaxRetryBackoff,
							EnableTLS:       flags.sessionDatabaseRedisEnableTLS,
						},
					},
				},
			},
		},
	}, nil
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
