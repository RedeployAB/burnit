package main

import (
	"github.com/RedeployAB/redeploy-secrets/secretgw/config"
	"github.com/RedeployAB/redeploy-secrets/secretgw/server"
)

func main() {
	// Setup config.
	conf := config.Configure()
	srv := server.NewServer(conf)
	// Start server.
	srv.Start()
}
