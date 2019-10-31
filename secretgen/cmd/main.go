package main

import (
	"flag"
	"log"

	"github.com/RedeployAB/burnit/secretgen/app"
	"github.com/RedeployAB/burnit/secretgen/config"
	"github.com/gorilla/mux"
)

func main() {
	configPath := flag.String("config", "", "path to configuration file")
	flag.Parse()

	// Setup config.
	conf, err := config.Configure(*configPath)
	if err != nil {
		log.Fatalf("%v", err)
	}

	r := mux.NewRouter()
	srv := app.NewServer(conf, r)
	// Start server.
	srv.Start()
}
