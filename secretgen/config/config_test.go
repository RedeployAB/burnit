package config

import (
	"os"
	"testing"
)

func TestConfigureFromEnv(t *testing.T) {
	confDefault := configureFromEnv()
	if confDefault.Port != "3002" {
		t.Errorf("default port value is incorrect, got %s, want: 3002", confDefault.Port)
	}

	os.Setenv("SECRET_GEN_PORT", "5000")
	confEnv := configureFromEnv()
	if confEnv.Port != "5000" {
		t.Errorf("port value is incorrect, got %s, want: 5000", confEnv.Port)
	}
}

func TestConfigureFromFile(t *testing.T) {
	configPath := "../test/config.yaml"
	conf, err := configureFromFile(configPath)
	if err != nil {
		t.Errorf("%v", err)
	}

	if conf.Port != "3003" {
		t.Errorf("port value is incorrect, got %s, want: 3003", conf.Port)
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
		t.Errorf("error: %v", err)
	}
	// Handle whene no configuration exists.
	_, err = Configure("../test/nofile.yml")
	if err == nil {
		t.Errorf("error in test, should return an error")
	}
}
