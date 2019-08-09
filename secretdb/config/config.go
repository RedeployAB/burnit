package config

import (
	"os"
	"strconv"
)

// Server represents server part of configuration.
type Server struct {
	Port             string
	GeneratorBaseURL string
	GeneratorPath    string
}

// Database represents database part of configuration.
type Database struct {
	Address  string
	Database string
	Username string
	Password string
	SSL      bool
}

// Config represents a configuration.
type Config struct {
	Server
	Database
}

// Configure performs the necessary steps
// for server/app configuration.
func Configure() Config {
	port := os.Getenv("SECRET_DB_PORT")
	if port == "" {
		port = "4000"
	}

	generatorBaseURL := os.Getenv("SECRET_GENERATOR_BASE_URL")
	if generatorBaseURL == "" {
		generatorBaseURL = "http://localhost:3000"
	}

	generatorPath := os.Getenv("SECRET_GENERATOR_REQUEST_PATH")
	if generatorPath == "" {
		generatorPath = "/api/v1/secret"
	}

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

	config := Config{
		Server{
			Port:             port,
			GeneratorBaseURL: generatorBaseURL,
			GeneratorPath:    generatorPath,
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
