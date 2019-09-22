package main

import (
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/common/auth"
	"github.com/RedeployAB/redeploy-secrets/secretdb/api"
	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
	"github.com/RedeployAB/redeploy-secrets/secretdb/db"
)

var apiVer = "v1"
var conf = config.Config

func main() {
	// Setup TokenStore.
	ts := auth.NewMemoryTokenStore()
	ts.Set(conf.DBAPIKey, "app")

	log.Printf("connecting to db server: %s...\n", conf.Database.Address)
	client, err := db.Connect(conf.Database)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("connected.\n")

	r := api.NewRouter(ts, client)
	// Start server.
	log.Printf("server listening on: %s", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, r))
}
