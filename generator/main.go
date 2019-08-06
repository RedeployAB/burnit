package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/generator/api"
	"github.com/RedeployAB/redeploy-secrets/generator/config"
	"github.com/gorilla/mux"
)

var conf config.Config
var apiVer = "v1"

func init() {
	conf = config.Configure()
	fmt.Printf("Server listening on: %s\n", conf.Port)
}

func main() {
	router := mux.NewRouter()
	// Routes.
	router.HandleFunc("/api/"+apiVer+"/secret", api.GenerateSecretHandler).Methods("GET")
	// Start server.
	log.Fatal(http.ListenAndServe(":"+conf.Port, router))
}
