package server

import (
	"github.com/RedeployAB/burnit/common/middleware"
)

func (s *Server) routes() {
	s.router.HandleFunc("/secret", s.generateSecret).Methods("GET")
	s.router.PathPrefix("/").HandlerFunc(s.notFound)
	s.router.Use(middleware.Logger)
}
