package config

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	configPath = "../test/config.yaml"
)

func TestNewConfiguration(t *testing.T) {
	want := &Configuration{
		Server: Server{
			Host: defaultListenHost,
			Port: defaultListenPort,
		},
		Database: Database{
			Driver:   defaultDBDriver,
			Address:  defaultDBAddress,
			Database: defaultDB,
			SSL:      true,
		},
	}

	got := newConfiguration()

	if !cmp.Equal(want, got) {
		t.Log(cmp.Diff(want, got))
		t.Errorf("results differ in test\n")
	}
}

func TestConfigure(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			envVars map[string]string
			flags   flags
		}
		want    *Configuration
		wantErr error
	}{
		{
			name: "from file",
			input: struct {
				envVars map[string]string
				flags   flags
			}{
				envVars: map[string]string{},
				flags: flags{
					ConfigPath: configPath,
				},
			},
			want: &Configuration{
				Server: Server{
					Host: "127.0.0.1",
					Port: "3003",
					Security: Security{
						Encryption: Encryption{
							Key: "secretstring",
						},
						TLS: TLS{
							Certificate: "path/to/cert",
							Key:         "path/to/key",
						},
						CORS: CORS{
							Origin: "http://localhost",
						},
					},
				},
				Database: Database{
					Driver:        "mongo",
					Address:       "localhost:27017",
					Database:      "burnit_db",
					Username:      "dbuser",
					Password:      "dbpassword",
					SSL:           false,
					DirectConnect: true,
					URI:           "mongodb://localhost:27017",
				},
			},
		},
		{
			name: "from file and environment",
			input: struct {
				envVars map[string]string
				flags   flags
			}{
				envVars: map[string]string{
					"BURNIT_LISTEN_HOST":     "127.0.0.2",
					"BURNIT_LISTEN_PORT":     "3004",
					"BURNIT_ENCRYPTION_KEY":  "secretstring1",
					"BURNIT_TLS_CERTIFICATE": "path/to/cert1",
					"BURNIT_TLS_KEY":         "path/to/key1",
					"BURNIT_CORS_ORIGIN":     "http://localhost1",
					"DB_DRIVER":              "redis",
					"DB_HOST":                "localhost:27018",
					"DB":                     "burnit_db1",
					"DB_CONNECTION_URI":      "mongodb://localhost:27018",
					"DB_USER":                "dbuser1",
					"DB_PASSWORD":            "dbpassword1",
					"DB_SSL":                 "true",
					"DB_DIRECT_CONNECT":      "false",
				},
				flags: flags{},
			},
			want: &Configuration{
				Server: Server{
					Host: "127.0.0.2",
					Port: "3004",
					Security: Security{
						Encryption: Encryption{
							Key: "secretstring1",
						},
						TLS: TLS{
							Certificate: "path/to/cert1",
							Key:         "path/to/key1",
						},
						CORS: CORS{
							Origin: "http://localhost1",
						},
					},
				},
				Database: Database{
					Driver:        "redis",
					Address:       "mongodb://localhost:27018",
					Database:      "burnit_db1",
					Username:      "dbuser1",
					Password:      "dbpassword1",
					SSL:           true,
					DirectConnect: false,
					URI:           "mongodb://localhost:27018",
				},
			},
		},
		{
			name: "from file, environment and flags",
			input: struct {
				envVars map[string]string
				flags   flags
			}{
				envVars: map[string]string{},
				flags: flags{
					ConfigPath:      configPath,
					Host:            "127.0.0.1",
					Port:            "3003",
					TLSCertificate:  "path/to/cert",
					TLSKey:          "path/to/key",
					CORSOrigin:      "http://localhost",
					EncryptionKey:   "secretstring",
					Driver:          "mongo",
					DBAddress:       "localhost:27017",
					DBURI:           "mongodb://localhost:27017",
					DB:              "burnit_db",
					DBUser:          "dbuser",
					DBPassword:      "dbpassword",
					DisableDBSSL:    true,
					DBDirectConnect: true,
				},
			},
			want: &Configuration{
				Server: Server{
					Host: "127.0.0.1",
					Port: "3003",
					Security: Security{
						Encryption: Encryption{
							Key: "secretstring",
						},
						TLS: TLS{
							Certificate: "path/to/cert",
							Key:         "path/to/key",
						},
						CORS: CORS{
							Origin: "http://localhost",
						},
					},
				},
				Database: Database{
					Driver:        "mongo",
					Address:       "localhost:27017",
					Database:      "burnit_db",
					Username:      "dbuser",
					Password:      "dbpassword",
					SSL:           false,
					DirectConnect: true,
					URI:           "mongodb://localhost:27017",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnvVars(test.input.envVars)

			got, gotErr := configure(test.input.flags)

			handleTestResults(t, test.name, got, gotErr, test.want, test.wantErr)

			unsetEnvVars(test.input.envVars)
		})
	}
}

