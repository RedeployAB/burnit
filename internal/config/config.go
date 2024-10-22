package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"time"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

const (
	// defaultListenHost is the default host to listen on.
	defaultListenHost = "0.0.0.0"
	// defaultListenPort is the default port to listen on.
	defaultListenPort = 3000
)

const (
	// defaultConfigPath is the default path to the configuration file.
	defaultConfigPath = "config.yaml"
)

const (
	// defaultSecretServiceTimeout is the default timeout for the secret service.
	defaultSecretServiceTimeout = 10 * time.Second
)

const (
	// defaultDatabaseTimeout is the default timeout for the database.
	defaultDatabaseTimeout = 10 * time.Second
	// defaultDatabaseConnectTimeout is the default connect timeout for the database.
	defaultDatabaseConnectTimeout = 10 * time.Second
	// defaultDatabaseName is the default name for the database.
	defaultDatabaseName = "burnit"
)

const (
	// defaultFrontendRuntimeRenderTemplateDir is the default directory for the runtime render templates.
	defaultFrontendRuntimeRenderTemplateDir = "internal/frontend/templates"
	// defaultFrontendRuntimeRenderStaticDir is the default directory for the runtime render static files.
	defaultFrontendRuntimeRenderStaticDir = "internal/frontend/static"
)

// ConfigOruration contains the configuration for the application.
type Configuration struct {
	Server   Server   `yaml:"server"`
	Services Services `yaml:"services"`
	Frontend Frontend `yaml:"frontend"`
}

// Server contains the configuration for the server.
type Server struct {
	Host        string      `env:"LISTEN_HOST" yaml:"host"`
	Port        int         `env:"LISTEN_PORT" yaml:"port"`
	TLS         TLS         `yaml:"tls"`
	CORS        CORS        `yaml:"cors"`
	RateLimiter RateLimiter `yaml:"rateLimiter"`
	BackendOnly bool        `env:"BACKEND_ONLY" yaml:"backendOnly"`
}

// MarshalJSON returns the JSON encoding of Server. A custom marshalling method
// is defined to hide sensitive values. The reason for not just using the struct tag
// `json:"-"` is that this way we must explicitly set the properties to be marshalled
// and thus output to the logs.
func (s Server) MarshalJSON() ([]byte, error) {
	var rateLimiter *RateLimiter
	if s.RateLimiter.Burst > 0 || s.RateLimiter.Rate > 0 || s.RateLimiter.TTL > 0 || s.RateLimiter.CleanupInterval > 0 {
		rateLimiter = &s.RateLimiter
	}

	return json.Marshal(struct {
		Host        string       `json:",omitempty"`
		Port        int          `json:",omitempty"`
		TLS         *TLS         `json:",omitempty"`
		CORS        *CORS        `json:",omitempty"`
		RateLimiter *RateLimiter `json:",omitempty"`
	}{
		Host:        s.Host,
		Port:        s.Port,
		TLS:         &s.TLS,
		CORS:        &s.CORS,
		RateLimiter: rateLimiter,
	})
}

// TLS contains the configuration server TLS.
type TLS struct {
	CertFile string `env:"TLS_CERT_FILE" yaml:"certFile"`
	KeyFile  string `env:"TLS_KEY_FILE" yaml:"keyFile"`
}

// CORS contains the configuration for CORS.
type CORS struct {
	Origin string `env:"CORS_ORIGIN" yaml:"origin"`
}

// RateLimiter contains the configuration for the rate limiter.
type RateLimiter struct {
	Rate            float64       `env:"RATE_LIMITER_RATE" yaml:"rate"`
	Burst           int           `env:"RATE_LIMITER_BURST" yaml:"burst"`
	TTL             time.Duration `env:"RATE_LIMITER_TTL" yaml:"ttl"`
	CleanupInterval time.Duration `env:"RATE_LIMITER_CLEANUP_INTERVAL" yaml:"cleanupInterval"`
}

// Services contains the configuration for the services.
type Services struct {
	Secret   Secret   `yaml:"secret"`
	Database Database `yaml:"database"`
}

// Secret contains the configuration for the secret service.
type Secret struct {
	Timeout time.Duration `env:"SECRETS_TIMEOUT" yaml:"timeout"`
}

