package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secret-generator/api"
	"github.com/RedeployAB/redeploy-secrets/secret-generator/config"
	"github.com/gorilla/mux"
)

var conf config.Config

func init() {
	conf = config.Configure()
	fmt.Printf("Server listening on: %s\n", conf.Port)
}

func main() {
	router := mux.NewRouter()
	// Routes.
	router.HandleFunc("/api/v1/secret", api.HandleGenerateSecret).Methods("GET")
	// Start server.
	log.Fatal(http.ListenAndServe(":"+conf.Port, router))
}
