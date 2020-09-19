package config

import (
	"os"
	"testing"
)

func TestConfigureDefault(t *testing.T) {
	expectedPort := "3000"
	expectedGenAddress := "http://localhost:3002"
	expectedGenSvcPath := "/generate"
	expectedDBAddress := "http://localhost:3001"
	expectedDBSvcPath := "/secrets"
	expectedDBAPIKey := ""

	var flags Flags
	config, err := Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorAddress != expectedGenAddress {
		t.Errorf("Generator Address is incorrect, got: %s, want: %s", config.GeneratorAddress, expectedGenAddress)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBAddress != expectedDBAddress {
		t.Errorf("DB Address is incorrect, got: %s, want: %s", config.DBAddress, expectedDBAddress)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}
}

func TestConfigureFromFile(t *testing.T) {
	expectedPort := "3003"
	expectedGenAddress := "http://localhost:3003"
	expectedGenSvcPath := "/v1/generate"
	expectedDBAddress := "http://localhost:3003"
	expectedDBSvcPath := "/v1/secrets"
	expectedDBAPIKey := "aabbcc"
	configPath := "../test/config.yaml"

	config := &Configuration{}
	if err := configureFromFile(config, configPath); err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorAddress != expectedGenAddress {
		t.Errorf("Generator Address is incorrect, got: %s, want: %s", config.GeneratorAddress, expectedGenAddress)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBAddress != expectedDBAddress {
		t.Errorf("DB Address is incorrect, got: %s, want: %s", config.DBAddress, expectedDBAddress)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}

	if err := configureFromFile(config, "nonexisting.yml"); err == nil {
		t.Errorf("should return error if file does not exist")
	}
}

func TestConfigureFromEnv(t *testing.T) {
	expectedPort := "3003"
	expectedGenAddress := "http://someurl:3002"
	expectedGenSvcPath := "/v1/generate"
	expectedDBAddress := "http://someurl:3001"
	expectedDBSvcPath := "/v1/secrets"
	expectedDBAPIKey := "AAAA"

	config := &Configuration{}
	os.Setenv("BURNITGW_LISTEN_PORT", expectedPort)
	os.Setenv("BURNITGEN_ADDRESS", expectedGenAddress)
	os.Setenv("BURNITGEN_PATH", expectedGenSvcPath)
	os.Setenv("BURNITDB_ADDRESS", expectedDBAddress)
	os.Setenv("BURNITDB_PATH", expectedDBSvcPath)
	os.Setenv("BURNITDB_API_KEY", expectedDBAPIKey)
	configureFromEnv(config)

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorAddress != expectedGenAddress {
		t.Errorf("Generator Address is incorrect, got: %s, want: %s", config.GeneratorAddress, expectedGenAddress)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBAddress != expectedDBAddress {
		t.Errorf("DB Address is incorrect, got: %s, want: %s", config.DBAddress, expectedDBAddress)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}
	os.Setenv("BURNITGW_LISTEN_PORT", "")
	os.Setenv("BURNITGEN_ADDRESS", "")
	os.Setenv("BURNITGEN_PATH", "")
	os.Setenv("BURNITDB_ADDRESS", "")
	os.Setenv("BURNITDB_PATH", "")
	os.Setenv("BURNITDB_API_KEY", "")
}

func TestConfigureFromFlags(t *testing.T) {
	expectedPort := "4000"
	expectedGenAddress := "http://someurl:4002"
	expectedGenSvcPath := "/v1/generate"
	expectedDBAddress := "http://someurl:4003"
	expectedDBSvcPath := "/v1/secrets"
	expectedDBAPIKey := "ccaabb"

	flags := Flags{
		Port:                 expectedPort,
		GeneratorAddress:     expectedGenAddress,
		GeneratorServicePath: expectedGenSvcPath,
		DBAddress:            expectedDBAddress,
		DBServicePath:        expectedDBSvcPath,
		DBAPIKey:             expectedDBAPIKey,
	}

	config := &Configuration{}
	configureFromFlags(config, flags)

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorAddress != expectedGenAddress {
		t.Errorf("Generator Address is incorrect, got: %s, want: %s", config.GeneratorAddress, expectedGenAddress)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBAddress != expectedDBAddress {
		t.Errorf("DB Address is incorrect, got: %s, want: %s", config.DBAddress, expectedDBAddress)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}
}

func TestConfigure(t *testing.T) {
	configPath := "../test/config.yaml"
	// Test default configuration.
	expectedPort := "3000"
	expectedGenAddress := "http://localhost:3002"
	expectedGenSvcPath := "/generate"
	expectedDBAddress := "http://localhost:3001"
	expectedDBSvcPath := "/secrets"
	expectedDBAPIKey := ""

	var flags Flags
	config, err := Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorAddress != expectedGenAddress {
		t.Errorf("Generator Address is incorrect, got: %s, want: %s", config.GeneratorAddress, expectedGenAddress)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBAddress != expectedDBAddress {
		t.Errorf("DB Address is incorrect, got: %s, want: %s", config.DBAddress, expectedDBAddress)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}

	// Test with file. Should override default.
	expectedPort = "3003"
	expectedGenAddress = "http://localhost:3003"
	expectedGenSvcPath = "/v1/generate"
	expectedDBAddress = "http://localhost:3003"
	expectedDBSvcPath = "/v1/secrets"
	expectedDBAPIKey = "aabbcc"

	flags.ConfigPath = configPath
	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorAddress != expectedGenAddress {
		t.Errorf("Generator Address is incorrect, got: %s, want: %s", config.GeneratorAddress, expectedGenAddress)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBAddress != expectedDBAddress {
		t.Errorf("DB Address is incorrect, got: %s, want: %s", config.DBAddress, expectedDBAddress)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}

	// Test with environment variables. Should override file.
	expectedPort = "3003"
	expectedGenAddress = "http://someurl:3002"
	expectedGenSvcPath = "/v1/generate"
	expectedDBAddress = "http://someurl:3001"
	expectedDBSvcPath = "/v1/secrets"
	expectedDBAPIKey = "AAAA"

	os.Setenv("BURNITGW_LISTEN_PORT", expectedPort)
	os.Setenv("BURNITGEN_ADDRESS", expectedGenAddress)
	os.Setenv("BURNITGEN_PATH", expectedGenSvcPath)
	os.Setenv("BURNITDB_ADDRESS", expectedDBAddress)
	os.Setenv("BURNITDB_PATH", expectedDBSvcPath)
	os.Setenv("BURNITDB_API_KEY", expectedDBAPIKey)

	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorAddress != expectedGenAddress {
		t.Errorf("Generator Address is incorrect, got: %s, want: %s", config.GeneratorAddress, expectedGenAddress)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBAddress != expectedDBAddress {
		t.Errorf("DB Address is incorrect, got: %s, want: %s", config.DBAddress, expectedDBAddress)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}

	// Test with flags. Should override file and envrionment variables.
	expectedPort = "4000"
	expectedGenAddress = "http://someurl:4002"
	expectedGenSvcPath = "/v1/generate"
	expectedDBAddress = "http://someurl:4003"
	expectedDBSvcPath = "/v1/secrets"
	expectedDBAPIKey = "ccaabb"

	flags = Flags{
		ConfigPath:           configPath,
		Port:                 expectedPort,
		GeneratorAddress:     expectedGenAddress,
		GeneratorServicePath: expectedGenSvcPath,
		DBAddress:            expectedDBAddress,
		DBServicePath:        expectedDBSvcPath,
		DBAPIKey:             expectedDBAPIKey,
	}

	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorAddress != expectedGenAddress {
		t.Errorf("Generator Address is incorrect, got: %s, want: %s", config.GeneratorAddress, expectedGenAddress)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBAddress != expectedDBAddress {
		t.Errorf("DB Address is incorrect, got: %s, want: %s", config.DBAddress, expectedDBBaseURL)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}

	os.Setenv("BURNITGW_LISTEN_PORT", "")
	os.Setenv("BURNITGEN_ADDRESS", "")
	os.Setenv("BURNITGEN_PATH", "")
	os.Setenv("BURNITDB_ADDRESS", "")
	os.Setenv("BURNITDB_PATH", "")
	os.Setenv("BURNITDB_API_KEY", "")
}