// MarshalJSON returns the JSON encoding of Secret. A custom marshalling method
// is defined to hide sensitive values. The reason for not just using the struct tag
// `json:"-"` is that this way we must explicitly set the properties to be marshalled
// and thus output to the logs.
func (s Secret) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timeout time.Duration `json:",omitempty"`
	}{
		Timeout: s.Timeout,
	})
}

// Database contains the configuration for the database.
type Database struct {
	Driver         string        `env:"DATABASE_DRIVER" yaml:"driver"`
	URI            string        `env:"DATABASE_URI" yaml:"uri"`
	Address        string        `env:"DATABASE_ADDRESS" yaml:"address"`
	Database       string        `env:"DATABASE" yaml:"database"`
	Username       string        `env:"DATABASE_USERNAME" yaml:"username"`
	Password       string        `env:"DATABASE_PASSWORD" yaml:"password"`
	Timeout        time.Duration `env:"DATABASE_TIMEOUT" yaml:"timeout"`
	ConnectTimeout time.Duration `env:"DATABASE_CONNECT_TIMEOUT" yaml:"connectTimeout"`
	Mongo          Mongo         `yaml:"mongo"`
	Postgres       Postgres      `yaml:"postgres"`
	MSSQL          MSSQL         `yaml:"mssql"`
	SQLite         SQLite        `yaml:"sqlite"`
	Redis          Redis         `yaml:"redis"`
}

// MarshalJSON returns the JSON encoding of Database. A custom marshalling method
// is defined to hide sensitive values. The reason for not just using the struct tag
// `json:"-"` is that this way we must explicitly set the properties to be marshalled
// and thus output to the logs.
func (d Database) MarshalJSON() ([]byte, error) {
	var uri string
	reg := regexp.MustCompile(`://.*:.*@`)
	if len(d.URI) > 0 && reg.MatchString(d.URI) {
		uri = reg.ReplaceAllString(d.URI, "://***:***@")
	}

	var mongo *Mongo
	if d.Mongo.EnableTLS != nil {
		mongo = &d.Mongo
	}
	var postgres *Postgres
	if len(d.Postgres.SSLMode) > 0 {
		postgres = &d.Postgres
	}
	var mssql *MSSQL
	if len(d.MSSQL.Encrypt) > 0 {
		mssql = &d.MSSQL
	}
	var sqlite *SQLite
	if len(d.SQLite.File) > 0 || d.SQLite.InMemory != nil {
		sqlite = &d.SQLite
	}
	var redis *Redis
	if d.Redis.DialTimeout > 0 || d.Redis.MaxRetries > 0 || d.Redis.MinRetryBackoff > 0 || d.Redis.MaxRetryBackoff > 0 || d.Redis.EnableTLS != nil {
		redis = &d.Redis
	}

	return json.Marshal(struct {
		Driver         string        `json:",omitempty"`
		URI            string        `json:",omitempty"`
		Address        string        `json:",omitempty"`
		Database       string        `json:",omitempty"`
		Timeout        time.Duration `json:",omitempty"`
		ConnectTimeout time.Duration `json:",omitempty"`
		Mongo          *Mongo        `json:",omitempty"`
		Postgres       *Postgres     `json:",omitempty"`
		MSSQL          *MSSQL        `json:",omitempty"`
		SQLite         *SQLite       `json:",omitempty"`
		Redis          *Redis        `json:",omitempty"`
	}{
		Driver:         d.Driver,
		URI:            uri,
		Address:        d.Address,
		Database:       d.Database,
		Timeout:        d.Timeout,
		ConnectTimeout: d.ConnectTimeout,
		Mongo:          mongo,
		Postgres:       postgres,
		MSSQL:          mssql,
		SQLite:         sqlite,
		Redis:          redis,
	})
}

// Mongo contains the configuration for the Mongo database.
type Mongo struct {
	EnableTLS *bool `env:"DATABASE_MONGO_ENABLE_TLS" yaml:"enableTLS"`
}

// Postgres contains the configuration for the Postgres database.
type Postgres struct {
	SSLMode string `env:"DATABASE_POSTGRES_SSL_MODE" yaml:"sslMode"`
}

// MSSQL contains the configuration for the MSSQL database.
type MSSQL struct {
	Encrypt string `env:"DATABASE_MSSQL_ENCRYPT" yaml:"encrypt"`
}