func TestFromFile(t *testing.T) {
	var tests = []struct {
		name    string
		want    *Configuration
		wantErr error
	}{
		{
			name: "set from file",
			want: &Configuration{
				Server: Server{
					Host: "127.0.0.1",
					Port: "3003",
					Security: Security{
						Encryption: Encryption{
							Key: "secretstring",
						},
						TLS: TLS{
							Certificate: "path/to/cert",
							Key:         "path/to/key",
						},
						CORS: CORS{
							Origin: "http://localhost",
						},
					},
				},
				Database: Database{
					Driver:        "mongo",
					Address:       "localhost:27017",
					Database:      "burnit_db",
					Username:      "dbuser",
					Password:      "dbpassword",
					SSL:           false,
					DirectConnect: true,
					URI:           "mongodb://localhost:27017",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := newConfiguration()
			gotErr := fromFile(cfg, configPath)

			handleTestResults(t, test.name, cfg, gotErr, test.want, test.wantErr)
		})
	}
}

func TestFromEnv(t *testing.T) {
	var tests = []struct {
		name  string
		input map[string]string
		want  *Configuration
	}{
		{
			name: "set from environment",
			input: map[string]string{
				"BURNIT_LISTEN_HOST":     "127.0.0.1",
				"BURNIT_LISTEN_PORT":     "3003",
				"BURNIT_ENCRYPTION_KEY":  "secretstring",
				"BURNIT_TLS_CERTIFICATE": "path/to/cert",
				"BURNIT_TLS_KEY":         "path/to/key",
				"BURNIT_CORS_ORIGIN":     "http://localhost",
				"DB_DRIVER":              "mongo",
				"DB_HOST":                "localhost:27017",
				"DB":                     "burnit_db",
				"DB_CONNECTION_URI":      "mongodb://localhost:27017",
				"DB_USER":                "dbuser",
				"DB_PASSWORD":            "dbpassword",
				"DB_SSL":                 "false",
				"DB_DIRECT_CONNECT":      "true",
			},
			want: &Configuration{
				Server: Server{
					Host: "127.0.0.1",
					Port: "3003",
					Security: Security{
						Encryption: Encryption{
							Key: "secretstring",
						},
						TLS: TLS{
							Certificate: "path/to/cert",
							Key:         "path/to/key",
						},
						CORS: CORS{
							Origin: "http://localhost",
						},
					},
				},
				Database: Database{
					Driver:        "mongo",
					Address:       "localhost:27017",
					Database:      "burnit_db",
					Username:      "dbuser",
					Password:      "dbpassword",
					SSL:           false,
					DirectConnect: true,
					URI:           "mongodb://localhost:27017",
				},
			},
		},
		{
			name:  "set from environment - keep values if empty",
			input: map[string]string{},
			want: &Configuration{
				Server: Server{
					Host: defaultListenHost,
					Port: defaultListenPort,
				},
				Database: Database{
					Driver:   defaultDBDriver,
					Address:  defaultDBAddress,
					Database: defaultDB,
					SSL:      defaultDBSSL,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnvVars(test.input)

			cfg := newConfiguration()
			fromEnv(cfg)

			handleTestResults(t, test.name, cfg, nil, test.want, nil)

			unsetEnvVars(test.input)
		})
	}
}

func TestFromFlags(t *testing.T) {
	var tests = []struct {
		name  string
		input flags
		want  *Configuration
	}{
		{
			name: "set from flags",
			input: flags{
				ConfigPath:      "../test/config.yaml",
				Host:            "127.0.0.1",
				Port:            "3003",
				TLSCertificate:  "path/to/cert",
				TLSKey:          "path/to/key",
				CORSOrigin:      "http://localhost",
				EncryptionKey:   "secretstring",
				Driver:          "mongo",
				DBAddress:       "localhost:27017",
				DBURI:           "mongodb://localhost:27017",
				DB:              "burnit_db",
				DBUser:          "dbuser",
				DBPassword:      "dbpassword",
				DisableDBSSL:    true,
				DBDirectConnect: true,
			},
			want: &Configuration{
				Server: Server{
					Host: "127.0.0.1",
					Port: "3003",
					Security: Security{
						Encryption: Encryption{
							Key: "secretstring",
						},
						TLS: TLS{
							Certificate: "path/to/cert",
							Key:         "path/to/key",
						},
						CORS: CORS{
							Origin: "http://localhost",
						},
					},
				},
				Database: Database{
					Driver:        "mongo",
					Address:       "localhost:27017",
					Database:      "burnit_db",
					Username:      "dbuser",
					Password:      "dbpassword",
					SSL:           false,
					DirectConnect: true,
					URI:           "mongodb://localhost:27017",
				},
			},
		},
		{
			name:  "set from flags - keep values if empty",
			input: flags{},
			want: &Configuration{
				Server: Server{
					Host: defaultListenHost,
					Port: defaultListenPort,
				},
				Database: Database{
					Driver:   defaultDBDriver,
					Address:  defaultDBAddress,
					Database: defaultDB,
					SSL:      defaultDBSSL,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := newConfiguration()
			fromFlags(cfg, test.input)

			handleTestResults(t, test.name, cfg, nil, test.want, nil)
		})
	}
}

func TestAddressFromMongoURI(t *testing.T) {
	var tests = []struct {
		uri string
	}{
		{uri: "mongodb://localhost:27017"},
		{uri: "mongodb://localhost:27017/?ssl=true"},
		{uri: "mongodb://user:pass@localhost:27017"},
		{uri: "mongodb://user:pass@localhost:27017/?ssl=true"},
	}

	expected := "localhost:27017"
	for _, test := range tests {
		addr := AddressFromMongoURI(test.uri)
		if addr != expected {
			t.Errorf("incorrect value, got: %s, want: %s", addr, expected)
		}
	}
}

func TestAddressFromRedisURI(t *testing.T) {
	var tests = []struct {
		uri string
	}{
		{uri: "localhost:6379"},
		{uri: "redis://localhost:6379"},
		{uri: "rediss://localhost:6379"},
		{uri: "localhost:6379,password=1234,ssl=true"},
		{uri: "redis://localhost:6379,password=1234,ssl=true"},
		{uri: "rediss://localhost:6379,password=1234,ssl=true"},
	}

	expected := "localhost:6379"
	for _, test := range tests {
		addr := AddressFromRedisURI(test.uri)
		if addr != expected {
			t.Errorf("incorrect value, got: %s, want: %s", addr, expected)
		}
	}
}

func handleTestResults(t *testing.T, name string, got *Configuration, gotErr error, want *Configuration, wantErr error) {
	if wantErr == nil && gotErr != nil {
		t.Errorf("error in test: %s, should not return error, error returned: %v\n", name, gotErr)
	}
	if wantErr != nil && gotErr == nil {
		t.Errorf("error in test %s, should return error\n", name)
	}

	if !cmp.Equal(want, got) {
		t.Log(cmp.Diff(want, got))
		t.Errorf("results differ in test: %s\n", name)
	}

	if wantErr != nil && gotErr != nil {
		if !cmp.Equal(wantErr.Error(), gotErr.Error(), cmpopts.EquateErrors()) {
			t.Log(cmp.Diff(wantErr.Error(), gotErr.Error(), cmpopts.EquateErrors()))
			t.Errorf("results differ in test: %s\n", name)
		}
	}
}

func setEnvVars(envVars map[string]string) {
	for k, v := range envVars {
		os.Setenv(k, v)
	}
}

func unsetEnvVars(envVars map[string]string) {
	for k := range envVars {
		os.Unsetenv(k)
	}
}
