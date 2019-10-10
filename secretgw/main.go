package main

import (
	"github.com/RedeployAB/burnit/secretgw/app"
	"github.com/RedeployAB/burnit/secretgw/config"
	"github.com/gorilla/mux"
)

func main() {
	// Setup config.
	config := config.Configure()

	r := mux.NewRouter()
	srv := app.NewServer(config, r)
	// Start server.
	srv.Start()
}
