package main

import (
	"flag"
	"log"

	"github.com/RedeployAB/burnit/burnitgw/config"
	"github.com/RedeployAB/burnit/burnitgw/server"
	"github.com/gorilla/mux"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	// Setup config.
	conf, err := config.Configure(*configPath)
	if err != nil {
		log.Fatalf("configuration: %v", err)
	}

	r := mux.NewRouter()
	srv := server.New(conf, r)
	// Start server.
	srv.Start()
}
