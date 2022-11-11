package config

import "flag"

// Flags is parsed flags.
type Flags struct {
	ConfigPath      string
	Host            string
	Port            string
	TLSCertificate  string
	TLSKey          string
	CORSOrigin      string
	APIKey          string
	EncryptionKey   string
	Driver          string
	DBAddress       string
	DBURI           string
	DB              string
	DBUser          string
	DBPassword      string
	DisableDBSSL    bool
	DBDirectConnect bool
}

// ParseFlags runs flag.Parse and returns a flag object.
func ParseFlags() Flags {
	configPath := flag.String("config", "", "Path to configuration file")
	listenHost := flag.String("host", "", "Host to listen on")
	listenPort := flag.String("port", "", "Port to listen on")
	tlsCert := flag.String("tls-certificate", "", "Path to TLS certificate file")
	tlsKey := flag.String("tls-key", "", "Path to TLS key file")
	corsOrigin := flag.String("cors-origin", "", "Enable and set CORS origin")
	apiKey := flag.String("api-key", "", "API key for database endpoints")
	encryptionKey := flag.String("encryption-key", "", "Encryption key for secrets in database")
	driver := flag.String("driver", "", "Database driver for storage of secrets: redis|mongo")
	dbAddress := flag.String("db-address", "", "Host name and port for database")
	dbURI := flag.String("db-uri", "", "URI for database connection")
	db := flag.String("db", "", "Database name")
	dbUser := flag.String("db-user", "", "User for database connections")
	dbPassword := flag.String("db-password", "", "Password for user for database connections")
	disableDBSSL := flag.Bool("disable-db-ssl", false, "Disable SSL for database connections")
	dbDirectConnect := flag.Bool("db-direct-connect", false, "Enable direct connect (mongodb only)")
	flag.Parse()

	return Flags{
		ConfigPath:      *configPath,
		Host:            *listenHost,
		Port:            *listenPort,
		TLSCertificate:  *tlsCert,
		TLSKey:          *tlsKey,
		CORSOrigin:      *corsOrigin,
		APIKey:          *apiKey,
		EncryptionKey:   *encryptionKey,
		Driver:          *driver,
		DBAddress:       *dbAddress,
		DBURI:           *dbURI,
		DB:              *db,
		DBUser:          *dbUser,
		DBPassword:      *dbPassword,
		DisableDBSSL:    *disableDBSSL,
		DBDirectConnect: *dbDirectConnect,
	}
}
