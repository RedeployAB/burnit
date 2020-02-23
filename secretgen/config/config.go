package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Configuration represents a configuration.
type Configuration struct {
	Port string `yaml:"port"`
}

// Configure calls configureFromEnv or
// configureFromFile depending on the parameters
// passed in.
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
// for server/app configuration from environment
// variables.
func configureFromEnv() Configuration {
	port := os.Getenv("BURNITGEN_PORT")
	if len(port) == 0 {
		port = "3002"
	}
	return Configuration{Port: port}
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
