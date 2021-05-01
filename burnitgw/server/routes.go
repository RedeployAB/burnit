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

	hmw := middleware.HeaderStrip{Exceptions: []string{"Passphrase"}}
	d.Use(hmw.Strip)
	if len(s.middlewareConfig.dbAPIkey) > 0 {
		amw := middleware.AuthHeader{Token: s.middlewareConfig.dbAPIkey}
		d.Use(amw.AddAuthHeader)
	}

	s.router.PathPrefix("/").HandlerFunc(s.notFound)
	s.router.Use(middleware.Logger)
}
