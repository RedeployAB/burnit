package main

import (
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secretgen/api"
	"github.com/RedeployAB/redeploy-secrets/secretgen/config"
)

var apiVer = "v1"
var conf = config.Config

func main() {
	r := api.NewRouter(conf.Server)
	// Start server.
	log.Printf("server listening on: %s", conf.Port)
	log.Fatal(http.ListenAndServe(":"+conf.Port, r))
}
