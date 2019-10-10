package app

import (
	"github.com/RedeployAB/burnit/common/middleware"
)

var genAPIVer = "v0"
var dbAPIVer = "v0"

func (s *Server) routes() {
	// Generator routes.
	g := s.router.PathPrefix("/api/" + genAPIVer).Subrouter()
	g.Handle("/generate", s.generateSecret()).Methods("GET")

	// DB routes.
	d := s.router.PathPrefix("/api/" + dbAPIVer).Subrouter()
	d.Handle("/secrets/{id}", s.getSecret()).Methods("GET")
	d.Handle("/secrets", s.createSecret()).Methods("POST")
	// Init middleware for all db routes.
	amw := middleware.AuthHeader{Token: s.middlewareConfig.dbAPIkey}
	hmw := middleware.HeaderStrip{Exceptions: []string{"X-Passphrase"}}
	d.Use(hmw.Strip, amw.AddAuthHeader)
	// All other routes.
	s.router.PathPrefix("/").HandlerFunc(s.notFound)
	// Attach logger.
	s.router.Use(middleware.Logger)
}
