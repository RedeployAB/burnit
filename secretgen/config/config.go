package config

import "os"

// Config exports a configuration to be used by application.
var Config Configuration

func init() {
	Config = Configure()
}

// Server represents server part of configuration.
type Server struct {
	Port string
}

// Configuration represents a configuration.
type Configuration struct {
	Server
}

// Configure performs the necessary steps
// for server/app configuration.
func Configure() Configuration {
	port := os.Getenv("SECRET_GENERATOR_PORT")
	if port == "" {
		port = "3002"
	}

	config := Configuration{Server{Port: port}}

	return config
}
