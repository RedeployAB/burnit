package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConfigureDefault(t *testing.T) {
	var flags Flags
	config, err := Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	expected := &Configuration{
		Server: Server{
			Port:                 "3000",
			GeneratorAddress:     "http://localhost:3002",
			GeneratorServicePath: "/secret",
			DBAddress:            "http://localhost:3001",
			DBServicePath:        "/secrets",
		},
	}

	if !cmp.Equal(expected, config) {
		t.Log(cmp.Diff(expected, config))
		t.Errorf("incorrect, configurations differ")
	}
}

func TestConfigureFromFile(t *testing.T) {
	configPath := "../test/config.yaml"
	expected := &Configuration{
		Server: Server{
			Port:                 "3003",
			GeneratorAddress:     "http://localhost:3003",
			GeneratorServicePath: "/v1/secret",
			DBAddress:            "http://localhost:3003",
			DBServicePath:        "/v1/secrets",
			DBAPIKey:             "aabbcc",
			TLS: TLS{
				Certificate: "path/to/cert",
				Key:         "path/to/key",
			},
		},
	}

	config := &Configuration{}
	if err := configureFromFile(config, configPath); err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if !cmp.Equal(expected, config) {
		t.Log(cmp.Diff(expected, config))
		t.Errorf("incorrect, configurations differ")
	}
}

func TestConfigureFromEnv(t *testing.T) {
	expected := &Configuration{
		Server: Server{
			Port:                 "3003",
			GeneratorAddress:     "http://someurl:3002",
			GeneratorServicePath: "/v1/secret",
			DBAddress:            "http://someurl:3001",
			DBServicePath:        "/v1/secrets",
			DBAPIKey:             "AAAA",
			TLS: TLS{
				Certificate: "path/to/cert",
				Key:         "path/to/key",
			},
		},
	}

	config := &Configuration{}
	os.Setenv("BURNITGW_LISTEN_PORT", expected.Port)
	os.Setenv("BURNITGEN_ADDRESS", expected.GeneratorAddress)
	os.Setenv("BURNITGEN_PATH", expected.GeneratorServicePath)
	os.Setenv("BURNITDB_ADDRESS", expected.DBAddress)
	os.Setenv("BURNITDB_PATH", expected.DBServicePath)
	os.Setenv("BURNITDB_API_KEY", expected.DBAPIKey)
	os.Setenv("BURNITGW_TLS_CERTIFICATE", expected.TLS.Certificate)
	os.Setenv("BURNITGW_TLS_KEY", expected.TLS.Key)
	configureFromEnv(config)

	if !cmp.Equal(expected, config) {
		t.Log(cmp.Diff(expected, config))
		t.Errorf("incorrect, configurations differ")
	}

	os.Setenv("BURNITGW_LISTEN_PORT", "")
	os.Setenv("BURNITGEN_ADDRESS", "")
	os.Setenv("BURNITGEN_PATH", "")
	os.Setenv("BURNITDB_ADDRESS", "")
	os.Setenv("BURNITDB_PATH", "")
	os.Setenv("BURNITDB_API_KEY", "")
	os.Setenv("BURNITGW_TLS_CERTIFICATE", "")
	os.Setenv("BURNITGW_TLS_KEY", "")
}

func TestConfigureFromFlags(t *testing.T) {
	expected := &Configuration{
		Server: Server{
			Port:                 "4000",
			GeneratorAddress:     "http://someurl:4002",
			GeneratorServicePath: "/v1/secret",
			DBAddress:            "http://someurl:4003",
			DBServicePath:        "/v1/secrets",
			DBAPIKey:             "ccaabb",
			TLS: TLS{
				Certificate: "path/to/cert",
				Key:         "path/to/key",
			},
		},
	}

	flags := Flags{
		Port:                 expected.Port,
		GeneratorAddress:     expected.GeneratorAddress,
		GeneratorServicePath: expected.GeneratorServicePath,
		DBAddress:            expected.DBAddress,
		DBServicePath:        expected.DBServicePath,
		DBAPIKey:             expected.DBAPIKey,
		TLSCert:              expected.TLS.Certificate,
		TLSKey:               expected.TLS.Key,
	}

	config := &Configuration{}
	configureFromFlags(config, flags)

	if !cmp.Equal(expected, config) {
		t.Log(cmp.Diff(expected, config))
		t.Errorf("incorrect, configurations differ")
	}
}

