package server

import (
	"testing"

	"github.com/RedeployAB/burnit/burnit/config"
	"github.com/gorilla/mux"
)

func SetupOptions() Options {
	conf := config.Configuration{
		Server: config.Server{
			Host: "0.0.0.0",
			Port: "5000",
			Security: config.Security{
				Encryption: config.Encryption{
					Key: "testphrase",
				},
			},
		},
		Database: config.Database{
			Address:  "mongo://db",
			Database: "secrets",
			Username: "user",
			Password: "password",
			SSL:      false,
			URI:      "",
		},
	}

	r := mux.NewRouter()
	srvOpts := Options{
		Configuration: &conf,
		Router:        r,
	}

	return srvOpts
}

func TestNew(t *testing.T) {

	srvOpts := SetupOptions()
	srv := New(srvOpts)

	if srv.httpServer.Addr != "0.0.0.0:5000" {
		t.Errorf("incorrect value, got: %s, want: 0.0.0.0:5000", srv.httpServer.Addr)
	}
}