// SQLite contains the configuration for the SQLite database.
type SQLite struct {
	File     string `env:"DATABASE_SQLITE_FILE" yaml:"file"`
	InMemory *bool  `env:"DATABASE_SQLITE_IN_MEMORY" yaml:"inMemory"`
}

// Redis contains the configuration for the Redis database.
type Redis struct {
	DialTimeout     time.Duration `env:"DATABASE_REDIS_DIAL_TIMEOUT" yaml:"dialTimeout"`
	MaxRetries      int           `env:"DATABASE_REDIS_MAX_RETRIES" yaml:"maxRetries"`
	MinRetryBackoff time.Duration `env:"DATABASE_REDIS_MIN_RETRY_BACKOFF" yaml:"minRetryBackoff"`
	MaxRetryBackoff time.Duration `env:"DATABASE_REDIS_MAX_RETRY_BACKOFF" yaml:"maxRetryBackoff"`
	EnableTLS       *bool         `env:"DATABASE_MONGO_ENABLE_TLS" yaml:"enableTLS"`
}

// Frontend contains the configuration for the frontend.
type Frontend struct {
	RuntimeRender bool `env:"FRONTEND_RUNTIME_RENDER" yaml:"runtimeRender"`
}

// New creates a new Configuration.
func New() (*Configuration, error) {
	flags, _, err := parseFlags(os.Args[1:])
	if err != nil {
		return nil, err
	}

	cfg := &Configuration{
		Server: Server{
			Host: defaultListenHost,
			Port: defaultListenPort,
		},
		Services: Services{
			Secret: Secret{
				Timeout: defaultSecretServiceTimeout,
			},
			Database: Database{
				Database:       defaultDatabaseName,
				Timeout:        defaultDatabaseTimeout,
				ConnectTimeout: defaultDatabaseConnectTimeout,
			},
		},
	}

	yamlCfg, err := configurationFromYAMLFile(flags.configPath)
	if err != nil {
		return nil, err
	}

	envCfg, err := configurationFromEnvironment()
	if err != nil {
		return nil, err
	}

	flagCfg, err := configurationFromFlags(&flags)
	if err != nil {
		return nil, err
	}

	if err := mergeConfigurations(cfg, &yamlCfg, &envCfg, &flagCfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// mergeConfigurations merges the src configuration into the dst configuration.
func mergeConfigurations(dst *Configuration, srcs ...*Configuration) error {
	dstv := reflect.ValueOf(dst)
	for _, src := range srcs {
		srcv := reflect.ValueOf(src)
		if srcv.Kind() != reflect.Ptr || dstv.Kind() != reflect.Ptr {
			return fmt.Errorf("dst and src must be pointers")
		}

		if err := merge(dstv.Elem(), srcv.Elem()); err != nil {
			return err
		}
	}

	return nil
}

// merge the src struct into the dst struct.
func merge(dstv, srcv reflect.Value) error {
	for i := 0; i < srcv.NumField(); i++ {
		srcField := srcv.Field(i)
		dstField := dstv.Field(i)

		if !dstField.CanSet() {
			continue
		}

		if srcField.Kind() == reflect.Ptr && srcField.Elem().Kind() == reflect.Struct {
			if err := merge(dstField.Elem(), srcField.Elem()); err != nil {
				return err
			}
		} else if srcField.Kind() == reflect.Struct {
			if err := merge(dstField, srcField); err != nil {
				return err
			}
		} else {
			if srcField.Kind() == reflect.Ptr && srcField.IsNil() {
				continue
			}
			if srcField.Kind() != reflect.Ptr && srcField.Kind() != reflect.Bool && srcField.IsZero() {
				continue
			}
			dstField.Set(srcField)
		}
	}
	return nil
}

// configurationFromYAMLFile reads the configuration from the given path.
func configurationFromYAMLFile(path string) (Configuration, error) {
	if len(path) == 0 {
		path = defaultConfigPath
	}

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && path == defaultConfigPath {
			return Configuration{}, nil
		}
		return Configuration{}, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return Configuration{}, err
	}

	var cfg Configuration
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Configuration{}, err
	}

	return cfg, nil
}

// configurationFromEnvironment reads the configuration from the environment.
func configurationFromEnvironment() (Configuration, error) {
	var cfg Configuration
	if err := env.ParseWithOptions(&cfg, env.Options{Prefix: "BURNIT_"}); err != nil {
		return Configuration{}, err
	}
	return cfg, nil
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
	}, nil
}
