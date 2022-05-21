package server

import (
	"testing"

	"github.com/RedeployAB/burnit/burnitgw/config"
	"github.com/gorilla/mux"
)

func TestNewServer(t *testing.T) {

	conf := &config.Configuration{Server: config.Server{Port: "5000"}}
	r := mux.NewRouter()
	srv := New(conf, r)

	expected := "0.0.0.0:5000"
	if srv.httpServer.Addr != expected {
		t.Errorf("incorrect value, got: %s, want: %s", srv.httpServer.Addr, expected)
	}
}
