package config

import (
	"os"
	"testing"
)

func TestConfigureFromEnv(t *testing.T) {
	confDefault := configureFromEnv()
	if confDefault.Server.Port != "3001" {
		t.Errorf("default port value is incorrect, got %s, want: 3001", confDefault.Server.Port)
	}
	if confDefault.Server.DBAPIKey != "" {
		t.Errorf("default passphrase value is incorrect, got %s, want: \"\"", confDefault.Server.DBAPIKey)
	}
	if confDefault.Server.Passphrase != "" {
		t.Errorf("default passphrase value is incorrect, got %s, want: \"\"", confDefault.Server.Passphrase)
	}
	if confDefault.Database.Address != "localhost" {
		t.Errorf("default address value is incorrect, got %s, want: localhost", confDefault.Database.Address)
	}
	if confDefault.Database.Database != "" {
		t.Errorf("default database value is incorrect, got %s, want: \"\"", confDefault.Database.Database)
	}
	if confDefault.Database.Username != "" {
		t.Errorf("default username value is incorrect, got %s, want: \"\"", confDefault.Database.Username)
	}
	if confDefault.Database.Password != "" {
		t.Errorf("default password value is incorrect, got %s, want: \"\"", confDefault.Database.Password)
	}
	if confDefault.Database.SSL != false {
		t.Errorf("default ssl value is incorrect, got %v, want: false", confDefault.Database.SSL)
	}
	if confDefault.Database.URI != "" {
		t.Errorf("default uri value is incorrect, got %v, want: \"\"", confDefault.Database.URI)
	}

	os.Setenv("BURNITDB_PORT", "6000")
	os.Setenv("BURNITDB_API_KEY", "aabbcc")
	os.Setenv("BURNITDB_PASSPHRASE", "secretstring")
	os.Setenv("DB_HOST", "localhost:27017")
	os.Setenv("DB", "burnit_db")
	os.Setenv("DB_USER", "dbuser")
	os.Setenv("DB_PASSWORD", "dbpassword")
	os.Setenv("DB_SSL", "true")
	os.Setenv("DB_CONNECTION_URI", "mongodb://localhost:27017")
	confEnv := configureFromEnv()
	if confEnv.Server.Port != "6000" {
		t.Errorf("port value is incorrect, got %s, want: 6000", confEnv.Server.Port)
	}
	if confEnv.Server.DBAPIKey != "aabbcc" {
		t.Errorf("passphrase value is incorrect, got %s, want: aabbcc", confEnv.Server.DBAPIKey)
	}
	if confEnv.Server.Passphrase != "secretstring" {
		t.Errorf("passphrase value is incorrect, got %s, want: secretstring", confEnv.Server.Passphrase)
	}

	if confEnv.Database.Address != "localhost:27017" {
		t.Errorf("address value is incorrect, got %s, want: localhost:27017", confEnv.Database.Address)
	}
	if confEnv.Database.Database != "burnit_db" {
		t.Errorf("database value is incorrect, got %s, want: burnit_db", confEnv.Database.Database)
	}
	if confEnv.Database.Username != "dbuser" {
		t.Errorf("username value is incorrect, got %s, want: dbuser", confEnv.Database.Username)
	}
	if confEnv.Database.Password != "dbpassword" {
		t.Errorf("password value is incorrect, got %s, want: dbpassword", confEnv.Database.Password)
	}
	if confEnv.Database.SSL != true {
		t.Errorf("ssl value is incorrect, got %v, want: true", confEnv.Database.SSL)
	}
	if confEnv.Database.URI != "mongodb://localhost:27017" {
		t.Errorf("uri value is incorrect, got %v, want: mongodb://localhost:27017", confEnv.Database.URI)
	}
}

func TestConfigureFromFile(t *testing.T) {
	configPath := "../test/config.yaml"
	conf, err := configureFromFile(configPath)
	if err != nil {
		t.Errorf("%v", err)
	}

	if conf.Server.Port != "3003" {
		t.Errorf("port value is incorrect, got %s, want: 3003", conf.Server.Port)
	}
	if conf.Server.DBAPIKey != "aabbcc" {
		t.Errorf("passphrase value is incorrect, got %s, want: aabbcc", conf.Server.DBAPIKey)
	}
	if conf.Server.Passphrase != "secretstring" {
		t.Errorf("passphrase value is incorrect, got %s, want: secretstring", conf.Server.Passphrase)
	}

	if conf.Database.Address != "localhost:27017" {
		t.Errorf("address value is incorrect, got %s, want: localhost:27017", conf.Database.Address)
	}
	if conf.Database.Database != "burnit_db" {
		t.Errorf("database value is incorrect, got %s, want: burnit_db", conf.Database.Database)
	}
	if conf.Database.Username != "dbuser" {
		t.Errorf("username value is incorrect, got %s, want: dbuser", conf.Database.Username)
	}
	if conf.Database.Password != "dbpassword" {
		t.Errorf("password value is incorrect, got %s, want: dbpassword", conf.Database.Password)
	}
	if conf.Database.SSL != true {
		t.Errorf("ssl value is incorrect, got %v, want: true", conf.Database.SSL)
	}
	if conf.Database.URI != "mongodb://localhost:27017" {
		t.Errorf("uri value is incorrect, got %v, want: mongodb://localhost:27017", conf.Database.URI)
	}
}

func TestConfigure(t *testing.T) {
	// Test Configure from environment.
	_, err := Configure("")
	if err != nil {
		t.Errorf("error: %v", err)
	}
	// Test Configure from file.
	_, err = Configure("../test/config.yaml")
	if err != nil {
		t.Errorf("eerror: %v", err)
	}
	// Handle whene no configuration exists.
	_, err = Configure("../test/nofile.yml")
	if err == nil {
		t.Errorf("incorrect, should return an error")
	}
}
