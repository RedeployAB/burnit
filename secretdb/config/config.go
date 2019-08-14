package config

import (
	"os"
	"strconv"
)

// Config exports a configuration to be used by application.
var Config Configuration

func init() {
	Config = Configure()
}

// Server represents server part of configuration.
type Server struct {
	Port       string
	Passphrase string
}

// Database represents database part of configuration.
type Database struct {
	Address  string
	Database string
	Username string
	Password string
	SSL      bool
}

// Configuration represents a configuration.
type Configuration struct {
	Server
	Database
}

// Configure performs the necessary steps
// for server/app configuration.
func Configure() Configuration {
	port := os.Getenv("SECRET_DB_SERVICE_PORT")
	if port == "" {
		port = "3001"
	}

	passphrase := os.Getenv("SECRET_DB_PASSPHRASE")

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	db := os.Getenv("DB")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	var dbSSL bool
	dbSSLStr := os.Getenv("DB_SSL")
	if dbSSLStr == "" {
		dbSSL = false
	} else {
		dbSSL, _ = strconv.ParseBool(dbSSLStr)
	}

	config := Configuration{
		Server{
			Port:       port,
			Passphrase: passphrase,
		},
		Database{
			Address:  dbHost,
			Database: db,
			Username: dbUser,
			Password: dbPassword,
			SSL:      dbSSL,
		},
	}

	return config
}
