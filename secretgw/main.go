package main

import (
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secretgw/api"
	"github.com/RedeployAB/redeploy-secrets/secretgw/config"
	"github.com/gorilla/mux"
)

var apiVer = "v1"
var port = config.Config.Port

func main() {
	// Setup routes.
	r := mux.NewRouter()
	s := r.PathPrefix("/api/" + apiVer).Subrouter()
	// Routes.
	s.Handle("/secret", api.GenerateSecretHandler()).Methods("GET")
	s.Handle("/secrets/{id}", api.ReadSecretHandler()).Methods("GET")
	s.Handle("/secrets", api.CreateSecretHandler()).Methods("POST")
	r.PathPrefix("/").HandlerFunc(api.NotFoundHandler)

	log.Printf("server listening on: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
