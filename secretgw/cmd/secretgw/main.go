package main

import (
	"flag"
	"log"

	"github.com/RedeployAB/burnit/burnitgw/app"
	"github.com/RedeployAB/burnit/burnitgw/config"
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
	srv := app.NewServer(conf, r)
	// Start server.
	srv.Start()
}