func TestConfigure(t *testing.T) {
	configPath := "../test/config.yaml"
	// Test default configuration.
	expected := &Configuration{
		Server: Server{
			Port:                 "3000",
			GeneratorAddress:     "http://localhost:3002",
			GeneratorServicePath: "/secret",
			DBAddress:            "http://localhost:3001",
			DBServicePath:        "/secrets",
			TLS:                  TLS{},
		},
	}

	var flags Flags
	config, err := Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if !cmp.Equal(expected, config) {
		t.Log(cmp.Diff(expected, config))
		t.Errorf("incorrect, configurations differ")
	}

	// Test with file. Should override default.
	expected = &Configuration{
		Server: Server{
			Port:                 "3003",
			GeneratorAddress:     "http://localhost:3003",
			GeneratorServicePath: "/v1/secret",
			DBAddress:            "http://localhost:3003",
			DBServicePath:        "/v1/secrets",
			DBAPIKey:             "aabbcc",
			TLS: TLS{
				Certificate: "path/to/cert",
				Key:         "path/to/key",
			},
		},
	}

	flags.ConfigPath = configPath
	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if !cmp.Equal(expected, config) {
		t.Log(cmp.Diff(expected, config))
		t.Errorf("incorrect, configurations differ")
	}

	// Test with environment variables. Should override file.
	expected = &Configuration{
		Server: Server{
			Port:                 "3003",
			GeneratorAddress:     "http://localhost:3002",
			GeneratorServicePath: "/v1/secret",
			DBAddress:            "http://localhost:3001",
			DBServicePath:        "/v1/secrets",
			DBAPIKey:             "AAAA",
			TLS: TLS{
				Certificate: "another/path/to/cert",
				Key:         "another/path/to/key",
			},
		},
	}

	os.Setenv("BURNITGW_LISTEN_PORT", expected.Port)
	os.Setenv("BURNITGEN_ADDRESS", expected.GeneratorAddress)
	os.Setenv("BURNITGEN_PATH", expected.GeneratorServicePath)
	os.Setenv("BURNITDB_ADDRESS", expected.DBAddress)
	os.Setenv("BURNITDB_PATH", expected.DBServicePath)
	os.Setenv("BURNITDB_API_KEY", expected.DBAPIKey)
	os.Setenv("BURNITGW_TLS_CERTIFICATE", expected.TLS.Certificate)
	os.Setenv("BURNITGW_TLS_KEY", expected.TLS.Key)

	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if !cmp.Equal(expected, config) {
		t.Log(cmp.Diff(expected, config))
		t.Errorf("incorrect, configurations differ")
	}

	// Test with flags. Should override file and envrionment variables.
	expected = &Configuration{
		Server: Server{
			Port:                 "4000",
			GeneratorAddress:     "http://localhost:4002",
			GeneratorServicePath: "/v1/secret",
			DBAddress:            "http://localhost:4003",
			DBServicePath:        "/v1/secrets",
			DBAPIKey:             "ccaabb",
			TLS: TLS{
				Certificate: "third/path/to/cert",
				Key:         "third/path/to/key",
			},
		},
	}

	flags = Flags{
		ConfigPath:           configPath,
		Port:                 expected.Port,
		GeneratorAddress:     expected.GeneratorAddress,
		GeneratorServicePath: expected.GeneratorServicePath,
		DBAddress:            expected.DBAddress,
		DBServicePath:        expected.DBServicePath,
		DBAPIKey:             expected.DBAPIKey,
		TLSCert:              expected.TLS.Certificate,
		TLSKey:               expected.TLS.Key,
	}

	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if !cmp.Equal(expected, config) {
		t.Log(cmp.Diff(expected, config))
		t.Errorf("incorrect, configurations differ")
	}

	os.Setenv("BURNITGW_LISTEN_PORT", "")
	os.Setenv("BURNITGEN_ADDRESS", "")
	os.Setenv("BURNITGEN_PATH", "")
	os.Setenv("BURNITDB_ADDRESS", "")
	os.Setenv("BURNITDB_PATH", "")
	os.Setenv("BURNITDB_API_KEY", "")
	os.Setenv("BURNITGW_TLS_CERTIFICATE", "")
	os.Setenv("BURNITGW_TLS_KEY", "")

	// Test with flags. Should override file and envrionment variables.
	// Should not have prefix http:// or https:// in address arguments.
	expected = &Configuration{
		Server: Server{
			Port:                 "4000",
			GeneratorAddress:     "http://someurl:4002",
			GeneratorServicePath: "/v1/secret",
			DBAddress:            "http://someurl:4003",
			DBServicePath:        "/v1/secrets",
			DBAPIKey:             "ccaabb",
			TLS: TLS{
				Certificate: "path/to/cert",
				Key:         "path/to/key",
			},
		},
	}

	flags = Flags{
		ConfigPath:           configPath,
		Port:                 expected.Port,
		GeneratorAddress:     "someurl:4002",
		GeneratorServicePath: expected.GeneratorServicePath,
		DBAddress:            "someurl:4003",
		DBServicePath:        expected.DBServicePath,
		DBAPIKey:             expected.DBAPIKey,
		TLSCert:              "",
		TLSKey:               "",
	}

	config, err = Configure(flags)
	if err != nil {
		t.Fatalf("error in test: %v", err)
	}

	if !cmp.Equal(expected, config) {
		t.Log(cmp.Diff(expected, config))
		t.Errorf("incorrect, configurations differ")
	}

	os.Setenv("BURNITGW_LISTEN_PORT", "")
	os.Setenv("BURNITGEN_ADDRESS", "")
	os.Setenv("BURNITGEN_PATH", "")
	os.Setenv("BURNITDB_ADDRESS", "")
	os.Setenv("BURNITDB_PATH", "")
	os.Setenv("BURNITDB_API_KEY", "")
	os.Setenv("BURNITGW_TLS_CERTIFICATE", "")
	os.Setenv("BURNITGW_TLS_KEY", "")
}
