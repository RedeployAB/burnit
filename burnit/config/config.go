package config

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	defaultHost            = "0.0.0.0"
	defaultListenPort      = "3000"
	defaultDriver          = "redis"
	defaultAddress         = "localhost"
	defaultDB              = "burnit"
	defaultDBSSL           = true
	defaultDBDirectConnect = false
)

// Server defines server part of configuration.
type Server struct {
	Host     string   `yaml:"host"`
	Port     string   `yaml:"port"`
	Security Security `yaml:"security"`
}

// Security defines security part of server configuration.
type Security struct {
	APIKey     string     `yaml:"apiKey"`
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

// Configuration represents a configuration.
type Configuration struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

// Configure performs the necessary steps
// for server configuration.
func Configure(f Flags) (*Configuration, error) {
	config := &Configuration{
		Server{
			Host: defaultHost,
			Port: defaultListenPort,
		},
		Database{
			Driver:        defaultDriver,
			Address:       defaultAddress,
			Database:      defaultDB,
			SSL:           defaultDBSSL,
			DirectConnect: defaultDBDirectConnect,
		},
	}

	if len(f.ConfigPath) > 0 {
		if err := configureFromFile(config, f.ConfigPath); err != nil {
			return nil, err
		}
	}

	configureFromEnv(config)
	configureFromFlags(config, f)

	if len(config.Server.Security.Encryption.Key) == 0 {
		return config, errors.New("encryption key must be set")
	}

	if len(config.Database.URI) > 0 {
		switch config.Database.Driver {
		case "redis":
			config.Database.Address = AddressFromRedisURI(config.Database.URI)
		case "mongo":
			config.Database.Address = AddressFromMongoURI(config.Database.URI)
		}
	}

	if config.Database.Driver == "redis" {
		re := regexp.MustCompile(`:\d+$`)
		if !re.MatchString(config.Database.Address) {
			config.Database.Address += ":6379"
		}
	}
	return config, nil
}

func mergeConfig(config *Configuration, srcCfg Configuration) {
	if len(srcCfg.Server.Host) > 0 {
		config.Server.Host = srcCfg.Server.Host
	}
	if len(srcCfg.Server.Port) > 0 {
		config.Server.Port = srcCfg.Server.Port
	}
	if len(srcCfg.Server.Security.TLS.Certificate) > 0 {
		config.Server.Security.TLS.Certificate = srcCfg.Server.Security.TLS.Certificate
	}
	if len(srcCfg.Server.Security.TLS.Key) > 0 {
		config.Server.Security.TLS.Key = srcCfg.Server.Security.TLS.Key
	}
	if len(srcCfg.Server.Security.CORS.Origin) > 0 {
		config.Server.Security.CORS.Origin = srcCfg.Server.Security.CORS.Origin
	}
	if len(srcCfg.Server.Security.APIKey) > 0 {
		config.Server.Security.APIKey = srcCfg.Server.Security.APIKey
	}
	if len(srcCfg.Server.Security.Encryption.Key) > 0 {
		config.Server.Security.Encryption.Key = srcCfg.Server.Security.Encryption.Key
	}

	if len(srcCfg.Database.Driver) > 0 {
		config.Database.Driver = strings.ToLower(srcCfg.Database.Driver)
	}
	if len(srcCfg.Database.Address) > 0 {
		config.Database.Address = srcCfg.Database.Address
	}
	if len(srcCfg.Database.URI) > 0 {
		config.Database.URI = srcCfg.Database.URI
	}
	if len(srcCfg.Database.Database) > 0 {
		config.Database.Database = srcCfg.Database.Database
	}
	if len(srcCfg.Database.Username) > 0 {
		config.Database.Username = srcCfg.Database.Username
	}
	if len(srcCfg.Database.Password) > 0 {
		config.Database.Password = srcCfg.Database.Password
	}
	config.Database.SSL = srcCfg.Database.SSL
	config.Database.DirectConnect = srcCfg.Database.DirectConnect
}

// configureFromFile performs the necessary steps
// for server configuration from environment
// variables.
func configureFromFile(config *Configuration, path string) error {
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

	mergeConfig(config, cfg)
	return nil
}

// configureFromEnv performs the necessary steps
// for server configuration from environment
// variables.
func configureFromEnv(config *Configuration) {
	var dbSSL bool
	var dbDirect bool

	dbSSLStr := os.Getenv("DB_SSL")
	if len(dbSSLStr) == 0 {
		dbSSL = config.Database.SSL
	} else {
		dbSSL, _ = strconv.ParseBool(dbSSLStr)
	}

	dbDirectStr := os.Getenv("DB_DIRECT_CONNECT")
	if len(dbDirectStr) == 0 {
		dbDirect = config.Database.DirectConnect
	} else {
		dbDirect, _ = strconv.ParseBool(dbDirectStr)
	}

	cfg := Configuration{
		Server{
			Host: os.Getenv("BURNIT_LISTEN_HOST"),
			Port: os.Getenv("BURNIT_LISTEN_PORT"),
			Security: Security{
				APIKey: os.Getenv("BURNIT_API_KEY"),
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
			DirectConnect: dbDirect,
		},
	}
	mergeConfig(config, cfg)
}

// configureFromFlags takes incoming flags and creates
// a configuration object from it.
func configureFromFlags(config *Configuration, f Flags) {
	dbSSL := config.Database.SSL
	if f.DisableDBSSL {
		dbSSL = false
	}

	dbDirect := config.Database.DirectConnect
	if f.DBDirectConnect {
		dbDirect = true
	}

	cfg := Configuration{
		Server{
			Host: f.Host,
			Port: f.Port,
			Security: Security{
				APIKey: f.APIKey,
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
			DirectConnect: dbDirect,
		},
	}
	mergeConfig(config, cfg)
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
