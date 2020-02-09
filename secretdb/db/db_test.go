package db

import (
	"testing"

	"github.com/RedeployAB/burnit/secretdb/config"
)

func TestToURI(t *testing.T) {
	dbConf1 := config.Database{
		Address:  "hostname",
		Database: "db",
		Username: "user",
		Password: "password",
		SSL:      true,
	}
	dbURI1 := "mongodb://user:password@hostname/?ssl=true"

	dbConf2 := config.Database{
		Address:  "hostname",
		Database: "db",
		Username: "user",
	}
	dbURI2 := "mongodb://user@hostname"

	dbConf3 := config.Database{
		Address: "hostname",
	}
	dbURI3 := "mongodb://hostname"

	tests := []struct {
		input config.Database
		want  string
	}{
		{dbConf1, dbURI1},
		{dbConf2, dbURI2},
		{dbConf3, dbURI3},
	}

	for _, test := range tests {
		got := toURI(test.input)
		if got != test.want {
			t.Errorf("incorrect value, got: %s, want: %s", got, test.want)
		}
	}
}
