package config

import "os"

// Config exports a configuration to be used by application.
var Config Configuration

func init() {
	Config = Configure()
}

// Server represents server part of configuration.
type Server struct {
	Port             string
	GeneratorBaseURL string
	DBBaseURL        string
	DBAPIKey         string
}

// Configuration represents a configuration.
type Configuration struct {
	Server
}

// Configure performs the necessary steps
// for server/app configuration.
func Configure() Configuration {
	port := os.Getenv("SECRET_API_SERVICE_PORT")
	if port == "" {
		port = "3000"
	}

	genBaseURL := os.Getenv("SECRET_GENERATOR_SERVICE_BASE_URL")
	if genBaseURL == "" {
		genBaseURL = "http://localhost:3002"
	}

	dbBaseURL := os.Getenv("SECRET_DB_SERVICE_BASE_URL")
	if dbBaseURL == "" {
		dbBaseURL = "http://localhost:3001"
	}

	dbAPIKey := os.Getenv("SECRET_DB_SERVICE_API_KEY")

	config := Configuration{
		Server{
			Port:             port,
			GeneratorBaseURL: genBaseURL,
			DBBaseURL:        dbBaseURL,
			DBAPIKey:         dbAPIKey,
		},
	}

	return config
}
