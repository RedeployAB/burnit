package main

import (
	"github.com/RedeployAB/redeploy-secrets/secretgen/api"
	"github.com/RedeployAB/redeploy-secrets/secretgen/config"
	"github.com/RedeployAB/redeploy-secrets/secretgen/server"
)

var apiVer = "v1"

func main() {
	// Setup config.
	conf := config.Configure()
	r := api.NewRouter()
	srv := server.NewServer(conf, r)
	// Start server.
	srv.Start()
}
