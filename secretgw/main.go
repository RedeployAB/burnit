package main

import (
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secretgw/api"
	"github.com/RedeployAB/redeploy-secrets/secretgw/config"
)

var conf = config.Config

func main() {
	r := api.NewRouter()
	// Start server.
	log.Printf("server listening on: %s", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, r))
}
