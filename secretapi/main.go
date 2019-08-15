package main

import (
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secretapi/api"
	"github.com/RedeployAB/redeploy-secrets/secretapi/config"
)

var conf = config.Config

func main() {
	r := api.NewRouter()
	// Start server.
	log.Printf("server listening on: %s", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, r))
}
