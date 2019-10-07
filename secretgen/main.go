package main

import (
	"log"

	"github.com/RedeployAB/redeploy-secrets/secretgen/api"
	"github.com/RedeployAB/redeploy-secrets/secretgen/config"
	"github.com/RedeployAB/redeploy-secrets/secretgen/server"
)

var apiVer = "v1"

func main() {
	// Setup config.
	conf := config.Configure()
	r := api.NewRouter()
	srv := server.NewServer(conf, r)

	done := make(chan bool, 1)
	go srv.AddShutdownHook(done)
	// Start server.
	srv.Start()
	<-done
	log.Println("server has been stopped")
}
