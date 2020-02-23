package app

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/RedeployAB/burnit/burnitdb/config"
	"github.com/RedeployAB/burnit/burnitdb/db"
	"github.com/RedeployAB/burnit/burnitdb/internal/models"
	"github.com/RedeployAB/burnit/common/auth"
	"github.com/gorilla/mux"
)

type mockConnection struct {
}

func (c *mockConnection) Connect(context.Context) error {
	return nil
}

func (c *mockConnection) Disconnect(context.Context) error {
	return nil
}

func (c *mockConnection) Close(context.Context) error {
	return nil
}

func (c *mockConnection) Database(name string) db.Database {
	return nil
}

type mockRepository struct {
}

func (r *mockRepository) Find(id string) (*models.Secret, error) {
	return &models.Secret{}, nil
}

func (r *mockRepository) Insert(s *models.Secret) (*models.Secret, error) {
	return &models.Secret{}, nil
}

func (r *mockRepository) Delete(id string) (int64, error) {
	return 0, nil
}

func (r *mockRepository) DeleteExpired() (int64, error) {
	return 0, nil
}

func SetupOptions() ServerOptions {
	conf := config.Configuration{
		Server: config.Server{
			Port:       "5000",
			DBAPIKey:   "ABCDEF",
			Passphrase: "testphrase",
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

	connection := &mockConnection{}
	repo := &mockRepository{}
	ts := auth.NewMemoryTokenStore()
	ts.Set(conf.Server.DBAPIKey, "app")

	r := mux.NewRouter()
	srvOpts := ServerOptions{
		Config:     conf,
		Router:     r,
		Connection: connection,
		Repository: repo,
		TokenStore: ts,
	}

	return srvOpts
}

func TestNewServer(t *testing.T) {

	srvOpts := SetupOptions()
	srv := NewServer(srvOpts)

	if srv.httpServer.Addr != "0.0.0.0:5000" {
		t.Errorf("incorrect value, got: %s, want: 0.0.0.0:5000", srv.httpServer.Addr)
	}
}

func TestStartAndShutdown(t *testing.T) {
	srvOpts := SetupOptions()
	srv := NewServer(srvOpts)

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
		t.Errorf("incorrect value, got %d, want: -1", exitCode)
	}
}
