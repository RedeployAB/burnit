package main

import (
	"github.com/RedeployAB/redeploy-secrets/secretgen/config"
	"github.com/RedeployAB/redeploy-secrets/secretgen/server"
)

func main() {
	// Setup config.
	conf := config.Configure()
	srv := server.NewServer(conf)
	// Start server.
	srv.Start()
}
