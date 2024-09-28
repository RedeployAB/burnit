package config

import (
	"bytes"
	"flag"
	"strconv"
	"time"
)

// flags contains the flags.
type flags struct {
	configPath             string
	host                   string
	port                   int
	tlsCertFile            string
	tlsKeyFile             string
	encryptionKey          string
	corsOrigin             string
	timeout                time.Duration
	databaseConnStr        string
	databaseAddr           string
	database               string
	databaseUser           string
	databasePass           string
	databaseTimeout        time.Duration
	databaseConnectTimeout time.Duration
	databaseEnableTLS      *bool
}

// parseFlags parses the flags.
func parseFlags(args []string) (flags, string, error) {
	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	var buf bytes.Buffer
	fs.SetOutput(&buf)

	var f flags
	var enableDBTLS boolFlag

	fs.StringVar(&f.configPath, "config-path", "", "Optional. Path to a configuration file. Defaults to: "+defaultConfigPath+".")
	fs.StringVar(&f.host, "host", "", "Optional. Host to listen on. Defaults to: "+defaultListenHost+".")
	fs.IntVar(&f.port, "port", 0, "Optional. Port to listen on. Defaults to: "+strconv.Itoa(defaultListenPort)+".")
	fs.StringVar(&f.tlsCertFile, "tls-cert-file", "", "Optional. TLS certificate file.")
	fs.StringVar(&f.tlsKeyFile, "tls-key-file", "", "Optional. TLS key file.")
	fs.StringVar(&f.encryptionKey, "encryption-key", "", "Optional. Default encryption key for the secrets service.")
	fs.StringVar(&f.corsOrigin, "cors-origin", "", "Optional. CORS origin.")
	fs.DurationVar(&f.timeout, "timeout", 0, "Optional. Default timeout for the service. Defaults to: "+defaultSecretServiceTimeout.String()+".")
	fs.StringVar(&f.databaseConnStr, "database-connection-string", "", "Optional. Connection string for the database.")
	fs.StringVar(&f.databaseAddr, "database-address", "", "Optional. Address for the database.")
	fs.StringVar(&f.database, "database", "", "Optional. Database name.")
	fs.StringVar(&f.databaseUser, "database-user", "", "Optional. Database username.")
	fs.StringVar(&f.databasePass, "database-password", "", "Optional. Database password.")
	fs.DurationVar(&f.databaseTimeout, "database-timeout", 0, "Optional. Timeout for the database. Defaults to: "+defaultDatabaseTimeout.String()+".")
	fs.DurationVar(&f.databaseConnectTimeout, "database-connect-timeout", 0, "Optional. Connect timeout for the database. Defaults to: "+defaultDatabaseConnectTimeout.String()+".")
	fs.Var(&enableDBTLS, "database-enable-tls", "Optional. Enable TLS for the database. Defaults to true.")

	if err := fs.Parse(args); err != nil {
		return f, buf.String(), err
	}

	if enableDBTLS.isSet {
		f.databaseEnableTLS = &enableDBTLS.value
	}

	return f, buf.String(), nil
}

// boolFlag is a flag for bool values that keeps track if it was set.
type boolFlag struct {
	value bool
	isSet bool
}

// Set sets the value of the boolFlag.
func (f *boolFlag) Set(value string) error {
	var v bool
	if value == "true" {
		v = true
	}
	f.value = v
	f.isSet = true
	return nil
}

// String returns the string representation of the boolFlag.
func (f *boolFlag) String() string {
	if f.value {
		return "true"
	}
	return "false"
}
