package config

import (
	"os"
	"testing"
)

func TestConfigureDefault(t *testing.T) {
	expectedPort := "3001"
	expectedAPIKey := ""
	expectedEncryptionKey := ""
	expectedDriver := "redis"
	expectedDBAddress := "localhost"
	expectedDBURI := ""
	expectedDB := "burnit"
	expectedDBUser := ""
	expectedDBPassword := ""
	expectedDBSSL := true
	exepctedDBDirectConnect := false

	var flags Flags
	config, err := Configure(flags)
	if err == nil {
		t.Fatalf("error in test, encryption key must be set: %v", err)
	}

	flags.EncryptionKey = "aabbcc"
	if config.Server.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Server.Port, expectedPort)
	}
	if config.Server.Security.APIKey != expectedAPIKey {
		t.Errorf("API Key was incorrect, got :%s, want: %s", config.Server.Security.APIKey, expectedAPIKey)
	}
	if config.Server.Security.Encryption.Key != expectedEncryptionKey {
		t.Errorf("Encryption key was incorrect, got: %s, want: %s", config.Server.Security.Encryption.Key, expectedEncryptionKey)
	}
	if config.Database.Driver != expectedDriver {
		t.Errorf("Driver weas incorrect, got: %s, want: %s", config.Database.Driver, expectedDriver)
	}
	if config.Database.Address != expectedDBAddress {
		t.Errorf("DB Address was incorrect, got: %s, want: %s", config.Database.Address, expectedDBAddress)
	}
	if config.Database.URI != expectedDBURI {
		t.Errorf("DB URI was incorrect, got: %s, want: %s", config.Database.URI, expectedDBURI)
	}
	if config.Database.Database != expectedDB {
		t.Errorf("DB was incorrect, got: %s, want: %s", config.Database.Database, expectedDB)
	}
	if config.Database.Username != expectedDBUser {
		t.Errorf("DB User was incorrect, got: %s, want: %s", config.Database.Username, expectedDBUser)
	}
	if config.Database.Password != expectedDBPassword {
		t.Errorf("DB Password was incorrect, got: %s, want: %s", config.Database.Password, expectedDBPassword)
	}
	if config.Database.SSL != expectedDBSSL {
		t.Errorf("DB SSL was incorrect, got: %t, want: %t", config.Database.SSL, expectedDBSSL)
	}
	if config.Database.DirectConnect != exepctedDBDirectConnect {
		t.Errorf("DB Direct Connect was incorrect, got: %t, want: %t", config.Database.DirectConnect, exepctedDBDirectConnect)
	}
}

func TestConfigureFromFile(t *testing.T) {
	expectedPort := "3003"
	expectedAPIKey := "aabbcc"
	expectedEncryptionKey := "secretstring"
	expectedDriver := "mongo"
	expectedDBAddress := "localhost:27017"
	expectedDBURI := "mongodb://localhost:27017"
	expectedDB := "burnit_db"
	expectedDBUser := "dbuser"
	expectedDBPassword := "dbpassword"
	expectedDBSSL := false
	exepctedDBDirectConnect := true
	configPath := "../test/config.yaml"

	config := &Configuration{}
	if err := configureFromFile(config, configPath); err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Server.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Server.Port, expectedPort)
	}
	if config.Server.Security.APIKey != expectedAPIKey {
		t.Errorf("API Key was incorrect, got :%s, want: %s", config.Server.Security.APIKey, expectedAPIKey)
	}
	if config.Server.Security.Encryption.Key != expectedEncryptionKey {
		t.Errorf("Encryption key was incorrect, got: %s, want: %s", config.Server.Security.Encryption.Key, expectedEncryptionKey)
	}
	if config.Database.Driver != expectedDriver {
		t.Errorf("Driver was incorrect, got: %s, want: %s", config.Database.Driver, expectedDriver)
	}
	if config.Database.Address != expectedDBAddress {
		t.Errorf("DB Address was incorrect, got: %s, want: %s", config.Database.Address, expectedDBAddress)
	}
	if config.Database.URI != expectedDBURI {
		t.Errorf("DB URI was incorrect, got: %s, want: %s", config.Database.URI, expectedDBURI)
	}
	if config.Database.Database != expectedDB {
		t.Errorf("DB was incorrect, got: %s, want: %s", config.Database.Database, expectedDB)
	}
	if config.Database.Username != expectedDBUser {
		t.Errorf("DB User was incorrect, got: %s, want: %s", config.Database.Username, expectedDBUser)
	}
	if config.Database.Password != expectedDBPassword {
		t.Errorf("DB Password was incorrect, got: %s, want: %s", config.Database.Password, expectedDBPassword)
	}
	if config.Database.SSL != expectedDBSSL {
		t.Errorf("DB SSL was incorrect, got: %t, want: %t", config.Database.SSL, expectedDBSSL)
	}
	if config.Database.DirectConnect != exepctedDBDirectConnect {
		t.Errorf("DB Direct Connect was incorrect, got: %t, want: %t", config.Database.DirectConnect, exepctedDBDirectConnect)
	}

	if err := configureFromFile(config, "nonexisting.yml"); err == nil {
		t.Errorf("should return error if file does not exist")
	}
}

