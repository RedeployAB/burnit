package config

import (
	"os"
	"testing"
)

func TestConfigure(t *testing.T) {
	confDefault := Configure()
	if confDefault.Port != "3001" {
		t.Errorf("Default port value is incorrect, got %s, want: 3001", confDefault.Port)
	}

	os.Setenv("SECRET_DB_PORT", "6000")
	os.Setenv("SECRET_GENERATOR_BASE_URL", "http://generator:7000")
	os.Setenv("SECRET_GENERATOR_REQUEST_PATH", "/api/v2/secret")
	confEnv := Configure()
	if confEnv.Port != "6000" {
		t.Errorf("Port value is incorrect, got %s, want: 6000", confEnv.Port)
	}
}
