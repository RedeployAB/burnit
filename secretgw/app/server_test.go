package app

import (
	"os"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/burnitgw/config"
	"github.com/gorilla/mux"
)

func TestNewServer(t *testing.T) {

	conf := config.Configuration{Server: config.Server{Port: "5000"}}
	r := mux.NewRouter()
	srv := NewServer(conf, r)

	expected := "0.0.0.0:5000"
	if srv.httpServer.Addr != expected {
		t.Errorf("incorrect value, got: %s, want: %s", srv.httpServer.Addr, expected)
	}
}

func TestStartAndShutdown(t *testing.T) {
	conf := config.Configuration{Server: config.Server{Port: "5000"}}
	r := mux.NewRouter()
	srv := NewServer(conf, r)

	proc, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatal("error in getting running process")
	}

	var result *os.ProcessState

	go func() {
		srv.Start()
		result, _ = proc.Wait()
	}()

	time.Sleep(3 * time.Second)
	proc.Signal(os.Interrupt)

	exitCode := result.ExitCode()
	if exitCode != -1 {
		t.Errorf("incorrect value, got: %d, want: -1", exitCode)
	}
}
