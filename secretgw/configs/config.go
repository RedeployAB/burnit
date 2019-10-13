package configs

import "os"

// Server represents server part of configuration.
type Server struct {
	Port                 string
	GeneratorBaseURL     string
	GeneratorServicePath string
	DBBaseURL            string
	DBServicePath        string
	DBAPIKey             string
}

// Configuration represents a configuration.
type Configuration struct {
	Server
}

// Configure performs the necessary steps
// for server/app configuration.
func Configure() Configuration {
	port := os.Getenv("SECRET_GW_SERVICE_PORT")
	if port == "" {
		port = "3000"
	}

	genBaseURL := os.Getenv("SECRET_GENERATOR_SERVICE_BASE_URL")
	if genBaseURL == "" {
		genBaseURL = "http://localhost:3002"
	}
	genSvcPath := os.Getenv("SECRET_GENERATOR_SERVICE_PATH")
	if genSvcPath == "" {
		genSvcPath = "/api/v0/generate"
	}

	dbBaseURL := os.Getenv("SECRET_DB_SERVICE_BASE_URL")
	if dbBaseURL == "" {
		dbBaseURL = "http://localhost:3001"
	}
	dbSvcPath := os.Getenv("SECRET_DB_SERVICE_PATH")
	if dbSvcPath == "" {
		dbSvcPath = "/api/v0/secrets"
	}

	dbAPIKey := os.Getenv("SECRET_DB_SERVICE_API_KEY")

	config := Configuration{
		Server{
			Port:                 port,
			GeneratorBaseURL:     genBaseURL,
			GeneratorServicePath: genSvcPath,
			DBBaseURL:            dbBaseURL,
			DBServicePath:        dbSvcPath,
			DBAPIKey:             dbAPIKey,
		},
	}

	return config
}
