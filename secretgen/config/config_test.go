package config

import (
	"os"
	"testing"
)

func TestConfigureFromEnv(t *testing.T) {
	confDefault := configureFromEnv()
	expected := "3002"
	if confDefault.Port != expected {
		t.Errorf("default port value is incorrect, got %s, want: %s", confDefault.Port, expected)
	}

	expected = "5000"
	os.Setenv("SECRET_GEN_PORT", expected)
	confEnv := configureFromEnv()
	if confEnv.Port != expected {
		t.Errorf("port value is incorrect, got %s, want: %s", confEnv.Port, expected)
	}
}

func TestConfigureFromFile(t *testing.T) {
	configPath := "../test/config.yaml"
	conf, err := configureFromFile(configPath)
	if err != nil {
		t.Errorf("%v", err)
	}

	expected := "3003"
	if conf.Port != expected {
		t.Errorf("port value is incorrect, got %s, want: %s", conf.Port, expected)
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
