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
	// Connect to db
	db.Connect(conf.Database)
	session, err := db.GetSession()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer session.Close()
	// Connect to collection.
	collection := session.DB("secretdb").C("secrets")
	// Create new router.
	r := api.NewRouter(ts, collection)
	// Start server.
	log.Printf("server listening on: %s", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, r))
}
