package main

import (
	"github.com/RedeployAB/burnit/secretgen/app"
	"github.com/RedeployAB/burnit/secretgen/config"
	"github.com/gorilla/mux"
)

func main() {
	// Setup config.
	conf := config.Configure()

	r := mux.NewRouter()
	srv := app.NewServer(conf, r)
	// Start server.
	srv.Start()
}
