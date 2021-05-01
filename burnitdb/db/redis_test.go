package db

import (
	"testing"

	"github.com/RedeployAB/burnit/burnitdb/config"
)

var opts = config.Configuration{
	Server: config.Server{},
	Database: config.Database{
		Driver:   "redis",
		Address:  "",
		URI:      "",
		Database: "",
		Username: "",
		Password: "",
		SSL:      false,
	},
}

func TestFromURI(t *testing.T) {
	var tests = []struct {
		uri              string
		expectedAddress  string
		expectedPassword string
		expectedSSL      bool
	}{
		{uri: "localhost:6379", expectedAddress: "localhost:6379", expectedPassword: "", expectedSSL: false},
		{uri: "redis://localhost:6379", expectedAddress: "localhost:6379", expectedPassword: "", expectedSSL: false},
		{uri: "rediss://localhost:6379", expectedAddress: "localhost:6379", expectedPassword: "", expectedSSL: false},
		{uri: "localhost:6379,password=1234,ssl=true", expectedAddress: "localhost:6379", expectedPassword: "1234", expectedSSL: true},
		{uri: "redis://localhost:6379,password=1234,ssl=true", expectedAddress: "localhost:6379", expectedPassword: "1234", expectedSSL: true},
		{uri: "rediss://localhost:6379,password=1234,ssl=true", expectedAddress: "localhost:6379", expectedPassword: "1234", expectedSSL: true},
	}

	for _, test := range tests {
		opts.Database.URI = test.uri
		opts.Database = fromURI(opts.Database)

		if opts.Database.Address != test.expectedAddress {
			t.Errorf("incorrect value, got: %s, want: %s", opts.Database.Address, test.expectedAddress)
		}

		if opts.Database.Password != test.expectedPassword {
			t.Errorf("incorrect value, got: %s, want: %s", opts.Database.Password, test.expectedPassword)
		}

		if opts.Database.SSL != test.expectedSSL {
			t.Errorf("incorrect value, got: %v, want: %v", opts.Database.SSL, test.expectedSSL)
		}
	}
}
