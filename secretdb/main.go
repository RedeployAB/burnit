package main

import (
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secretdb/api"
	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
	"github.com/RedeployAB/redeploy-secrets/secretdb/db"
)

var apiVer = "v1"
var conf = config.Config

func main() {
	// Connect to db
	db.Connect(config.Config.Database)
	session, err := db.GetSession()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer session.Close()
	// Connect to collection.
	collection := session.DB("secretdb").C("secrets")
	// Create new router.
	r := api.NewRouter(conf.Server, collection)
	// Start server.
	log.Printf("server listening on: %s", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, r))
}
