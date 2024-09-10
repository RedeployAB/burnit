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
	defaultListenPort = "3000"
)

const (
	defaultConfigPath = "config.yaml"
)

const (
	defaultDatabaseTimeout        = 10 * time.Second
	defaultDatabaseConnectTimeout = 10 * time.Second
)

const (
	// defaultSecretServiceTimeout is the default timeout for the secret service.
	defaultSecretServiceTimeout = 10 * time.Second
)

// Configuration contains the configuration for the application.
type Configuration struct {
	Server   Server
	Services Services
}

// Server contains the configuration for the server.
type Server struct {
	Host string
	Port string
}

// Services contains the configuration for the services.
type Services struct {
	Secrets  Secrets  `yaml:"secrets"`
	Database Database `yaml:"database"`
}

// Secrets contains the configuration for the secret service.
type Secrets struct {
	EncryptionKey string        `env:"SECRETS_ENCRYPTION_KEY" yaml:"encryptionKey"`
	Timeout       time.Duration `env:"SECRETS_TIMEOUT" yaml:"timeout"`
}

// Database contains the configuration for the database.
type Database struct {
	URI            string        `env:"DATABASE_URI" yaml:"uri"`
	Address        string        `env:"DATABASE_ADDRESS" yaml:"address"`
	Database       string        `env:"DATABASE" yaml:"database"`
	Username       string        `env:"DATABASE_USER" yaml:"username"`
	Password       string        `env:"DATABASE_PASSWORD" yaml:"password"`
	Collection     string        `env:"DATABASE_COLLECTION" yaml:"collection"`
	Timeout        time.Duration `env:"DATABASE_TIMEOUT" yaml:"timeout"`
	ConnectTimeout time.Duration `env:"DATABASE_CONNECT_TIMEOUT" yaml:"connectTimeout"`
	EnableTLS      *bool         `env:"DATABASE_ENABLE_TLS" yaml:"enableTLS,omitempty"`
}

// New creates a new Configuration.
func New() (*Configuration, error) {
	cfg := &Configuration{
		Server: Server{
			Host: defaultListenHost,
			Port: defaultListenPort,
		},
		Services: Services{
			Secrets: Secrets{
				Timeout: defaultSecretServiceTimeout,
			},
			Database: Database{
				ConnectTimeout: defaultDatabaseConnectTimeout,
			},
		},
	}

	// Load YAML configuration.

	if err := env.ParseWithOptions(cfg, env.Options{Prefix: "BURNIT_"}); err != nil {
		return nil, err
	}

	// Parse flags.

	return cfg, nil
}

// mergeConfigurations merges the src configuration into the dst configuration.
func mergeConfigurations(dst, src *Configuration) error {
	dstv := reflect.ValueOf(dst)
	srcv := reflect.ValueOf(src)

	if srcv.Kind() != reflect.Ptr || dstv.Kind() != reflect.Ptr {
		return fmt.Errorf("src and dst must be pointers")
	}

	if err := merge(dstv.Elem(), srcv.Elem()); err != nil {
		return err
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

// ReadYAMLConfiguration reads the configuration from the given path.
func ReadYAMLConfiguration(path string) (*Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) && path == defaultConfigPath {
			return &Configuration{}, nil
		}
		return nil, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var cfg Configuration
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
