package db

import (
	"testing"

	"github.com/RedeployAB/burnit/burnitdb/config"
)

func TestToURI(t *testing.T) {
	cfg1 := config.Configuration{
		Server: config.Server{},
		Database: config.Database{
			Address: "localhost",
		},
	}
	cfg2 := config.Configuration{
		Server: config.Server{},
		Database: config.Database{
			Address: "localhost:27017",
		},
	}
	cfg3 := config.Configuration{
		Server: config.Server{},
		Database: config.Database{
			Address:  "localhost:27017",
			Username: "user",
		},
	}
	cfg4 := config.Configuration{
		Server: config.Server{},
		Database: config.Database{
			Address:  "localhost:27017",
			Username: "user",
			Password: "1234",
		},
	}
	cfg5 := config.Configuration{
		Server: config.Server{},
		Database: config.Database{
			Address:  "localhost:27017",
			Username: "user",
			Password: "1234",
			Database: "db",
		},
	}
	cfg6 := config.Configuration{
		Server: config.Server{},
		Database: config.Database{
			Address:  "localhost:27017",
			Username: "user",
			Password: "1234",
			Database: "db",
			SSL:      true,
		},
	}
	var tests = []struct {
		config      config.Configuration
		expectedURI string
	}{
		{config: cfg1, expectedURI: "mongodb://localhost"},
		{config: cfg2, expectedURI: "mongodb://localhost:27017"},
		{config: cfg3, expectedURI: "mongodb://user@localhost:27017"},
		{config: cfg4, expectedURI: "mongodb://user:1234@localhost:27017"},
		{config: cfg5, expectedURI: "mongodb://user:1234@localhost:27017/db"},
		{config: cfg6, expectedURI: "mongodb://user:1234@localhost:27017/db?ssl=true"},
	}

	for i, test := range tests {
		got := toURI(test.config.Database)
		if got != test.expectedURI {
			t.Errorf("[%d] incorrect value, got: %s, want: %s", i+1, got, test.expectedURI)
		}
	}
}