func TestConfigureFromEnv(t *testing.T) {
	expectedPort := "3004"
	expectedAPIKey := "aabbcc"
	expectedEncryptionKey := "secretstring"
	expectedDBAddress := "localhost:27017"
	expectedDBURI := "mongodb://localhost:27017"
	expectedDB := "burnit"
	expectedDBUser := "dbuser"
	expectedDBPassword := "dbpassword"
	expectedDriver := "mongo"
	exepctedDBDirectConnect := true
	expectedDBSSL := false

	config := &Configuration{}
	os.Setenv("BURNITDB_LISTEN_PORT", expectedPort)
	os.Setenv("BURNITDB_API_KEY", expectedAPIKey)
	os.Setenv("BURNITDB_ENCRYPTION_KEY", expectedEncryptionKey)
	os.Setenv("DB_DRIVER", expectedDriver)
	os.Setenv("DB_HOST", expectedDBAddress)
	os.Setenv("DB_CONNECTION_URI", expectedDBURI)
	os.Setenv("DB", expectedDB)
	os.Setenv("DB_USER", expectedDBUser)
	os.Setenv("DB_PASSWORD", expectedDBPassword)
	os.Setenv("DB_SSL", "false")
	os.Setenv("DB_DIRECT_CONNECT", "true")

	configureFromEnv(config)

	if config.Server.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Server.Port, expectedPort)
	}
	if config.Server.Security.APIKey != expectedAPIKey {
		t.Errorf("API Key was incorrect, got :%s, want: %s", config.Server.Security.APIKey, expectedAPIKey)
	}
	if config.Server.Security.Encryption.Key != expectedEncryptionKey {
		t.Errorf("Encryption key was incorrect, got: %s, want: %s", config.Server.Security.Encryption.Key, expectedEncryptionKey)
	}
	if config.Database.Driver != expectedDriver {
		t.Errorf("Driver was incorrect, got: %s, want: %s", config.Database.Driver, expectedDriver)
	}
	if config.Database.Address != expectedDBAddress {
		t.Errorf("DB Address was incorrect, got: %s, want: %s", config.Database.Address, expectedDBAddress)
	}
	if config.Database.URI != expectedDBURI {
		t.Errorf("DB URI was incorrect, got: %s, want: %s", config.Database.URI, expectedDBURI)
	}
	if config.Database.Database != expectedDB {
		t.Errorf("DB was incorrect, got: %s, want: %s", config.Database.Database, expectedDB)
	}
	if config.Database.Username != expectedDBUser {
		t.Errorf("DB User was incorrect, got: %s, want: %s", config.Database.Username, expectedDBUser)
	}
	if config.Database.Password != expectedDBPassword {
		t.Errorf("DB Password was incorrect, got: %s, want: %s", config.Database.Password, expectedDBPassword)
	}
	if config.Database.SSL != expectedDBSSL {
		t.Errorf("DB SSL was incorrect, got: %t, want: %t", config.Database.SSL, expectedDBSSL)
	}
	if config.Database.DirectConnect != exepctedDBDirectConnect {
		t.Errorf("DB Direct Connect was incorrect, got: %t, want: %t", config.Database.DirectConnect, exepctedDBDirectConnect)
	}

	os.Setenv("BURNITDB_LISTEN_PORT", "")
	os.Setenv("BURNITDB_API_KEY", "")
	os.Setenv("BURNITDB_ENCRYPTION_KEY", "")
	os.Setenv("BURNITDB_HASH_METHOD", "")
	os.Setenv("DB_DRIVER", "")
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_CONNECTION_URI", "")
	os.Setenv("DB", "")
	os.Setenv("DB_USER", "")
	os.Setenv("DB_PASSWORD", "")
	os.Setenv("DB_SSL", "")
	os.Setenv("DB_DIRECT_CONNECT", "")
}

