package server

import (
	"context"
	"testing"

	"github.com/RedeployAB/burnit/burnitdb/config"
	"github.com/RedeployAB/burnit/burnitdb/db"
	"github.com/RedeployAB/burnit/common/auth"
	"github.com/gorilla/mux"
)

type mockClient struct {
}

func (c *mockClient) Connect(context.Context) error {
	return nil
}

func (c *mockClient) Disconnect(context.Context) error {
	return nil
}

func (c *mockClient) GetAddress() string {
	return ""
}

func (c *mockClient) Close(context.Context) error {
	return nil
}

func (c *mockClient) Database(name string) db.Database {
	return nil
}

type mockRepository struct {
}

func (r *mockRepository) Find(id string) (*db.Secret, error) {
	return &db.Secret{}, nil
}

func (r *mockRepository) Insert(s *db.Secret) (*db.Secret, error) {
	return &db.Secret{}, nil
}

func (r *mockRepository) Delete(id string) (int64, error) {
	return 0, nil
}

func (r *mockRepository) DeleteExpired() (int64, error) {
	return 0, nil
}

func SetupOptions() Options {
	conf := config.Configuration{
		Server: config.Server{
			Port: "5000",
			Security: config.Security{
				APIKey: "ABCDEF",
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

	client := &mockClient{}
	repo := &mockRepository{}
	ts := auth.NewMemoryTokenStore()
	ts.Set(conf.Server.Security.APIKey, "server")

	r := mux.NewRouter()
	srvOpts := Options{
		Config:     &conf,
		Router:     r,
		DBClient:   client,
		Repository: repo,
		TokenStore: ts,
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
