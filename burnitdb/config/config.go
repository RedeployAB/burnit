package config

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Server represents server part of configuration.
type Server struct {
	Port       string `yaml:"port"`
	DBAPIKey   string `yaml:"dbApiKey"`
	Passphrase string `yaml:"passphrase"`
}

// Database represents database part of configuration.
type Database struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSL      bool   `yaml:"ssl"`
	URI      string `yaml:"uri"`
}

// Configuration represents a configuration.
type Configuration struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

// Configure performs the necessary steps
// for server/app configuration.
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

	if len(config.Server.Passphrase) == 0 {
		return Configuration{}, errors.New("encryption passphrase must be set")
	}
	return config, nil
}

// configureFromEnv performs the necessary steps
// for server/app configuration from environment
// variables.
func configureFromEnv() Configuration {
	// Server variables.
	port := os.Getenv("BURNITDB_PORT")
	if len(port) == 0 {
		port = "3001"
	}
	dbAPIkey := os.Getenv("BURNITDB_API_KEY")
	passphrase := os.Getenv("BURNITDB_PASSPHRASE")

	// Database variables.
	dbHost := os.Getenv("DB_HOST")
	if len(dbHost) == 0 {
		dbHost = "localhost"
	}

	db := os.Getenv("DB")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	var dbSSL bool
	dbSSLStr := os.Getenv("DB_SSL")
	if len(dbSSLStr) == 0 {
		dbSSL = false
	} else {
		dbSSL, _ = strconv.ParseBool(dbSSLStr)
	}

	uri := os.Getenv("DB_CONNECTION_URI")

	config := Configuration{
		Server{
			Port:       port,
			DBAPIKey:   dbAPIkey,
			Passphrase: passphrase,
		},
		Database{
			Address:  dbHost,
			Database: db,
			Username: dbUser,
			Password: dbPassword,
			SSL:      dbSSL,
			URI:      uri,
		},
	}
	return config
}

// configureFromFile performs the necessary steps
// for server/app configuration from environment
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
		config.Server.Port = "3001"
	}
	if len(config.Database.Address) == 0 && len(config.Database.URI) == 0 {
		config.Database.URI = "mongodb://localhost"
	}
	return config, nil
}
