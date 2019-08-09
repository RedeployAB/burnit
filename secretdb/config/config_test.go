package config

import (
	"os"
	"testing"
)

func TestConfigure(t *testing.T) {
	confDefault := Configure()
	if confDefault.Port != "4000" {
		t.Errorf("Default port value is incorrect, got %s, want: 4000", confDefault.Port)
	}
	if confDefault.GeneratorBaseURL != "http://localhost:3000" {
		t.Errorf("Default Generator base url is incorrect, got: '%s' want: 'http://localhost:3000'", confDefault.GeneratorBaseURL)
	}
	if confDefault.GeneratorPath != "/api/v1/secret" {
		t.Errorf("Default Generator URL path is incorrect, got: '%s', want: '/api/v1/secret'", confDefault.GeneratorPath)
	}

	os.Setenv("SECRET_DB_PORT", "6000")
	os.Setenv("SECRET_GENERATOR_BASE_URL", "http://generator:7000")
	os.Setenv("SECRET_GENERATOR_REQUEST_PATH", "/api/v2/secret")
	confEnv := Configure()
	if confEnv.Port != "6000" {
		t.Errorf("Port value is incorrect, got %s, want: 6000", confEnv.Port)
	}
	if confEnv.GeneratorBaseURL != "http://generator:7000" {
		t.Errorf("Generator base url is incorrect, got: '%s', want: 'http://generator:7000'", confEnv.GeneratorBaseURL)
	}
	if confEnv.GeneratorPath != "/api/v2/secret" {
		t.Errorf("Default Generator URL path is incorrect, got: '%s', want: '/api/v2/secret'", confEnv.GeneratorPath)
	}

}