func TestConfigureFromFlags(t *testing.T) {
	expectedPort := "4001"
	expectedAPIKey := "ccaabb"
	expectedEncryptionKey := "stringsecret"
	expectedDBAddress := "localhost:6379"
	expectedDBURI := "localhost:6379"
	expectedDB := "burnit"
	expectedDBUser := "dbuser"
	expectedDBPassword := "dbpassword"
	expectedDriver := "redis"
	expectedDBSSL1 := false
	expectedDBSSL2 := true
	exepctedDBDirectConnect := true

	flags := Flags{
		Port:            expectedPort,
		APIKey:          expectedAPIKey,
		EncryptionKey:   expectedEncryptionKey,
		Driver:          expectedDriver,
		DBAddress:       expectedDBAddress,
		DBURI:           expectedDBURI,
		DB:              expectedDB,
		DBUser:          expectedDBUser,
		DBPassword:      expectedDBPassword,
		DisableDBSSL:    true,
		DBDirectConnect: true,
	}

	config := &Configuration{}
	configureFromFlags(config, flags)

	if config.Server.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Server.Port, expectedPort)
	}
	if config.Server.Security.APIKey != expectedAPIKey {
		t.Errorf("API Key was incorrect, got :%s, want: %s", config.Server.Security.APIKey, expectedAPIKey)
	}
	if config.Server.Security.Encryption.Key != expectedEncryptionKey {
		t.Errorf("Encryption key was incorrect, got: %s, want: %s", config.Server.Security.Encryption.Key, expectedEncryptionKey)
	}
	if config.Database.Driver != expectedDriver {
		t.Errorf("Driver was incorrect, got: %s, want: %s", config.Database.Driver, expectedDriver)
	}
	if config.Database.Address != expectedDBAddress {
		t.Errorf("DB Address was incorrect, got: %s, want: %s", config.Database.Address, expectedDBAddress)
	}
	if config.Database.URI != expectedDBURI {
		t.Errorf("DB URI was incorrect, got: %s, want: %s", config.Database.URI, expectedDBURI)
	}
	if config.Database.Database != expectedDB {
		t.Errorf("DB was incorrect, got: %s, want: %s", config.Database.Database, expectedDB)
	}
	if config.Database.Username != expectedDBUser {
		t.Errorf("DB User was incorrect, got: %s, want: %s", config.Database.Username, expectedDBUser)
	}
	if config.Database.Password != expectedDBPassword {
		t.Errorf("DB Password was incorrect, got: %s, want: %s", config.Database.Password, expectedDBPassword)
	}
	if config.Database.SSL != expectedDBSSL1 {
		t.Errorf("DB SSL was incorrect, got: %t, want: %t", config.Database.SSL, expectedDBSSL1)
	}
	if config.Database.DirectConnect != exepctedDBDirectConnect {
		t.Errorf("DB Direct Connect was incorrect, got: %t, want: %t", config.Database.DirectConnect, exepctedDBDirectConnect)
	}

	config.Database.SSL = true
	flags.DisableDBSSL = false
	configureFromFlags(config, flags)

	if config.Database.SSL != expectedDBSSL2 {
		t.Errorf("DB SSL was incorrect, got: %t, want: %t", config.Database.SSL, expectedDBSSL2)
	}
}

