package config

import (
	"flag"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	defaultListenPort           = "3000"
	defaultGeneratorBaseURL     = "http://localhost:3002"
	defaultGeneratorServicePath = "/generate"
	defaultDBBaseURL            = "http://localhost:3001"
	defaultDBServicePath        = "/secrets"
	defaultDBAPIKey             = ""
)

// Flags is parsed flags.
type Flags struct {
	ConfigPath           string
	Port                 string
	GeneratorBaseURL     string
	GeneratorServicePath string
	DBBaseURL            string
	DBServicePath        string
	DBAPIKey             string
}

// ParseFlags runs flag.Parse and returns a flag object.
func ParseFlags() Flags {
	configPath := flag.String("config", "", "Path to configuration file")
	listenPort := flag.String("port", "", "Port to listen on")
	generatorBaseURL := flag.String("generator-base-url", "", "Base URL to generator service (burnitgen)")
	generatorServicePath := flag.String("generator-service-path", "", "Path to generator service endpoint (burnitgen)")
	dbBaseURL := flag.String("db-base-url", "", "Base URL to DB service (burnitdb)")
	dbServicePath := flag.String("db-service-path", "", "Path to DB service endpoint (burnitdb)")
	dbAPIKey := flag.String("db-api-key", "", "API Key to DB service")
	flag.Parse()

	return Flags{
		ConfigPath:           *configPath,
		Port:                 *listenPort,
		GeneratorBaseURL:     *generatorBaseURL,
		GeneratorServicePath: *generatorServicePath,
		DBBaseURL:            *dbBaseURL,
		DBServicePath:        *dbServicePath,
		DBAPIKey:             *dbAPIKey,
	}
}

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
func Configure(f Flags) (*Configuration, error) {
	// Set default configuration.
	config := &Configuration{
		Server{
			Port:                 defaultListenPort,
			GeneratorBaseURL:     defaultGeneratorBaseURL,
			GeneratorServicePath: defaultGeneratorServicePath,
			DBBaseURL:            defaultDBBaseURL,
			DBServicePath:        defaultDBServicePath,
			DBAPIKey:             defaultDBAPIKey,
		},
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
	if len(srcCfg.GeneratorBaseURL) > 0 {
		config.GeneratorBaseURL = srcCfg.GeneratorBaseURL
	}
	if len(srcCfg.GeneratorServicePath) > 0 {
		config.GeneratorServicePath = srcCfg.GeneratorServicePath
	}
	if len(srcCfg.DBBaseURL) > 0 {
		config.DBBaseURL = srcCfg.DBBaseURL
	}
	if len(srcCfg.DBServicePath) > 0 {
		config.DBServicePath = srcCfg.DBServicePath
	}
	if len(srcCfg.DBAPIKey) > 0 {
		config.DBAPIKey = srcCfg.DBAPIKey
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
	if err = yaml.Unmarshal(b, &cfg); err != nil {
		return err
	}

	mergeConfig(config, cfg)
	return nil
}

// configureFromEnv performs the necessary steps
// for server configuration from environment
// variables.
func configureFromEnv(config *Configuration) {
	cfg := Configuration{
		Server{
			Port:                 os.Getenv("BURNITGW_LISTEN_PORT"),
			GeneratorBaseURL:     os.Getenv("BURNITGEN_BASE_URL"),
			GeneratorServicePath: os.Getenv("BURNITGEN_PATH"),
			DBBaseURL:            os.Getenv("BURNITDB_BASE_URL"),
			DBServicePath:        os.Getenv("BURNITDB_PATH"),
			DBAPIKey:             os.Getenv("BURNITDB_API_KEY"),
		},
	}
	mergeConfig(config, cfg)
}

// configureFromFlags takes incoming flags and creates
// a configuration object from it.
func configureFromFlags(config *Configuration, f Flags) {
	cfg := Configuration{
		Server{
			Port:                 f.Port,
			GeneratorBaseURL:     f.GeneratorBaseURL,
			GeneratorServicePath: f.GeneratorServicePath,
			DBBaseURL:            f.DBBaseURL,
			DBServicePath:        f.DBServicePath,
			DBAPIKey:             f.DBAPIKey,
		},
	}
	mergeConfig(config, cfg)
}
