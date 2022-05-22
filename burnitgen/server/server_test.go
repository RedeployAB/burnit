package server

import (
	"net/http"
	"testing"

	"github.com/RedeployAB/burnit/burnitgen/config"
)

func TestNew(t *testing.T) {
	conf := &config.Configuration{Port: "5000"}
	r := http.NewServeMux()
	srv := New(conf, r)

	expected := "0.0.0.0:5000"
	if srv.httpServer.Addr != expected {
		t.Errorf("incorrect value, got: %s, want: %s", srv.httpServer.Addr, expected)
	}
}
