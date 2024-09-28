package config

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
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
	// defaultDatabaseTimeout is the default timeout for the database.
	defaultDatabaseTimeout = 10 * time.Second
	// defaultDatabaseConnectTimeout is the default connect timeout for the database.
	defaultDatabaseConnectTimeout = 10 * time.Second
	// defaultDatabaseName is the default name for the database.
	defaultDatabaseName = "burnit"
)

const (
	// defaultSecretServiceTimeout is the default timeout for the secret service.
	defaultSecretServiceTimeout = 10 * time.Second
)

// ConfigOruration contains the configuration for the application.
type Configuration struct {
	Server   Server   `yaml:"server"`
	Services Services `yaml:"services"`
}

// Server contains the configuration for the server.
type Server struct {
	Host string `env:"LISTEN_HOST" yaml:"host"`
	Port int    `env:"LISTEN_PORT" yaml:"port"`
	TLS  TLS    `yaml:"tls"`
	CORS CORS   `yaml:"cors"`
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

// Services contains the configuration for the services.
type Services struct {
	Secret   Secret   `yaml:"secret"`
	Database Database `yaml:"database"`
}

// Secret contains the configuration for the secret service.
type Secret struct {
	EncryptionKey string        `env:"SECRETS_ENCRYPTION_KEY" yaml:"encryptionKey"`
	Timeout       time.Duration `env:"SECRETS_TIMEOUT" yaml:"timeout"`
}

// Database contains the configuration for the database.
type Database struct {
	ConnectionString string        `env:"DATABASE_CONNECTION_STRING" yaml:"connectionString"`
	Address          string        `env:"DATABASE_ADDRESS" yaml:"address"`
	Database         string        `env:"DATABASE" yaml:"database"`
	Username         string        `env:"DATABASE_USERNAME" yaml:"username"`
	Password         string        `env:"DATABASE_PASSWORD" yaml:"password"`
	Timeout          time.Duration `env:"DATABASE_TIMEOUT" yaml:"timeout"`
	ConnectTimeout   time.Duration `env:"DATABASE_CONNECT_TIMEOUT" yaml:"connectTimeout"`
	EnableTLS        *bool         `env:"DATABASE_ENABLE_TLS" yaml:"enableTLS,omitempty"`
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
				EnableTLS:      toPtr(true),
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
		},
		Services: Services{
			Secret: Secret{
				EncryptionKey: flags.encryptionKey,
				Timeout:       flags.timeout,
			},
			Database: Database{
				ConnectionString: flags.databaseConnStr,
				Address:          flags.databaseAddr,
				Database:         flags.database,
				Username:         flags.databaseUser,
				Password:         flags.databasePass,
				Timeout:          flags.databaseTimeout,
				ConnectTimeout:   flags.databaseConnectTimeout,
				EnableTLS:        flags.databaseEnableTLS,
			},
		},
	}, nil
}

// toPtr returns a pointer to the value v.
func toPtr[T any](v T) *T {
	return &v
}
