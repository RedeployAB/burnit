package configs

import "os"

// Configuration represents a configuration.
type Configuration struct {
	Port string
}

// Configure performs the necessary steps
// for server/app configuration.
func Configure() Configuration {
	port := os.Getenv("SECRET_GEN_PORT")
	if port == "" {
		port = "3002"
	}

	config := Configuration{Port: port}

	return config
}
