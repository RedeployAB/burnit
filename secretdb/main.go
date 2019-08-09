package main

import (
	"log"
	"net/http"

	"github.com/RedeployAB/redeploy-secrets/secretdb/api"
	"github.com/RedeployAB/redeploy-secrets/secretdb/config"
	"github.com/RedeployAB/redeploy-secrets/secretdb/db"
	"github.com/RedeployAB/redeploy-secrets/secretdb/middleware"
	"github.com/gorilla/mux"
)

var conf config.Config
var apiVer = "v1"

func init() {
	conf = config.Configure()
	log.Printf("Server listening on: %s\n", conf.Port)
}

func main() {
	// Connect to db
	db.Connect(conf.Database)
	session, err := db.GetSession()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer session.Close()
	// Connect to collection.
	collection := session.DB("secretdb").C("secrets")

	// TODO: Move this into a seperate file?
	r := mux.NewRouter()
	s := r.PathPrefix("/api/" + apiVer).Subrouter()
	// Routes.
	s.Handle("/secrets/{id}", middleware.Chain(api.ReadSecretHandler(collection))).Methods("GET")
	s.Handle("/secrets", middleware.Chain(api.CreateSecretHandler(collection))).Methods("POST")
	s.Handle("/secrets/{id}", middleware.Chain(api.UpdateSecretHandler(collection))).Methods("PUT")
	s.Handle("/secrets/{id}", middleware.Chain(api.DeleteSecretHandler(collection))).Methods("DELETE")
	// All other routes.
	r.PathPrefix("/").HandlerFunc(api.NotFoundHandler)
	// Attach logger.
	r.Use(middleware.Logger)
	// Start server.
	log.Fatal(http.ListenAndServe(":"+conf.Server.Port, r))
}
