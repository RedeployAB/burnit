package config

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	defaultListenPort = "3001"
	defaultHashMethod = "md5"
	defaultAPIKey     = ""
	defaultDriver     = "mongo"
	defaultAddress    = "localhost" // Change when "default" driver is updated.
	defaultDBURI      = ""          // Change when "default" driver is updated.
	defaultDB         = "burnitdb"
	defaultDBSSL      = true
)

// Flags is parsed flags.
type Flags struct {
	ConfigPath    string
	Port          string
	APIKey        string
	HashMethod    string
	EncryptionKey string
	Driver        string
	DBAddress     string
	DBURI         string
	DB            string
	DBUser        string
	DBPassword    string
	DisableDBSSL  bool
}

// ParseFlags runs flag.Parse and returns a flag object.
func ParseFlags() Flags {
	configPath := flag.String("config", "", "Path to configuration file")
	listenPort := flag.String("port", "", "Port to listen on")
	apiKey := flag.String("api-key", "", "API key for database endpoints")
	hashMethod := flag.String("hash-method", "", "Hash method for passphrase protected secrets")
	encryptionKey := flag.String("encryption-key", "", "Encryption key for secrets in database")
	driver := flag.String("driver", "", "Database driver for storage of secrets: mongo|redis")
	dbAddress := flag.String("db-address", "", "Host name and port for database")
	dbURI := flag.String("db-uri", "", "URI for database connection")
	db := flag.String("db", "", "Database name")
	dbUser := flag.String("db-user", "", "User for database connections")
	dbPassword := flag.String("db-password", "", "Password for user for database connections")
	disableDBSSL := flag.Bool("disable-db-ssl", false, "Disable SSL for database connections")
	flag.Parse()

	return Flags{
		ConfigPath:    *configPath,
		Port:          *listenPort,
		APIKey:        *apiKey,
		HashMethod:    *hashMethod,
		EncryptionKey: *encryptionKey,
		Driver:        *driver,
		DBAddress:     *dbAddress,
		DBURI:         *dbURI,
		DB:            *db,
		DBUser:        *dbUser,
		DBPassword:    *dbPassword,
		DisableDBSSL:  *disableDBSSL,
	}
}

// Server defines server part of configuration.
type Server struct {
	Port     string   `yaml:"port"`
	Security Security `yaml:"security"`
}

// Security defines security part of server configuration.
type Security struct {
	APIKey     string     `yaml:"apiKey"`
	HashMethod string     `yaml:"hashMethod"`
	Encryption Encryption `yaml:"encryption"`
}

// Encryption defines encryption pat of security configuration.
type Encryption struct {
	Key string `yaml:"key"`
}

// Database represents database part of configuration.
type Database struct {
	Driver   string `yaml:"driver"`
	Address  string `yaml:"address"`
	URI      string `yaml:"uri"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	SSL      bool   `yaml:"ssl"`
}

// Configuration represents a configuration.
type Configuration struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}

// Configure performs the necessary steps
// for server configuration.
func Configure(f Flags) (*Configuration, error) {
	config := &Configuration{
		Server{
			Port: defaultListenPort,
			Security: Security{
				APIKey: defaultAPIKey,
				Encryption: Encryption{
					Key: "",
				},
				HashMethod: defaultHashMethod,
			},
		},
		Database{
			Driver:   defaultDriver,
			Address:  defaultAddress,
			URI:      defaultDBURI,
			Database: defaultDB,
			Username: "",
			Password: "",
			SSL:      defaultDBSSL,
		},
	}

	if len(f.ConfigPath) > 0 {
		if err := configureFromFile(config, f.ConfigPath); err != nil {
			return nil, err
		}
	}

	configureFromEnv(config)
	configureFromFlags(config, f)

	if len(config.Server.Security.Encryption.Key) == 0 {
		return config, errors.New("encryption key must be set")
	}

	if config.Database.Driver == "redis" {
		re := regexp.MustCompile(`:\d+$`)
		if !re.MatchString(config.Database.Address) {
			config.Database.Address += ":6379"
		}
	}
	return config, nil
}

func mergeConfig(config *Configuration, srcCfg Configuration) {
	if len(srcCfg.Server.Port) > 0 {
		config.Server.Port = srcCfg.Server.Port
	}
	if len(srcCfg.Server.Security.APIKey) > 0 {
		config.Server.Security.APIKey = srcCfg.Server.Security.APIKey
	}
	if len(srcCfg.Server.Security.HashMethod) > 0 {
		config.Server.Security.HashMethod = srcCfg.Server.Security.HashMethod
	}
	if len(srcCfg.Server.Security.Encryption.Key) > 0 {
		config.Server.Security.Encryption.Key = srcCfg.Server.Security.Encryption.Key
	}
	if len(srcCfg.Database.Driver) > 0 {
		config.Database.Driver = strings.ToLower(srcCfg.Database.Driver)
	}
	if len(srcCfg.Database.Address) > 0 {
		config.Database.Address = srcCfg.Database.Address
	}
	if len(srcCfg.Database.URI) > 0 {
		config.Database.URI = srcCfg.Database.URI
	}
	if len(srcCfg.Database.Database) > 0 {
		config.Database.Database = srcCfg.Database.Database
	}
	if len(srcCfg.Database.Username) > 0 {
		config.Database.Username = srcCfg.Database.Username
	}
	if len(srcCfg.Database.Password) > 0 {
		config.Database.Password = srcCfg.Database.Password
	}
	config.Database.SSL = srcCfg.Database.SSL
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
	var dbSSL bool
	dbSSLStr := os.Getenv("DB_SSL")
	if len(dbSSLStr) == 0 {
		dbSSL = config.Database.SSL
	} else {
		dbSSL, _ = strconv.ParseBool(dbSSLStr)
	}

	cfg := Configuration{
		Server{
			Port: os.Getenv("BURNITDB_LISTEN_PORT"),
			Security: Security{
				APIKey: os.Getenv("BURNITDB_API_KEY"),
				Encryption: Encryption{
					Key: os.Getenv("BURNITDB_ENCRYPTION_KEY"),
				},
				HashMethod: strings.ToLower(os.Getenv("BURNITDB_HASH_METHOD")),
			},
		},
		Database{
			Driver:   strings.ToLower(os.Getenv("DB_DRIVER")),
			Address:  os.Getenv("DB_HOST"),
			URI:      os.Getenv("DB_CONNECTION_URI"),
			Database: os.Getenv("DB"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			SSL:      dbSSL,
		},
	}
	mergeConfig(config, cfg)
}

// configureFromFlags takes incoming flags and creates
// a configuration object from it.
func configureFromFlags(config *Configuration, f Flags) {
	dbSSL := config.Database.SSL
	if f.DisableDBSSL == true {
		dbSSL = false
	}

	cfg := Configuration{
		Server{
			Port: f.Port,
			Security: Security{
				APIKey: f.APIKey,
				Encryption: Encryption{
					Key: f.EncryptionKey,
				},
				HashMethod: f.HashMethod,
			},
		},
		Database{
			Driver:   f.Driver,
			Address:  f.DBAddress,
			URI:      f.DBURI,
			Database: f.DB,
			Username: f.DBUser,
			Password: f.DBPassword,
			SSL:      dbSSL,
		},
	}
	mergeConfig(config, cfg)
}
