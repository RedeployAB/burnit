package app

import (
	"context"
	"testing"

	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/secretdb/config"
	"github.com/RedeployAB/burnit/secretdb/internal/dto"
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

type mockRepository struct {
}

func (r *mockRepository) Find(id string) (*dto.Secret, error) {
	return &dto.Secret{}, nil
}

func (r *mockRepository) Delete(id string) (int64, error) {
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

	mockConnection := &mockConnection{}
	mockRepoository := &mockRepository{}
	ts := auth.NewMemoryTokenStore()
	ts.Set(conf.Server.DBAPIKey, "app")

	r := mux.NewRouter()
	srvOpts := ServerOptions{
		Config:     conf,
		Router:     r,
		Connection: mockConnection,
		Repository: mockRepoository,
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