func TestConfigure(t *testing.T) {
	configPath := "../test/config.yaml"
	// Test default configuration.
	expectedPort := "3001"
	expectedAPIKey := ""
	expectedEncryptionKey := ""
	expectedDriver := "redis"
	expectedDBAddress := "localhost"
	expectedDBURI := ""
	expectedDB := "burnit"
	expectedDBUser := ""
	expectedDBPassword := ""
	expectedDBSSL := true
	exepctedDBDirectConnect := false

	var flags Flags
	config, err := Configure(flags)
	if err == nil {
		t.Fatalf("error in test, encryption key must be set: %v", err)
	}

	flags.EncryptionKey = "aabbcc"
	if config.Server.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Server.Port, expectedPort)
	}
	if config.Server.Security.APIKey != expectedAPIKey {
		t.Errorf("API Key was incorrect, got :%s, want: %s", config.Server.Security.APIKey, expectedAPIKey)
	}
	if config.Server.Security.Encryption.Key != expectedEncryptionKey {
		t.Errorf("Encryption key was incorrect, got: %s, want: %s", config.Server.Security.Encryption.Key, expectedEncryptionKey)
	}
	if config.Database.Driver != expectedDriver {
		t.Errorf("Driver weas incorrect, got: %s, want: %s", config.Database.Driver, expectedDriver)
	}
	if config.Database.Address != expectedDBAddress {
		t.Errorf("DB Address was incorrect, got: %s, want: %s", config.Database.Address, expectedDBAddress)
	}
	if config.Database.URI != expectedDBURI {
		t.Errorf("DB URI was incorrect, got: %s, want: %s", config.Database.URI, expectedDBURI)
	}
	if config.Database.Database != expectedDB {
		t.Errorf("DB was incorrect, got: %s, want: %s", config.Database.Database, expectedDB)
	}
	if config.Database.Username != expectedDBUser {
		t.Errorf("DB User was incorrect, got: %s, want: %s", config.Database.Username, expectedDBUser)
	}
	if config.Database.Password != expectedDBPassword {
		t.Errorf("DB Password was incorrect, got: %s, want: %s", config.Database.Password, expectedDBPassword)
	}
	if config.Database.SSL != expectedDBSSL {
		t.Errorf("DB SSL was incorrect, got: %t, want: %t", config.Database.SSL, expectedDBSSL)
	}
	if config.Database.DirectConnect != exepctedDBDirectConnect {
		t.Errorf("DB Direct Connect was incorrect, got: %t, want: %t", config.Database.DirectConnect, exepctedDBDirectConnect)
	}
	flags.EncryptionKey = ""

	// Test with file. Should override default.
	expectedPort = "3003"
	expectedAPIKey = "aabbcc"
	expectedEncryptionKey = "secretstring"
	expectedDriver = "mongo"
	expectedDBAddress = "localhost:27017"
	expectedDBURI = "mongodb://localhost:27017"
	expectedDB = "burnit_db"
	expectedDBUser = "dbuser"
	expectedDBPassword = "dbpassword"
	expectedDBSSL = false
	exepctedDBDirectConnect = true

	flags.ConfigPath = configPath
	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Server.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Server.Port, expectedPort)
	}
	if config.Server.Security.APIKey != expectedAPIKey {
		t.Errorf("API Key was incorrect, got :%s, want: %s", config.Server.Security.APIKey, expectedAPIKey)
	}
	if config.Server.Security.Encryption.Key != expectedEncryptionKey {
		t.Errorf("Encryption key was incorrect, got: %s, want: %s", config.Server.Security.Encryption.Key, expectedEncryptionKey)
	}
	if config.Database.Driver != expectedDriver {
		t.Errorf("Driver weas incorrect, got: %s, want: %s", config.Database.Driver, expectedDriver)
	}
	if config.Database.Address != expectedDBAddress {
		t.Errorf("DB Address was incorrect, got: %s, want: %s", config.Database.Address, expectedDBAddress)
	}
	if config.Database.URI != expectedDBURI {
		t.Errorf("DB URI was incorrect, got: %s, want: %s", config.Database.URI, expectedDBURI)
	}
	if config.Database.Database != expectedDB {
		t.Errorf("DB was incorrect, got: %s, want: %s", config.Database.Database, expectedDB)
	}
	if config.Database.Username != expectedDBUser {
		t.Errorf("DB User was incorrect, got: %s, want: %s", config.Database.Username, expectedDBUser)
	}
	if config.Database.Password != expectedDBPassword {
		t.Errorf("DB Password was incorrect, got: %s, want: %s", config.Database.Password, expectedDBPassword)
	}
	if config.Database.SSL != expectedDBSSL {
		t.Errorf("DB SSL was incorrect, got: %t, want: %t", config.Database.SSL, expectedDBSSL)
	}
	if config.Database.DirectConnect != exepctedDBDirectConnect {
		t.Errorf("DB Direct Connect was incorrect, got: %t, want: %t", config.Database.DirectConnect, exepctedDBDirectConnect)
	}

	// Test with environment variables. Should override file.
	expectedPort = "3004"
	expectedAPIKey = "aabbcc"
	expectedEncryptionKey = "secretstring"
	expectedDBAddress = "localhost:27017"
	expectedDBURI = "mongodb://localhost:27017"
	expectedDB = "burnit"
	expectedDBUser = "dbuser"
	expectedDBPassword = "dbpassword"
	expectedDriver = "mongo"
	expectedDBSSL = false
	exepctedDBDirectConnect = false

	os.Setenv("BURNITDB_LISTEN_PORT", expectedPort)
	os.Setenv("BURNITDB_API_KEY", expectedAPIKey)
	os.Setenv("BURNITDB_ENCRYPTION_KEY", expectedEncryptionKey)
	os.Setenv("DB_DRIVER", expectedDriver)
	os.Setenv("DB_HOST", expectedDBAddress)
	os.Setenv("DB_CONNECTION_URI", expectedDBURI)
	os.Setenv("DB", expectedDB)
	os.Setenv("DB_USER", expectedDBUser)
	os.Setenv("DB_PASSWORD", expectedDBPassword)
	os.Setenv("DB_SSL", "false")
	os.Setenv("DB_DIRECT_CONNECT", "false")

	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Server.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Server.Port, expectedPort)
	}
	if config.Server.Security.APIKey != expectedAPIKey {
		t.Errorf("API Key was incorrect, got :%s, want: %s", config.Server.Security.APIKey, expectedAPIKey)
	}
	if config.Server.Security.Encryption.Key != expectedEncryptionKey {
		t.Errorf("Encryption key was incorrect, got: %s, want: %s", config.Server.Security.Encryption.Key, expectedEncryptionKey)
	}
	if config.Database.Driver != expectedDriver {
		t.Errorf("Driver weas incorrect, got: %s, want: %s", config.Database.Driver, expectedDriver)
	}
	if config.Database.Address != expectedDBAddress {
		t.Errorf("DB Address was incorrect, got: %s, want: %s", config.Database.Address, expectedDBAddress)
	}
	if config.Database.URI != expectedDBURI {
		t.Errorf("DB URI was incorrect, got: %s, want: %s", config.Database.URI, expectedDBURI)
	}
	if config.Database.Database != expectedDB {
		t.Errorf("DB was incorrect, got: %s, want: %s", config.Database.Database, expectedDB)
	}
	if config.Database.Username != expectedDBUser {
		t.Errorf("DB User was incorrect, got: %s, want: %s", config.Database.Username, expectedDBUser)
	}
	if config.Database.Password != expectedDBPassword {
		t.Errorf("DB Password was incorrect, got: %s, want: %s", config.Database.Password, expectedDBPassword)
	}
	if config.Database.SSL != expectedDBSSL {
		t.Errorf("DB SSL was incorrect, got: %t, want: %t", config.Database.SSL, expectedDBSSL)
	}
	if config.Database.DirectConnect != exepctedDBDirectConnect {
		t.Errorf("DB Direct Connect was incorrect, got: %t, want: %t", config.Database.DirectConnect, exepctedDBDirectConnect)
	}

	// Test with flags. Should override file and envrionment variables.
	expectedPort = "4001"
	expectedAPIKey = "ccaabb"
	expectedEncryptionKey = "stringsecret"
	expectedDBAddress = "localhost:6379"
	expectedDBURI = "localhost:6379"
	expectedDB = "burnit"
	expectedDBUser = "dbuser"
	expectedDBPassword = "dbpassword"
	expectedDriver = "redis"
	expectedDBSSL1 := false
	expectedDBSSL2 := true
	exepctedDBDirectConnect1 := false

	flags = Flags{
		Port:            expectedPort,
		APIKey:          expectedAPIKey,
		EncryptionKey:   expectedEncryptionKey,
		Driver:          expectedDriver,
		DBAddress:       expectedDBAddress,
		DBURI:           expectedDBURI,
		DB:              expectedDB,
		DBUser:          expectedDBUser,
		DBPassword:      expectedDBPassword,
		DisableDBSSL:    true,
		DBDirectConnect: false,
	}

	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Server.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Server.Port, expectedPort)
	}
	if config.Server.Security.APIKey != expectedAPIKey {
		t.Errorf("API Key was incorrect, got :%s, want: %s", config.Server.Security.APIKey, expectedAPIKey)
	}
	if config.Server.Security.Encryption.Key != expectedEncryptionKey {
		t.Errorf("Encryption key was incorrect, got: %s, want: %s", config.Server.Security.Encryption.Key, expectedEncryptionKey)
	}
	if config.Database.Driver != expectedDriver {
		t.Errorf("Driver was incorrect, got: %s, want: %s", config.Database.Driver, expectedDriver)
	}
	if config.Database.Address != expectedDBAddress {
		t.Errorf("DB Address was incorrect, got: %s, want: %s", config.Database.Address, expectedDBAddress)
	}
	if config.Database.URI != expectedDBURI {
		t.Errorf("DB URI was incorrect, got: %s, want: %s", config.Database.URI, expectedDBURI)
	}
	if config.Database.Database != expectedDB {
		t.Errorf("DB was incorrect, got: %s, want: %s", config.Database.Database, expectedDB)
	}
	if config.Database.Username != expectedDBUser {
		t.Errorf("DB User was incorrect, got: %s, want: %s", config.Database.Username, expectedDBUser)
	}
	if config.Database.Password != expectedDBPassword {
		t.Errorf("DB Password was incorrect, got: %s, want: %s", config.Database.Password, expectedDBPassword)
	}
	if config.Database.SSL != expectedDBSSL1 {
		t.Errorf("DB SSL was incorrect, got: %t, want: %t", config.Database.SSL, expectedDBSSL1)
	}
	if config.Database.DirectConnect != exepctedDBDirectConnect1 {
		t.Errorf("DB Direct Connect was incorrect, got: %t, want: %t", config.Database.DirectConnect, exepctedDBDirectConnect1)
	}

	config.Database.SSL = true
	flags.DisableDBSSL = false
	configureFromFlags(config, flags)

	if config.Database.SSL != expectedDBSSL2 {
		t.Errorf("DB SSL was incorrect, got: %t, want: %t", config.Database.SSL, expectedDBSSL2)
	}

	os.Setenv("BURNITDB_LISTEN_PORT", "")
	os.Setenv("BURNITDB_API_KEY", "")
	os.Setenv("BURNITDB_ENCRYPTION_KEY", "")
	os.Setenv("BURNITDB_HASH_METHOD", "")
	os.Setenv("DB_DRIVER", "")
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_CONNECTION_URI", "")
	os.Setenv("DB", "")
	os.Setenv("DB_USER", "")
	os.Setenv("DB_PASSWORD", "")
	os.Setenv("DB_SSL", "")
}

