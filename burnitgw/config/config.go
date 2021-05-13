package config

import (
	"crypto/tls"
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	defaultListenPort           = "3000"
	defaultGeneratorAddress     = "http://localhost:3002"
	defaultGeneratorServicePath = "/secret"
	defaultDBAddress            = "http://localhost:3001"
	defaultDBServicePath        = "/secrets"
	defaultDBAPIKey             = ""
	defaultTLSCert              = ""
	defaultTLSKey               = ""
)

// Flags is parsed flags.
type Flags struct {
	ConfigPath           string
	Port                 string
	GeneratorAddress     string
	GeneratorServicePath string
	DBAddress            string
	DBServicePath        string
	DBAPIKey             string
	TLSCert              string
	TLSKey               string
	CORSOrigin           string
}

// ParseFlags runs flag.Parse and returns a flag object.
func ParseFlags() Flags {
	configPath := flag.String("config", "", "Path to configuration file")
	listenPort := flag.String("port", "", "Port to listen on")
	generatorAddress := flag.String("generator-address", "", "Address to generator service (burnitgen)")
	generatorServicePath := flag.String("generator-service-path", "", "Path to generator service endpoint (burnitgen)")
	dbAddress := flag.String("db-address", "", "Address to DB service (burnitdb)")
	dbServicePath := flag.String("db-service-path", "", "Path to DB service endpoint (burnitdb)")
	dbAPIKey := flag.String("db-api-key", "", "API Key to DB service")
	tlsCert := flag.String("tls-certificate", "", "Path to TLS certificate file")
	tlsKey := flag.String("tls-key", "", "Path to TLS key file")
	corsOrigin := flag.String("cors-origin", "", "Enable CORS and set origin")

	flag.Parse()
	return Flags{
		ConfigPath:           *configPath,
		Port:                 *listenPort,
		GeneratorAddress:     *generatorAddress,
		GeneratorServicePath: *generatorServicePath,
		DBAddress:            *dbAddress,
		DBServicePath:        *dbServicePath,
		DBAPIKey:             *dbAPIKey,
		TLSCert:              *tlsCert,
		TLSKey:               *tlsKey,
		CORSOrigin:           *corsOrigin,
	}
}

// Server represents server part of configuration.
type Server struct {
	Port                 string `yaml:"port"`
	GeneratorAddress     string `yaml:"generatorAddress"`
	GeneratorServicePath string `yaml:"generatorServicePath"`
	DBAddress            string `yaml:"dbAddress"`
	DBServicePath        string `yaml:"dbServicePath"`
	DBAPIKey             string `yaml:"dbApiKey"`
	TLS                  `yaml:"tls"`
	CORS                 `yaml:"cors"`
}

// TLS represents TLS part of server configuration.
type TLS struct {
	Certificate string `yaml:"certificate"`
	Key         string `yaml:"key"`
}

// CORS represent CORS part of server configuration.
type CORS struct {
	Origin string `yaml:"origin,omitempty"`
}

// Configuration represents a configuration.
type Configuration struct {
	Server `yaml:"server"`
}

// Configure calls configureFromEnvor
// configureFromFile depending on the parameters
// passed  in.
func Configure(f Flags) (*Configuration, error) {
	config := &Configuration{
		Server{
			Port:                 defaultListenPort,
			GeneratorAddress:     defaultGeneratorAddress,
			GeneratorServicePath: defaultGeneratorServicePath,
			DBAddress:            defaultDBAddress,
			DBServicePath:        defaultDBServicePath,
			DBAPIKey:             defaultDBAPIKey,
			TLS: TLS{
				Certificate: defaultTLSCert,
				Key:         defaultTLSKey,
			},
		},
	}

	if len(f.ConfigPath) > 0 {
		if err := configureFromFile(config, f.ConfigPath); err != nil {
			return nil, err
		}
	}

	configureFromEnv(config)
	configureFromFlags(config, f)

	if !strings.HasPrefix(config.GeneratorAddress, "http://") && !strings.HasPrefix(config.GeneratorAddress, "https://") {
		config.GeneratorAddress = "http://" + config.GeneratorAddress
	}

	if !strings.HasPrefix(config.DBAddress, "http://") && !strings.HasPrefix(config.DBAddress, "https://") {
		config.DBAddress = "http://" + config.DBAddress
	}

	return config, nil
}

// mergeConfig merges a configuration with the base
// configuration.
func mergeConfig(config *Configuration, srcCfg Configuration) {
	if len(srcCfg.Port) > 0 {
		config.Port = srcCfg.Port
	}
	if len(srcCfg.GeneratorAddress) > 0 {
		config.GeneratorAddress = srcCfg.GeneratorAddress
	}
	if len(srcCfg.GeneratorServicePath) > 0 {
		config.GeneratorServicePath = srcCfg.GeneratorServicePath
	}
	if len(srcCfg.DBAddress) > 0 {
		config.DBAddress = srcCfg.DBAddress
	}
	if len(srcCfg.DBServicePath) > 0 {
		config.DBServicePath = srcCfg.DBServicePath
	}
	if len(srcCfg.DBAPIKey) > 0 {
		config.DBAPIKey = srcCfg.DBAPIKey
	}
	if len(srcCfg.TLS.Certificate) > 0 {
		config.TLS.Certificate = srcCfg.TLS.Certificate
	}
	if len(srcCfg.TLS.Key) > 0 {
		config.TLS.Key = srcCfg.TLS.Key
	}
	if len(srcCfg.CORS.Origin) > 0 {
		config.CORS.Origin = srcCfg.CORS.Origin
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
			GeneratorAddress:     os.Getenv("BURNITGEN_ADDRESS"),
			GeneratorServicePath: os.Getenv("BURNITGEN_PATH"),
			DBAddress:            os.Getenv("BURNITDB_ADDRESS"),
			DBServicePath:        os.Getenv("BURNITDB_PATH"),
			DBAPIKey:             os.Getenv("BURNITDB_API_KEY"),
			TLS: TLS{
				Certificate: os.Getenv("BURNITGW_TLS_CERTIFICATE"),
				Key:         os.Getenv("BURNITGW_TLS_KEY"),
			},
			CORS: CORS{
				Origin: os.Getenv("BURNITGW_CORS_ORIGIN"),
			},
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
			GeneratorAddress:     f.GeneratorAddress,
			GeneratorServicePath: f.GeneratorServicePath,
			DBAddress:            f.DBAddress,
			DBServicePath:        f.DBServicePath,
			DBAPIKey:             f.DBAPIKey,
			TLS: TLS{
				Certificate: f.TLSCert,
				Key:         f.TLSKey,
			},
			CORS: CORS{
				Origin: f.CORSOrigin,
			},
		},
	}
	mergeConfig(config, cfg)
}

func NewTLSConfig() *tls.Config {
	return &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
		},
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
}

func CORSHeaders() map[string][]string {
	return map[string][]string{
		"Access-Control-Allow-Headers": {"Content-Type"},
		"Access-Control-Allow-Methods": {"GET", "POST", "OPTIONS"},
	}
}
