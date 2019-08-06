package config

import (
	"os"
)

// Server represents server part of configuration.
type Server struct {
	Port string
}

// Config represents a configuration.
type Config struct {
	Server
}

// Configure performs the necessary steps
// for server/app configuration.
func Configure() Config {
	port := os.Getenv("SECRET_GENERATOR_PORT")
	if port == "" {
		port = "3000"
	}

	config := Config{Server{Port: port}}

	return config
}
