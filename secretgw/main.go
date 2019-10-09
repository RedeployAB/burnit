package main

import (
	"github.com/RedeployAB/burnit/secretgw/config"
	"github.com/RedeployAB/burnit/secretgw/server"
)

func main() {
	// Setup config.
	conf := config.Configure()
	srv := server.NewServer(conf)
	// Start server.
	srv.Start()
}
