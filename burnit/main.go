package main

import (
	"log"

	"github.com/RedeployAB/burnit/burnit/config"
	"github.com/RedeployAB/burnit/burnit/server"
	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Configuration: %v.\n", err)
	}

	secrets, err := SetupSecretService(cfg)
	if err != nil {
		log.Fatalf("Setting up service: %v.\n", err)
	}
	srv := server.New(server.Options{
		Configuration: cfg,
		Router:        mux.NewRouter(),
		Secrets:       secrets,
	})

	srv.Start()
}
