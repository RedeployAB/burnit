package server

import (
	"github.com/RedeployAB/burnit/common/auth"
	"github.com/RedeployAB/burnit/common/middleware"
)

func (s *Server) routes(ts auth.TokenStore) {
	s.router.Handle("/secrets/{id}", s.getSecret()).Methods("GET")
	s.router.Handle("/secrets", s.createSecret()).Methods("POST")
	s.router.Handle("/secrets/{id}", s.updateSecret()).Methods("PUT")
	s.router.Handle("/secrets/{id}", s.deleteSecret()).Methods("DELETE")
	s.router.PathPrefix("/").HandlerFunc(s.notFound)

	s.router.Use(middleware.Logger)
	if s.tokenStore.(*auth.MemoryTokenStore) != nil {
		amw := middleware.Authentication{TokenStore: ts}
		s.router.Use(amw.Authenticate)
	}
}
