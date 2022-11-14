package main

import (
	"log"

	"github.com/RedeployAB/burnit/burnit/config"
	"github.com/RedeployAB/burnit/burnit/db"
	"github.com/RedeployAB/burnit/burnit/secret"
	"github.com/RedeployAB/burnit/burnit/server"
	"github.com/gorilla/mux"
)

func main() {
	conf, err := config.New()
	if err != nil {
		log.Fatalf("configuration: %v", err)
	}

	conn := connectToDB(conf)

	r := mux.NewRouter()
	srv := server.New(server.Options{
		Config:   conf,
		Router:   r,
		DBClient: conn.client,
		Secrets: secret.NewService(
			conn.repository,
			secret.Options{
				EncryptionKey: conf.Server.Security.Encryption.Key,
			},
		),
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
		log.Fatalf("database: %v", err)
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
