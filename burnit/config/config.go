package config

import (
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var (
	defaultListenHost = "0.0.0.0"
	defaultListenPort = "3000"
	defaultDBDriver   = "redis"
	defaultDBAddress  = "localhost"
	defaultDB         = "burnit"
	defaultDBSSL      = true
)

const (
	DatabaseDriverRedis = "redis"
	DatabaseDriverMongo = "mongo"
)

// Configuration represents a configuration.
type Configuration struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

func newConfiguration() *Configuration {
	return &Configuration{
		Server: Server{
			Host: defaultListenHost,
			Port: defaultListenPort,
		},
		Database: Database{
			Driver:   defaultDBDriver,
			Address:  defaultDBAddress,
			Database: defaultDB,
			SSL:      defaultDBSSL,
		},
	}
}

// Server defines server part of configuration.
type Server struct {
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	Security Security `yaml:"security"`
}

// Security defines security part of server configuration.
type Security struct {
	Encryption Encryption `yaml:"encryption"`
	TLS        TLS        `yaml:"tls"`
	CORS       CORS       `yaml:"cors"`
}

// Encryption defines encryption pat of security configuration.
type Encryption struct {
	Key string `yaml:"key"`
}

// TLS contains settings for TLS.
type TLS struct {
	Certificate string `yaml:"certificate"`
	Key         string `yaml:"key"`
}

// CORS contains settings for CORS.
type CORS struct {
	Origin string `yaml:"origin,omitempty"`
}

// Database represents database part of configuration.
type Database struct {
	Driver        string `yaml:"driver"`
	Address       string `yaml:"address"`
	URI           string `yaml:"uri"`
	Database      string `yaml:"database"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	SSL           bool   `yaml:"ssl"`
	DirectConnect bool   `yaml:"directConnect"`
}

func New() (*Configuration, error) {
	f := parseFlags()
	return configure(f)
}

func configure(f flags) (*Configuration, error) {
	cfg := newConfiguration()

	if len(f.ConfigPath) != 0 {
		if err := fromFile(cfg, f.ConfigPath); err != nil {
			return nil, err
		}
	}
	fromEnv(cfg)
	fromFlags(cfg, f)

	if len(cfg.Database.URI) > 0 {
		switch cfg.Database.Driver {
		case DatabaseDriverRedis:
			cfg.Database.Address = AddressFromRedisURI(cfg.Database.URI)
		case DatabaseDriverMongo:
			cfg.Database.Address = AddressFromMongoURI(cfg.Database.URI)
		}
	}

	if cfg.Database.Driver == DatabaseDriverRedis {
		re := regexp.MustCompile(`:\d+$`)
		if !re.MatchString(cfg.Database.Address) {
			cfg.Database.Address += ":6379"
		}
	}

	return cfg, nil
}

// fromFile performs the necessary steps
// for server configuration from environment
// variables.
func fromFile(config *Configuration, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var cfg Configuration
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return err
	}

	Merge(config, cfg)
	return nil
}

// fromEnv performs the necessary steps
// for server configuration from environment
// variables.
func fromEnv(config *Configuration) {
	var dbSSL, dbDirectConnect bool
	var err error

	if dbSSLEnv, ok := os.LookupEnv("DB_SSL"); ok {
		dbSSL, err = strconv.ParseBool(dbSSLEnv)
		if err != nil {
			dbSSL = config.Database.SSL
		}
	} else {
		dbSSL = config.Database.SSL
	}

	if dbDirectConnectEnv, ok := os.LookupEnv("DB_DIRECT_CONNECT"); ok {
		dbDirectConnect, err = strconv.ParseBool(dbDirectConnectEnv)
		if err != nil {
			dbDirectConnect = config.Database.DirectConnect
		}
	} else {
		dbDirectConnect = config.Database.DirectConnect
	}

	cfg := Configuration{
		Server{
			Host: os.Getenv("BURNIT_LISTEN_HOST"),
			Port: os.Getenv("BURNIT_LISTEN_PORT"),
			Security: Security{
				Encryption: Encryption{
					Key: os.Getenv("BURNIT_ENCRYPTION_KEY"),
				},
				TLS: TLS{
					Certificate: os.Getenv("BURNIT_TLS_CERTIFICATE"),
					Key:         os.Getenv("BURNIT_TLS_KEY"),
				},
				CORS: CORS{
					Origin: os.Getenv("BURNIT_CORS_ORIGIN"),
				},
			},
		},
		Database{
			Driver:        strings.ToLower(os.Getenv("DB_DRIVER")),
			Address:       os.Getenv("DB_HOST"),
			URI:           os.Getenv("DB_CONNECTION_URI"),
			Database:      os.Getenv("DB"),
			Username:      os.Getenv("DB_USER"),
			Password:      os.Getenv("DB_PASSWORD"),
			SSL:           dbSSL,
			DirectConnect: dbDirectConnect,
		},
	}
	Merge(config, cfg)
}

// fromFlags takes incoming flags and creates
// a configuration object from it.
func fromFlags(config *Configuration, f flags) {
	dbSSL := config.Database.SSL
	if f.DisableDBSSL {
		dbSSL = false
	}

	dbDirectConnect := config.Database.DirectConnect
	if f.DBDirectConnect {
		dbDirectConnect = f.DBDirectConnect
	}

	cfg := Configuration{
		Server{
			Host: f.Host,
			Port: f.Port,
			Security: Security{
				Encryption: Encryption{
					Key: f.EncryptionKey,
				},
				TLS: TLS{
					Certificate: f.TLSCertificate,
					Key:         f.TLSKey,
				},
				CORS: CORS{
					Origin: f.CORSOrigin,
				},
			},
		},
		Database{
			Driver:        f.Driver,
			Address:       f.DBAddress,
			URI:           f.DBURI,
			Database:      f.DB,
			Username:      f.DBUser,
			Password:      f.DBPassword,
			SSL:           dbSSL,
			DirectConnect: dbDirectConnect,
		},
	}
	Merge(config, cfg)
}

// AddressFromMongoURI returns the address (<host>:<port>) from
// a mongodb connection string.
func AddressFromMongoURI(uri string) string {
	if !strings.HasSuffix(uri, "/") {
		uri += "/"
	}
	var address string
	if strings.Contains(uri, "@") {
		address = strings.Split(strings.Split(uri, "@")[1], "/")[0]
	} else {
		address = strings.Split(uri, "/")[2]
	}

	return address
}

// AddressFromRedisURI returns the address (<host>:<port>) from
// a redis connection string.
func AddressFromRedisURI(uri string) string {
	reg := regexp.MustCompile("^redis://|^rediss://")
	res := reg.ReplaceAllString(uri, "${1}")
	return strings.Split(res, ",")[0]
}

// NewTLSConfig creates and returns a new config for TLS.
func NewTLSConfig() *tls.Config {
	return &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

// CORSHeaders creates and returns headers for CORS.
func CORSHeaders() http.Header {
	return http.Header{
		"Access-Control-Allow-Headers": {"Content-Type", "Passphrase"},
		"Access-Control-Allow-Methods": {"GET", "POST", "OPTIONS"},
	}
}
