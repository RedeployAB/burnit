package configs

import (
	"os"
	"testing"
)

func TestConfigure(t *testing.T) {
	confDefault := Configure()
	if confDefault.Server.Port != "3001" {
		t.Errorf("Default port value is incorrect, got %s, want: 3001", confDefault.Server.Port)
	}

	os.Setenv("SECRET_DB_PORT", "6000")
	confEnv := Configure()
	if confEnv.Server.Port != "6000" {
		t.Errorf("Port value is incorrect, got %s, want: 6000", confEnv.Server.Port)
	}
}
