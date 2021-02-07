package main

import (
	"log"

	"github.com/RedeployAB/burnit/burnitdb/config"
	"github.com/RedeployAB/burnit/burnitdb/db"
	"github.com/RedeployAB/burnit/burnitdb/server"
	"github.com/RedeployAB/burnit/common/auth"
	"github.com/gorilla/mux"
)

func main() {
	flags := config.ParseFlags()
	conf, err := config.Configure(flags)
	if err != nil {
		log.Fatalf("configuration: %v", err)
	}

	var ts *auth.MemoryTokenStore
	if len(conf.Server.Security.APIKey) > 0 {
		ts = auth.NewMemoryTokenStore()
		ts.Set(conf.Server.Security.APIKey, "server")
	}

	conn := connectToDB(conf)

	r := mux.NewRouter()
	srv := server.New(server.Options{
		Config:     conf,
		Router:     r,
		DBClient:   conn.client,
		Repository: conn.repository,
		TokenStore: ts,
	})

	srv.Start()
}

type dbConnection struct {
	client     db.Client
	repository db.Repository
}

func connectToDB(conf *config.Configuration) *dbConnection {
	log.Printf("connecting to db, host: %s (driver: %s)...\n", conf.Database.Address, conf.Database.Driver)
	client, err := db.Connect(conf.Database)
	if err != nil {
		log.Printf("could not connect to %s", conf.Database.Address)
		log.Fatalf("error: %v", err)
	}
	log.Printf("connected to %s", client.GetAddress())
	repo := db.NewSecretRepository(
		client,
		&db.SecretRepositoryOptions{
			Driver: conf.Database.Driver,
		},
	)

	return &dbConnection{client: client, repository: repo}
}
