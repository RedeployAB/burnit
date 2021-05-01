package main

import (
	"log"
	"net/http"

	"github.com/RedeployAB/burnit/burnitgen/config"
	"github.com/RedeployAB/burnit/burnitgen/server"
)

func main() {
	flags := config.ParseFlags()

	conf, err := config.Configure(flags)
	if err != nil {
		log.Fatalf("configuration: %v", err)
	}

	r := http.NewServeMux()
	srv := server.New(conf, r)

	srv.Start()
}
