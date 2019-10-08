package main

import (
	"log"

	"github.com/RedeployAB/redeploy-secrets/common/auth"
	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
	"github.com/RedeployAB/redeploy-secrets/secretdb/db"
	"github.com/RedeployAB/redeploy-secrets/secretdb/server"
)

var apiVer = "v1"

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
