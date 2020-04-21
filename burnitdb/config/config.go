package config

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	defaultListenPort = "3001"
	defaultHashMethod = "md5"
	defaultDB         = "burnitdb"
	defaultDBURI      = "mongodb://localhost"
	defaultDriver     = "mongo"
)

// Server defines server part of configuration.
type Server struct {
	Port     string   `yaml:"port"`
	Security Security `yaml:"security"`
}

// Security defines security part of server configuration.
type Security struct {
	APIKey     string     `yaml:"apiKey"`
	Encryption Encryption `yaml:"encryption"`
	HashMethod string     `yaml:"hashMethod"`
}

// Encryption defines encryption pat of security configuration.
type Encryption struct {
	Key string `yaml:"key"`
}

// Database represents database part of configuration.
type Database struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSL      bool   `yaml:"ssl"`
	URI      string `yaml:"uri"`
	Driver   string `yaml:"driver"`
}

// Configuration represents a configuration.
type Configuration struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

// Configure performs the necessary steps
// for server configuration.
func Configure(path string) (Configuration, error) {
	var config Configuration
	var err error
	if len(path) == 0 {
		config = configureFromEnv()
	} else {
		config, err = configureFromFile(path)
		if err != nil {
			return config, err
		}
	}

	if len(config.Server.Security.Encryption.Key) == 0 {
		return Configuration{}, errors.New("encryption key must be set")
	}
	return config, nil
}

// configureFromEnv performs the necessary steps
// for server configuration from environment
// variables.
func configureFromEnv() Configuration {
	// Server variables.
	port := os.Getenv("BURNITDB_LISTEN_PORT")
	if len(port) == 0 {
		port = defaultListenPort
	}
	apiKey := os.Getenv("BURNITDB_API_KEY")
	encryptionKey := os.Getenv("BURNITDB_ENCRYPTION_KEY")
	hashMethod := strings.ToLower(os.Getenv("BURNITDB_HASH_METHOD"))
	if len(hashMethod) == 0 {
		hashMethod = defaultHashMethod
	}
	// Database variables.
	dbHost := os.Getenv("DB_HOST")
	if len(dbHost) == 0 {
		dbHost = "localhost"
	}

	db := os.Getenv("DB")
	if len(db) == 0 {
		db = defaultDB
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	var dbSSL bool
	dbSSLStr := os.Getenv("DB_SSL")
	if len(dbSSLStr) == 0 {
		dbSSL = false
	} else {
		dbSSL, _ = strconv.ParseBool(dbSSLStr)
	}

	driver := strings.ToLower(os.Getenv("DB_DRIVER"))
	if len(driver) == 0 {
		driver = "mongo"
	}

	uri := os.Getenv("DB_CONNECTION_URI")

	config := Configuration{
		Server{
			Port: port,
			Security: Security{
				APIKey: apiKey,
				Encryption: Encryption{
					Key: encryptionKey,
				},
				HashMethod: hashMethod,
			},
		},
		Database{
			Address:  dbHost,
			Database: db,
			Username: dbUser,
			Password: dbPassword,
			SSL:      dbSSL,
			URI:      uri,
			Driver:   driver,
		},
	}
	return config
}

// configureFromFile performs the necessary steps
// for server configuration from environment
// variables.
func configureFromFile(path string) (Configuration, error) {
	var config = Configuration{}
	f, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return config, err
	}

	if err = yaml.Unmarshal(b, &config); err != nil {
		return config, err
	}

	if len(config.Server.Port) == 0 {
		config.Server.Port = defaultListenPort
	}
	if len(config.Server.Security.HashMethod) == 0 {
		config.Server.Security.HashMethod = defaultHashMethod
	} else {
		config.Server.Security.HashMethod = strings.ToLower(config.Server.Security.HashMethod)
	}

	if len(config.Database.Database) == 0 {
		config.Database.Database = defaultDB
	}

	if len(config.Database.Address) == 0 && len(config.Database.URI) == 0 {
		config.Database.URI = defaultDBURI
	}

	if len(config.Database.Driver) == 0 {
		config.Database.Driver = defaultDriver
	} else {
		config.Database.Driver = strings.ToLower(config.Database.Driver)
	}
	return config, nil
}
