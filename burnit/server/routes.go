package server

import (
	"github.com/RedeployAB/burnit/burnit/middleware"
)

func (s *server) routes() {
	// Generate a secret.
	s.router.Handle("/secret", s.generateSecret()).Methods("GET")
	// Get, create and delete a secret.
	s.router.Handle("/secrets/{id}", s.getSecret()).Methods("GET")
	s.router.Handle("/secrets", s.createSecret()).Methods("POST")
	s.router.Handle("/secrets/{id}", s.deleteSecret()).Methods("DELETE")

	s.router.PathPrefix("/").HandlerFunc(s.notFound)

	s.router.Use(middleware.Logger)
	if s.middleware.cors.enabled {
		corsHandler := middleware.CORSHandler{
			Origin:  s.middleware.cors.origin,
			Headers: s.middleware.cors.headers,
		}
		s.router.Use(corsHandler.Handle)
	}
}
