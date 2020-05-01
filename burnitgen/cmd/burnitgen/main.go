package main

import (
	"log"

	"github.com/RedeployAB/burnit/burnitgen/config"
	"github.com/RedeployAB/burnit/burnitgen/server"
	"github.com/gorilla/mux"
)

func main() {
	flags := config.ParseFlags()
	// Setup config.
	conf, err := config.Configure(flags)
	if err != nil {
		log.Fatalf("configuration: %v", err)
	}

	r := mux.NewRouter()
	srv := server.New(conf, r)
	// Start server.
	srv.Start()
}
