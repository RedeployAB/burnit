package config

import (
	"flag"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	defaultListenPort = "3002"
)

// Flags is parsed flags.
type Flags struct {
	ConfigPath string
	Port       string
}

// ParseFlags runs flag.Parse and returns a flag object.
func ParseFlags() Flags {
	configPath := flag.String("config", "", "Path to configuration file")
	listenPort := flag.String("port", "", "Port to listen on")
	flag.Parse()

	return Flags{
		ConfigPath: *configPath,
		Port:       *listenPort,
	}
}

// Configuration represents a configuration.
type Configuration struct {
	Port string `yaml:"port"`
}

// Configure calls configureFromEnv or
// configureFromFile depending on the parameters
// passed in.
func Configure(f Flags) (*Configuration, error) {
	// Set default configuration.
	config := &Configuration{
		Port: defaultListenPort,
	}

	if len(f.ConfigPath) > 0 {
		if err := configureFromFile(config, f.ConfigPath); err != nil {
			return nil, err
		}
	}

	configureFromEnv(config)
	configureFromFlags(config, f)

	return config, nil
}

// mergeConfig merges a configuration with the base
// configuration.
func mergeConfig(config *Configuration, srcCfg Configuration) {
	if len(srcCfg.Port) > 0 {
		config.Port = srcCfg.Port
	}
}

// configureFromFile performs the necessary steps
// for server configuration from environment
// variables.
func configureFromFile(config *Configuration, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var cfg Configuration
	if err = yaml.Unmarshal(b, &config); err != nil {
		return err
	}
	// Merge configurations.
	mergeConfig(config, cfg)
	return nil
}

// configureFromEnv performs the necessary steps
// for server configuration from environment
// variables.
func configureFromEnv(config *Configuration) {
	cfg := Configuration{
		Port: os.Getenv("BURNITGEN_LISTEN_PORT"),
	}
	mergeConfig(config, cfg)
}

// configureFromFlags takes incoming flags and creates
// a configuration object from it.
func configureFromFlags(config *Configuration, f Flags) {
	cfg := Configuration{
		Port: f.Port,
	}
	mergeConfig(config, cfg)
}