func TestAddressFromMongoURI(t *testing.T) {
	var tests = []struct {
		uri string
	}{
		{uri: "mongodb://localhost:27017"},
		{uri: "mongodb://localhost:27017/?ssl=true"},
		{uri: "mongodb://user:pass@localhost:27017"},
		{uri: "mongodb://user:pass@localhost:27017/?ssl=true"},
	}

	expected := "localhost:27017"
	for _, test := range tests {
		addr := AddressFromMongoURI(test.uri)
		if addr != expected {
			t.Errorf("incorrect value, got: %s, want: %s", addr, expected)
		}
	}
}

func TestAddressFromRedisURI(t *testing.T) {
	var tests = []struct {
		uri string
	}{
		{uri: "localhost:6379"},
		{uri: "redis://localhost:6379"},
		{uri: "rediss://localhost:6379"},
		{uri: "localhost:6379,password=1234,ssl=true"},
		{uri: "redis://localhost:6379,password=1234,ssl=true"},
		{uri: "rediss://localhost:6379,password=1234,ssl=true"},
	}

	expected := "localhost:6379"
	for _, test := range tests {
		addr := AddressFromRedisURI(test.uri)
		if addr != expected {
			t.Errorf("incorrect value, got: %s, want: %s", addr, expected)
		}
	}
}
