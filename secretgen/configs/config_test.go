package configs

import (
	"os"
	"testing"
)

func TestConfigure(t *testing.T) {
	confDefault := Configure()
	if confDefault.Port != "3002" {
		t.Errorf("Default port value is incorrect, got %s, want: 3002", confDefault.Port)
	}

	os.Setenv("SECRET_GENERATOR_PORT", "5000")
	confEnv := Configure()
	if confEnv.Port != "5000" {
		t.Errorf("Port value is incorrect, got %s, want: 5000", confEnv.Port)
	}
}
