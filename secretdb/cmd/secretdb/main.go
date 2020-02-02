package main

import (
	"flag"
	"log"

	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/secretdb/app"
	"github.com/RedeployAB/burnit/secretdb/config"
	"github.com/RedeployAB/burnit/secretdb/db"
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
	ts.Set(conf.Server.DBAPIKey, "app")

	// Connect to database.
	log.Printf("connecting to db server: %s...\n", conf.Database.Address)
	connection, err := db.Connect(conf.Database)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	repo := db.NewSecretRepository(connection, conf.Server.Passphrase)

	r := mux.NewRouter()
	srv := app.NewServer(app.ServerOptions{
		Config:     conf,
		Router:     r,
		Connection: connection,
		Repository: repo,
		TokenStore: ts,
	})
	// Listen and serve.
	srv.Start()
}
