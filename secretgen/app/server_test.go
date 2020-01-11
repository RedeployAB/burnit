package app

import (
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"

	"github.com/RedeployAB/burnit/secretgen/config"
)

func TestNewServer(t *testing.T) {
	conf := config.Configuration{Port: "5000"}
	r := mux.NewRouter()
	srv := NewServer(conf, r)

	if srv.httpServer.Addr != "0.0.0.0:5000" {
		t.Errorf("Incorrect value, got: %s, want: 0.0.0.0:5000", srv.httpServer.Addr)
	}
}

func TestStart(t *testing.T) {
	conf := config.Configuration{Port: "5000"}
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
		t.Errorf("Incorrect value, got: %d, want: -1", exitCode)
	}

}
