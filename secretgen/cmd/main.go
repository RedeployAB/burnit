package main

import (
	"github.com/RedeployAB/burnit/secretgen/app"
	"github.com/RedeployAB/burnit/secretgen/configs"
	"github.com/gorilla/mux"
)

func main() {
	// Setup config.
	config := configs.Configure()

	r := mux.NewRouter()
	srv := app.NewServer(config, r)
	// Start server.
	srv.Start()
}
