package main

import (
	"github.com/RedeployAB/burnit/secretgen/config"
	"github.com/RedeployAB/burnit/secretgen/server"
)

func main() {
	// Setup config.
	conf := config.Configure()
	srv := server.NewServer(conf)
	// Start server.
	srv.Start()
}
