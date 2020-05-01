package config

import (
	"os"
	"testing"
)

func TestConfigureDefault(t *testing.T) {
	expectedPort := "3002"
	var flags Flags
	config, err := Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
}

func TestConfigureFromEnv(t *testing.T) {
	expectedPort := "5000"

	config := &Configuration{}
	os.Setenv("BURNITGEN_LISTEN_PORT", expectedPort)
	configureFromEnv(config)

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	os.Setenv("BURNITGEN_LISTEN_PORT", "")
}

func TestConfigureFromFile(t *testing.T) {
	expectedPort := "3003"
	configPath := "../test/config.yaml"

	config := &Configuration{}
	if err := configureFromFile(config, configPath); err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}

	if err := configureFromFile(config, "nonexisting.yml"); err == nil {
		t.Errorf("should return error if file does not exist")
	}
}

func TestConfigureFromFlags(t *testing.T) {
	expectedPort := "3004"
	flags := Flags{
		Port: expectedPort,
	}

	config := &Configuration{}
	configureFromFlags(config, flags)
	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
}

func TestConfigure(t *testing.T) {
	configPath := "../test/config.yaml"
	// Test default configuration.
	expectedPort := "3002"
	var flags Flags
	config, err := Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}

	expectedPort = "3003"
	// Test with file. Should override default.
	flags.ConfigPath = configPath
	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}

	// Test with environment variables. Should override file.
	expectedPort = "3004"
	os.Setenv("BURNITGEN_LISTEN_PORT", "3004")
	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
	// Test with flags. Should override file and envrionment variables.
	expectedPort = "3005"
	flags.Port = expectedPort
	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if config.Port != expectedPort {
		t.Errorf("Port was incorrect, got: %s, want: %s", config.Port, expectedPort)
	}
}
