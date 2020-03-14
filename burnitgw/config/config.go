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
	if len(path) == 0 {
		config = configureFromEnv()
	} else {
		config, err = configureFromFile(path)
		if err != nil {
			return config, err
		}
	}
	return config, nil
}

// configureFromEnv performs the necessary steps
// for server configuration from environment
// variables.
func configureFromEnv() Configuration {
	port := os.Getenv("BURNITGW_LISTEN_PORT")
	if len(port) == 0 {
		port = "3000"
	}

	genBaseURL := os.Getenv("BURNITGEN_BASE_URL")
	if len(genBaseURL) == 0 {
		genBaseURL = "http://localhost:3002"
	}
	genSvcPath := os.Getenv("BURNITGEN_PATH")
	if len(genSvcPath) == 0 {
		genSvcPath = "/api/generate"
	}

	dbBaseURL := os.Getenv("BURNITDB_BASE_URL")
	if len(dbBaseURL) == 0 {
		dbBaseURL = "http://localhost:3001"
	}
	dbSvcPath := os.Getenv("BURNITDB_PATH")
	if len(dbSvcPath) == 0 {
		dbSvcPath = "/api/secrets"
	}

	dbAPIKey := os.Getenv("BURNITDB_API_KEY")

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
// for server configuration from environment
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
