package config

import (
	"os"
	"testing"
)

func TestConfigureFromEnv(t *testing.T) {
	confDefault := configureFromEnv()
	if confDefault.Port != "3000" {
		t.Errorf("Default port value is incorrect, got %s, want: 3000", confDefault.Port)
	}
	if confDefault.GeneratorBaseURL != "http://localhost:3002" {
		t.Errorf("Generator Base URL is incorrect, got %s, want: http://localhost:3002", confDefault.GeneratorBaseURL)
	}
	if confDefault.GeneratorServicePath != "/api/v0/generate" {
		t.Errorf("Generator Service  Path is incorrect, got %s, want: /api/v0/generate", confDefault.GeneratorServicePath)
	}
	if confDefault.DBBaseURL != "http://localhost:3001" {
		t.Errorf("DB Base URL is incorrect, got %s, want http://localhost:3001", confDefault.DBBaseURL)
	}
	if confDefault.DBServicePath != "/api/v0/secrets" {
		t.Errorf("DB Service Path is incorrect, got %s, want /api/v0/secrets", confDefault.DBServicePath)
	}
	if confDefault.DBAPIKey != "" {
		t.Errorf("DB API Key is incorrect, got %s, want empty string", confDefault.DBAPIKey)
	}

	os.Setenv("SECRET_GW_PORT", "5000")
	os.Setenv("SECRET_GEN_BASE_URL", "http://someurl:3000")
	os.Setenv("SECRET_GEN_PATH", "/api/v1/generate")
	os.Setenv("SECRET_DB_BASE_URL", "http://someurl:3001")
	os.Setenv("SECRET_DB_PATH", "/api/v1/secrets")
	os.Setenv("SECRET_DB_API_KEY", "AAAA")
	confEnv := configureFromEnv()

	if confEnv.Port != "5000" {
		t.Errorf("Port value is incorrect, got %s, want: 5000", confEnv.Port)
	}
	if confEnv.GeneratorBaseURL != "http://someurl:3000" {
		t.Errorf("Generator Base URL is incorrect, got %s, want: http://someurl:3000", confEnv.GeneratorBaseURL)
	}
	if confEnv.GeneratorServicePath != "/api/v1/generate" {
		t.Errorf("Generator Service  Path is incorrect, got %s, want: /api/v1/generate", confEnv.GeneratorServicePath)
	}
	if confEnv.DBBaseURL != "http://someurl:3001" {
		t.Errorf("DB Base URL is incorrect, got %s, want: http://someurl:3001", confEnv.DBBaseURL)
	}
	if confEnv.DBServicePath != "/api/v1/secrets" {
		t.Errorf("DB Service Path is incorrect, got %s, want: /api/v1/secrets", confEnv.DBServicePath)
	}
	if confEnv.DBAPIKey != "AAAA" {
		t.Errorf("DB API Key is incorrect, got %s, want: AAAA", confEnv.DBAPIKey)
	}
}

func TestConfigureFromFile(t *testing.T) {
	configPath := "../test/config.yaml"
	conf, err := configureFromFile(configPath)
	if err != nil {
		t.Errorf("%v", err)
	}

	if conf.Port != "3003" {
		t.Errorf("Port value is incorrect, got %s, want: 3003", conf.Port)
	}
	if conf.GeneratorBaseURL != "http://localhost:3003" {
		t.Errorf("Generator Base URL is incorrect, got %s, want: http://localhost:3003", conf.GeneratorBaseURL)
	}
	if conf.GeneratorServicePath != "/api/v1/generate" {
		t.Errorf("Generator Service Path is incorrect, got %s, want: /api/v1/generate", conf.GeneratorServicePath)
	}
	if conf.DBBaseURL != "http://localhost:3003" {
		t.Errorf("DB Base URL is incorrect, got %s, want: http://localhost:3003", conf.DBBaseURL)
	}
	if conf.DBServicePath != "/api/v1/secrets" {
		t.Errorf("DB Service Path is incorrect, got %s, want: /api/v1/secrets", conf.DBServicePath)
	}
	if conf.DBAPIKey != "aabbcc" {
		t.Errorf("DB API Key is incorrect, got %s, want: aabbcc", conf.DBAPIKey)
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
