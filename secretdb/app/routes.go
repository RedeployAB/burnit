package app

import (
	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/common/middleware"
)

var apiVer = "v0"

func (s *Server) routes(ts auth.TokenStore) {
	// Setup sub eouter.
	sr := s.router.PathPrefix("/api/" + apiVer).Subrouter()
	// Routes.
	sr.Handle("/secrets/{id}", s.getSecret()).Methods("GET")
	sr.Handle("/secrets", s.createSecret()).Methods("POST")
	//sr.Handle("/secrets/{id}", s.updateSecret()).Methods("PUT")
	//sr.Handle("/secrets/{id}", s.deleteSecret()).Methods("DELETE")
	//sr.Handle("/maintenance/cleanup", s.deleteExpiredSecrets()).Methods("DELETE")
	// All other routes.
	s.router.PathPrefix("/").HandlerFunc(s.notFound)
	// Attach middleware.
	amw := middleware.Authentication{TokenStore: ts}
	s.router.Use(middleware.Logger, amw.Authenticate)
}
