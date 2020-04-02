package main

import (
	"flag"
	"log"

	"github.com/RedeployAB/burnit/burnitdb/config"
	"github.com/RedeployAB/burnit/burnitdb/db"
	"github.com/RedeployAB/burnit/burnitdb/server"
	"github.com/RedeployAB/burnit/common/auth"
	"github.com/gorilla/mux"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Setup configuration.
	conf, err := config.Configure(*configPath)
	if err != nil {
		log.Fatalf("configuration: %v", err)
	}

	ts := auth.NewMemoryTokenStore()
	ts.Set(conf.Server.DBAPIKey, "server")

	// Connect to database.
	log.Printf("connecting to db server: %s...\n", conf.Database.Address)
	connection, err := db.Connect(conf.Database)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	repo := db.NewSecretRepository(
		connection,
		&db.SecretRepositoryOptions{
			EncryptionKey: conf.Server.Security.Encryption.Key,
			HashMethod:    conf.Server.Security.HashMethod,
		},
	)

	r := mux.NewRouter()
	srv := server.New(server.Options{
		Config:     conf,
		Router:     r,
		Connection: connection,
		Repository: repo,
		TokenStore: ts,
	})
	// Listen and serve.
	srv.Start()
}
