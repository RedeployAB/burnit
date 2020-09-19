package server

import (
	"github.com/RedeployAB/burnit/common/middleware"
)

func (s *Server) routes() {
	g := s.router.PathPrefix("/").Subrouter()
	g.Handle("/secret", s.generateSecret()).Methods("GET")

	d := s.router.PathPrefix("/").Subrouter()
	d.Handle("/secrets/{id}", s.getSecret()).Methods("GET")
	d.Handle("/secrets", s.createSecret()).Methods("POST")

	amw := middleware.AuthHeader{Token: s.middlewareConfig.dbAPIkey}
	hmw := middleware.HeaderStrip{Exceptions: []string{"X-Passphrase"}}
	d.Use(hmw.Strip, amw.AddAuthHeader)

	s.router.PathPrefix("/").HandlerFunc(s.notFound)
	s.router.Use(middleware.Logger)
}
