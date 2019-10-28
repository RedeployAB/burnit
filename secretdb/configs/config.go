package configs

import (
	"os"
	"strconv"
)

// Server represents server part of configuration.
type Server struct {
	Port       string
	DBAPIKey   string
	Passphrase string
}

// Database represents database part of configuration.
type Database struct {
	Address  string
	Database string
	Username string
	Password string
	SSL      bool
	URI      string
}

// Configuration represents a configuration.
type Configuration struct {
	Server   Server
	Database Database
}

// Configure performs the necessary steps
// for server/app configuration.
func Configure() Configuration {
	// Server variables.
	port := os.Getenv("SECRET_DB_PORT")
	if port == "" {
		port = "3001"
	}
	dbAPIkey := os.Getenv("SECRET_DB_API_KEY")
	passphrase := os.Getenv("SECRET_DB_PASSPHRASE")

	// Database variables.
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
