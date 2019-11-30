package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Server represents server part of configuration.
type Server struct {
	Port                 string `yaml:"port"`
	GeneratorBaseURL     string `yaml:"generatorBaseUrl"`
	GeneratorServicePath string `yaml:"generatorServicePath"`
	DBBaseURL            string `yaml:"dbBaseUrl"`
	DBServicePath        string `yaml:"dbServicePath"`
	DBAPIKey             string `yaml:"dbApiKey"`
}

// Configuration represents a configuration.
type Configuration struct {
	Server `yaml:"server"`
}

// Configure calls configureFromEnvor
// configureFromFile depending on the parameters
// passed  in.
func Configure(path string) (Configuration, error) {
	var config Configuration
	var err error
	if path == "" {
		config = configureFromEnv()
	} else {
		config, err = configureFromFile(path)
	}
	return config, err
}

// configureFromEnv performs the necessary steps
// for server/app configuration from environment
// variables.
func configureFromEnv() Configuration {
	port := os.Getenv("SECRET_GW_PORT")
	if port == "" {
		port = "3000"
	}

	genBaseURL := os.Getenv("SECRET_GEN_BASE_URL")
	if genBaseURL == "" {
		genBaseURL = "http://localhost:3002"
	}
	genSvcPath := os.Getenv("SECRET_GEN_PATH")
	if genSvcPath == "" {
		genSvcPath = "/api/v0/generate"
	}

	dbBaseURL := os.Getenv("SECRET_DB_BASE_URL")
	if dbBaseURL == "" {
		dbBaseURL = "http://localhost:3001"
	}
	dbSvcPath := os.Getenv("SECRET_DB_PATH")
	if dbSvcPath == "" {
		dbSvcPath = "/api/v0/secrets"
	}

	dbAPIKey := os.Getenv("SECRET_DB_API_KEY")

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

	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
