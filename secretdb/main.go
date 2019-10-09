package main

import (
	"log"

	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/secretdb/config"
	"github.com/RedeployAB/burnit/secretdb/db"
	"github.com/RedeployAB/burnit/secretdb/server"
)

func main() {
	// Setup TokenStore.
	conf := config.Configure()
	ts := auth.NewMemoryTokenStore()
	ts.Set(conf.Server.DBAPIKey, "app")

	log.Printf("connecting to db server: %s...\n", conf.Database.Address)
	client, err := db.Connect(conf.Database)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	srv := server.NewServer(conf, ts, client)
	// Start server.
	srv.Start()
}
