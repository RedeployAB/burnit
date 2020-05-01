package config

import (
	"os"
	"testing"
)

func TestConfigureDefault(t *testing.T) {
	expectedPort := "3000"
	expectedGenBaseURL := "http://localhost:3002"
	expectedGenSvcPath := "/api/generate"
	expectedDBBaseURL := "http://localhost:3001"
	expectedDBSvcPath := "/api/secrets"
	expectedDBAPIKey := ""

	var flags Flags
	config, err := Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got: %s, want: %s", config.GeneratorBaseURL, expectedGenBaseURL)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got: %s, want: %s", config.DBBaseURL, expectedDBBaseURL)
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
	expectedGenBaseURL := "http://localhost:3003"
	expectedGenSvcPath := "/api/v1/generate"
	expectedDBBaseURL := "http://localhost:3003"
	expectedDBSvcPath := "/api/v1/secrets"
	expectedDBAPIKey := "aabbcc"
	configPath := "../test/config.yaml"

	config := &Configuration{}
	if err := configureFromFile(config, configPath); err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got: %s, want: %s", config.GeneratorBaseURL, expectedGenBaseURL)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got: %s, want: %s", config.DBBaseURL, expectedDBBaseURL)
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
	expectedGenBaseURL := "http://someurl:3002"
	expectedGenSvcPath := "/api/v1/generate"
	expectedDBBaseURL := "http://someurl:3001"
	expectedDBSvcPath := "/api/v1/secrets"
	expectedDBAPIKey := "AAAA"

	config := &Configuration{}
	os.Setenv("BURNITGW_LISTEN_PORT", expectedPort)
	os.Setenv("BURNITGEN_BASE_URL", expectedGenBaseURL)
	os.Setenv("BURNITGEN_PATH", expectedGenSvcPath)
	os.Setenv("BURNITDB_BASE_URL", expectedDBBaseURL)
	os.Setenv("BURNITDB_PATH", expectedDBSvcPath)
	os.Setenv("BURNITDB_API_KEY", expectedDBAPIKey)
	configureFromEnv(config)

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got: %s, want: %s", config.GeneratorBaseURL, expectedGenBaseURL)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got: %s, want: %s", config.DBBaseURL, expectedDBBaseURL)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}
	os.Setenv("BURNITGW_LISTEN_PORT", "")
	os.Setenv("BURNITGEN_BASE_URL", "")
	os.Setenv("BURNITGEN_PATH", "")
	os.Setenv("BURNITDB_BASE_URL", "")
	os.Setenv("BURNITDB_PATH", "")
	os.Setenv("BURNITDB_API_KEY", "")
}

func TestConfigureFromFlags(t *testing.T) {
	expectedPort := "4000"
	expectedGenBaseURL := "http://someurl:4002"
	expectedGenSvcPath := "/api/v1/generate"
	expectedDBBaseURL := "http://someurl:4003"
	expectedDBSvcPath := "/api/v1/secrets"
	expectedDBAPIKey := "ccaabb"

	flags := Flags{
		Port:                 expectedPort,
		GeneratorBaseURL:     expectedGenBaseURL,
		GeneratorServicePath: expectedGenSvcPath,
		DBBaseURL:            expectedDBBaseURL,
		DBServicePath:        expectedDBSvcPath,
		DBAPIKey:             expectedDBAPIKey,
	}

	config := &Configuration{}
	configureFromFlags(config, flags)

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got: %s, want: %s", config.GeneratorBaseURL, expectedGenBaseURL)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got: %s, want: %s", config.DBBaseURL, expectedDBBaseURL)
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
	// Test defaul t configuration.
	expectedPort := "3000"
	expectedGenBaseURL := "http://localhost:3002"
	expectedGenSvcPath := "/api/generate"
	expectedDBBaseURL := "http://localhost:3001"
	expectedDBSvcPath := "/api/secrets"
	expectedDBAPIKey := ""

	var flags Flags
	config, err := Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got: %s, want: %s", config.GeneratorBaseURL, expectedGenBaseURL)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got: %s, want: %s", config.DBBaseURL, expectedDBBaseURL)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}

	// Test with file. Should override default.
	expectedPort = "3003"
	expectedGenBaseURL = "http://localhost:3003"
	expectedGenSvcPath = "/api/v1/generate"
	expectedDBBaseURL = "http://localhost:3003"
	expectedDBSvcPath = "/api/v1/secrets"
	expectedDBAPIKey = "aabbcc"

	flags.ConfigPath = configPath
	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got: %s, want: %s", config.GeneratorBaseURL, expectedGenBaseURL)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got: %s, want: %s", config.DBBaseURL, expectedDBBaseURL)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}

	// Test with environment variables. Should override file.
	expectedPort = "3003"
	expectedGenBaseURL = "http://someurl:3002"
	expectedGenSvcPath = "/api/v1/generate"
	expectedDBBaseURL = "http://someurl:3001"
	expectedDBSvcPath = "/api/v1/secrets"
	expectedDBAPIKey = "AAAA"

	os.Setenv("BURNITGW_LISTEN_PORT", expectedPort)
	os.Setenv("BURNITGEN_BASE_URL", expectedGenBaseURL)
	os.Setenv("BURNITGEN_PATH", expectedGenSvcPath)
	os.Setenv("BURNITDB_BASE_URL", expectedDBBaseURL)
	os.Setenv("BURNITDB_PATH", expectedDBSvcPath)
	os.Setenv("BURNITDB_API_KEY", expectedDBAPIKey)

	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	if config.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got: %s, want: %s", config.GeneratorBaseURL, expectedGenBaseURL)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got: %s, want: %s", config.DBBaseURL, expectedDBBaseURL)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}

	// Test with flags. Should override file and envrionment variables.
	expectedPort = "4000"
	expectedGenBaseURL = "http://someurl:4002"
	expectedGenSvcPath = "/api/v1/generate"
	expectedDBBaseURL = "http://someurl:4003"
	expectedDBSvcPath = "/api/v1/secrets"
	expectedDBAPIKey = "ccaabb"

	flags = Flags{
		ConfigPath:           configPath,
		Port:                 expectedPort,
		GeneratorBaseURL:     expectedGenBaseURL,
		GeneratorServicePath: expectedGenSvcPath,
		DBBaseURL:            expectedDBBaseURL,
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
	if config.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got: %s, want: %s", config.GeneratorBaseURL, expectedGenBaseURL)
	}
	if config.GeneratorServicePath != expectedGenSvcPath {
		t.Errorf("Generator Service Path is incorrect, got: %s, want: %s", config.GeneratorServicePath, expectedGenSvcPath)
	}
	if config.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got: %s, want: %s", config.DBBaseURL, expectedDBBaseURL)
	}
	if config.DBServicePath != expectedDBSvcPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", config.DBServicePath, expectedDBSvcPath)
	}
	if config.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got :%s, want: %s", config.DBAPIKey, expectedDBAPIKey)
	}

	os.Setenv("BURNITGW_LISTEN_PORT", "")
	os.Setenv("BURNITGEN_BASE_URL", "")
	os.Setenv("BURNITGEN_PATH", "")
	os.Setenv("BURNITDB_BASE_URL", "")
	os.Setenv("BURNITDB_PATH", "")
	os.Setenv("BURNITDB_API_KEY", "")
}
