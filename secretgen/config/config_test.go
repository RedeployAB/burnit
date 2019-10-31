package config

import (
	"os"
	"testing"
)

func TestConfigureFromEnv(t *testing.T) {
	confDefault := configureFromEnv()
	if confDefault.Port != "3002" {
		t.Errorf("Default port value is incorrect, got %s, want: 3002", confDefault.Port)
	}

	os.Setenv("SECRET_GEN_PORT", "5000")
	confEnv := configureFromEnv()
	if confEnv.Port != "5000" {
		t.Errorf("Port value is incorrect, got %s, want: 5000", confEnv.Port)
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
}

func TestConfigure(t *testing.T) {

	// First test configuration from file.
	configPath := "../test/config.yaml"
	conf, err := Configure(configPath)
	if err != nil {
		t.Errorf("%v", err)
	}

	if conf.Port != "3003" {
		t.Errorf("Port value is incorrect, got %s, want: 3003", conf.Port)
	}

	// Test from set environment variable.
	os.Setenv("SECRET_GEN_PORT", "5000")
	configPath = ""
	conf, err = Configure(configPath)
	if err != nil {
		t.Errorf("%v", err)
	}

	if conf.Port != "5000" {
		t.Errorf("Port value is incorrect, got %s, want: 5000", conf.Port)
	}
}
