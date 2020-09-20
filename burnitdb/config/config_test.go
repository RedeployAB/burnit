package config

import (
	"os"
	"testing"
)

func TestConfigureDefault(t *testing.T) {
	expectedPort := "3001"
	expectedAPIKey := ""
	expectedHashMethod := "bcrypt"
	expectedEncryptionKey := ""
	expectedDriver := "redis"
	expectedDBAddress := "localhost"
	expectedDBURI := ""
	expectedDB := "burnitdb"
	expectedDBUser := ""
	expectedDBPassword := ""
	expectedDBSSL := true

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
	if config.Server.Security.HashMethod != expectedHashMethod {
		t.Errorf("Hash method was incorrect, got: %s, want: %s", config.Server.Security.HashMethod, expectedHashMethod)
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
}

func TestConfigureFromFile(t *testing.T) {
	expectedPort := "3003"
	expectedAPIKey := "aabbcc"
	expectedHashMethod := "bcrypt"
	expectedEncryptionKey := "secretstring"
	expectedDriver := "mongo"
	expectedDBAddress := "localhost:27017"
	expectedDBURI := "mongodb://localhost:27017"
	expectedDB := "burnit_db"
	expectedDBUser := "dbuser"
	expectedDBPassword := "dbpassword"
	expectedDBSSL := false
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
	if config.Server.Security.HashMethod != expectedHashMethod {
		t.Errorf("Hash method was incorrect, got: %s, want: %s", config.Server.Security.HashMethod, expectedHashMethod)
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

	if err := configureFromFile(config, "nonexisting.yml"); err == nil {
		t.Errorf("should return error if file does not exist")
	}
}

func TestConfigureFromEnv(t *testing.T) {
	expectedPort := "3004"
	expectedAPIKey := "aabbcc"
	expectedHashMethod := "bcrypt"
	expectedEncryptionKey := "secretstring"
	expectedDBAddress := "localhost:27017"
	expectedDBURI := "mongodb://localhost:27017"
	expectedDB := "burnitdb"
	expectedDBUser := "dbuser"
	expectedDBPassword := "dbpassword"
	expectedDriver := "mongo"
	expectedDBSSL := false

	config := &Configuration{}
	os.Setenv("BURNITDB_LISTEN_PORT", expectedPort)
	os.Setenv("BURNITDB_API_KEY", expectedAPIKey)
	os.Setenv("BURNITDB_ENCRYPTION_KEY", expectedEncryptionKey)
	os.Setenv("BURNITDB_HASH_METHOD", expectedHashMethod)
	os.Setenv("DB_DRIVER", expectedDriver)
	os.Setenv("DB_HOST", expectedDBAddress)
	os.Setenv("DB_CONNECTION_URI", expectedDBURI)
	os.Setenv("DB", expectedDB)
	os.Setenv("DB_USER", expectedDBUser)
	os.Setenv("DB_PASSWORD", expectedDBPassword)
	os.Setenv("DB_SSL", "false")

	configureFromEnv(config)

	if config.Server.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Server.Port, expectedPort)
	}
	if config.Server.Security.APIKey != expectedAPIKey {
		t.Errorf("API Key was incorrect, got :%s, want: %s", config.Server.Security.APIKey, expectedAPIKey)
	}
	if config.Server.Security.HashMethod != expectedHashMethod {
		t.Errorf("Hash method was incorrect, got: %s, want: %s", config.Server.Security.HashMethod, expectedHashMethod)
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

func TestConfigureFromFlags(t *testing.T) {
	expectedPort := "4001"
	expectedAPIKey := "ccaabb"
	expectedHashMethod := "md5"
	expectedEncryptionKey := "stringsecret"
	expectedDBAddress := "localhost:6379"
	expectedDBURI := "localhost:6379"
	expectedDB := "burnitdb"
	expectedDBUser := "dbuser"
	expectedDBPassword := "dbpassword"
	expectedDriver := "redis"
	expectedDBSSL1 := false
	expectedDBSSL2 := true

	flags := Flags{
		Port:          expectedPort,
		APIKey:        expectedAPIKey,
		HashMethod:    expectedHashMethod,
		EncryptionKey: expectedEncryptionKey,
		Driver:        expectedDriver,
		DBAddress:     expectedDBAddress,
		DBURI:         expectedDBURI,
		DB:            expectedDB,
		DBUser:        expectedDBUser,
		DBPassword:    expectedDBPassword,
		DisableDBSSL:  true,
	}

	config := &Configuration{}
	configureFromFlags(config, flags)

	if config.Server.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Server.Port, expectedPort)
	}
	if config.Server.Security.APIKey != expectedAPIKey {
		t.Errorf("API Key was incorrect, got :%s, want: %s", config.Server.Security.APIKey, expectedAPIKey)
	}
	if config.Server.Security.HashMethod != expectedHashMethod {
		t.Errorf("Hash method was incorrect, got: %s, want: %s", config.Server.Security.HashMethod, expectedHashMethod)
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
	expectedHashMethod := "bcrypt"
	expectedEncryptionKey := ""
	expectedDriver := "redis"
	expectedDBAddress := "localhost"
	expectedDBURI := ""
	expectedDB := "burnitdb"
	expectedDBUser := ""
	expectedDBPassword := ""
	expectedDBSSL := true

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
	if config.Server.Security.HashMethod != expectedHashMethod {
		t.Errorf("Hash method was incorrect, got: %s, want: %s", config.Server.Security.HashMethod, expectedHashMethod)
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
	flags.EncryptionKey = ""

	// Test with file. Should override default.
	expectedPort = "3003"
	expectedAPIKey = "aabbcc"
	expectedHashMethod = "bcrypt"
	expectedEncryptionKey = "secretstring"
	expectedDriver = "mongo"
	expectedDBAddress = "localhost:27017"
	expectedDBURI = "mongodb://localhost:27017"
	expectedDB = "burnit_db"
	expectedDBUser = "dbuser"
	expectedDBPassword = "dbpassword"
	expectedDBSSL = false

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
	if config.Server.Security.HashMethod != expectedHashMethod {
		t.Errorf("Hash method was incorrect, got: %s, want: %s", config.Server.Security.HashMethod, expectedHashMethod)
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

	// Test with environment variables. Should override file.
	expectedPort = "3004"
	expectedAPIKey = "aabbcc"
	expectedHashMethod = "bcrypt"
	expectedEncryptionKey = "secretstring"
	expectedDBAddress = "localhost:27017"
	expectedDBURI = "mongodb://localhost:27017"
	expectedDB = "burnitdb"
	expectedDBUser = "dbuser"
	expectedDBPassword = "dbpassword"
	expectedDriver = "mongo"
	expectedDBSSL = false

	os.Setenv("BURNITDB_LISTEN_PORT", expectedPort)
	os.Setenv("BURNITDB_API_KEY", expectedAPIKey)
	os.Setenv("BURNITDB_ENCRYPTION_KEY", expectedEncryptionKey)
	os.Setenv("BURNITDB_HASH_METHOD", expectedHashMethod)
	os.Setenv("DB_DRIVER", expectedDriver)
	os.Setenv("DB_HOST", expectedDBAddress)
	os.Setenv("DB_CONNECTION_URI", expectedDBURI)
	os.Setenv("DB", expectedDB)
	os.Setenv("DB_USER", expectedDBUser)
	os.Setenv("DB_PASSWORD", expectedDBPassword)
	os.Setenv("DB_SSL", "false")

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
	if config.Server.Security.HashMethod != expectedHashMethod {
		t.Errorf("Hash method was incorrect, got: %s, want: %s", config.Server.Security.HashMethod, expectedHashMethod)
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

	// Test with flags. Should override file and envrionment variables.
	expectedPort = "4001"
	expectedAPIKey = "ccaabb"
	expectedHashMethod = "md5"
	expectedEncryptionKey = "stringsecret"
	expectedDBAddress = "localhost:6379"
	expectedDBURI = "localhost:6379"
	expectedDB = "burnitdb"
	expectedDBUser = "dbuser"
	expectedDBPassword = "dbpassword"
	expectedDriver = "redis"
	expectedDBSSL1 := false
	expectedDBSSL2 := true

	flags = Flags{
		Port:          expectedPort,
		APIKey:        expectedAPIKey,
		HashMethod:    expectedHashMethod,
		EncryptionKey: expectedEncryptionKey,
		Driver:        expectedDriver,
		DBAddress:     expectedDBAddress,
		DBURI:         expectedDBURI,
		DB:            expectedDB,
		DBUser:        expectedDBUser,
		DBPassword:    expectedDBPassword,
		DisableDBSSL:  true,
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
	if config.Server.Security.HashMethod != expectedHashMethod {
		t.Errorf("Hash method was incorrect, got: %s, want: %s", config.Server.Security.HashMethod, expectedHashMethod)
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
