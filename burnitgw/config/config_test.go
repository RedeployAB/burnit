package config

import (
	"os"
	"testing"
)

func TestConfigureFromEnv(t *testing.T) {
	confDefault := configureFromEnv()
	expectedPort := "3000"
	if confDefault.Port != expectedPort {
		t.Errorf("default port value is incorrect, got: %s, want: %s", confDefault.Port, expectedPort)
	}
	expectedGenBaseURL := "http://localhost:3002"
	if confDefault.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got: %s, want: %s", confDefault.GeneratorBaseURL, expectedGenBaseURL)
	}
	expectedGenPath := "/api/generate"
	if confDefault.GeneratorServicePath != expectedGenPath {
		t.Errorf("Generator Service  Path is incorrect, got: %s, want: %s", confDefault.GeneratorServicePath, expectedGenPath)
	}
	expectedDBBaseURL := "http://localhost:3001"
	if confDefault.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got: %s, want: %s", confDefault.DBBaseURL, expectedDBBaseURL)
	}
	expectedDBPath := "/api/secrets"
	if confDefault.DBServicePath != expectedDBPath {
		t.Errorf("DB Service Path is incorrect, got: %s, want: %s", confDefault.DBServicePath, expectedDBPath)
	}
	expectedDBAPIKey := ""
	if confDefault.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got: %s, want: %s", confDefault.DBAPIKey, expectedDBAPIKey)
	}

	os.Setenv("BURNITGW_LISTEN_PORT", "5000")
	os.Setenv("BURNITGEN_BASE_URL", "http://someurl:3000")
	os.Setenv("BURNITGEN_PATH", "/api/v1/generate")
	os.Setenv("BURNITDB_BASE_URL", "http://someurl:3001")
	os.Setenv("BURNITDB_PATH", "/api/v1/secrets")
	os.Setenv("BURNITDB_API_KEY", "AAAA")
	confEnv := configureFromEnv()

	expectedPort = "5000"
	if confEnv.Port != expectedPort {
		t.Errorf("Port value is incorrect, got: %s, want: %s", confEnv.Port, expectedPort)
	}
	expectedGenBaseURL = "http://someurl:3000"
	if confEnv.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got %s, want: %s", confEnv.GeneratorBaseURL, expectedGenBaseURL)
	}
	expectedGenPath = "/api/v1/generate"
	if confEnv.GeneratorServicePath != expectedGenPath {
		t.Errorf("Generator Service  Path is incorrect, got %s, want: %s", confEnv.GeneratorServicePath, expectedGenPath)
	}
	expectedDBBaseURL = "http://someurl:3001"
	if confEnv.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got %s, want: %s", confEnv.DBBaseURL, expectedDBBaseURL)
	}
	expectedDBPath = "/api/v1/secrets"
	if confEnv.DBServicePath != expectedDBPath {
		t.Errorf("DB Service Path is incorrect, got %s, want: %s", confEnv.DBServicePath, expectedDBPath)
	}
	expectedDBAPIKey = "AAAA"
	if confEnv.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got %s, want: %s", confEnv.DBAPIKey, expectedDBAPIKey)
	}
}

func TestConfigureFromFile(t *testing.T) {
	configPath := "../test/config.yaml"
	conf, err := configureFromFile(configPath)
	if err != nil {
		t.Errorf("%v", err)
	}

	expectedPort := "3003"
	if conf.Port != expectedPort {
		t.Errorf("Port value is incorrect, got %s, want: %s", conf.Port, expectedPort)
	}
	expectedGenBaseURL := "http://localhost:3003"
	if conf.GeneratorBaseURL != expectedGenBaseURL {
		t.Errorf("Generator Base URL is incorrect, got %s, want: %s", conf.GeneratorBaseURL, expectedGenBaseURL)
	}
	expectedGenPath := "/api/v1/generate"
	if conf.GeneratorServicePath != expectedGenPath {
		t.Errorf("Generator Service Path is incorrect, got %s, want: %s", conf.GeneratorServicePath, expectedGenPath)
	}
	expectedDBBaseURL := "http://localhost:3003"
	if conf.DBBaseURL != expectedDBBaseURL {
		t.Errorf("DB Base URL is incorrect, got %s, want: %s", conf.DBBaseURL, expectedDBBaseURL)
	}
	expectedDBPath := "/api/v1/secrets"
	if conf.DBServicePath != expectedDBPath {
		t.Errorf("DB Service Path is incorrect, got %s, want: %s", conf.DBServicePath, expectedDBPath)
	}
	expectedDBAPIKey := "aabbcc"
	if conf.DBAPIKey != expectedDBAPIKey {
		t.Errorf("DB API Key is incorrect, got %s, want: %s", conf.DBAPIKey, expectedDBAPIKey)
	}
}

func TestConfigure(t *testing.T) {
	// Test Configure from environment.
	_, err := Configure("")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	// Test Configure from file.
	_, err = Configure("../test/config.yaml")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	// Handle whene no configuration exists.
	_, err = Configure("../test/nofile.yml")
	if err == nil {
		t.Errorf("Incorrect, should have returned an error")
	}
}
