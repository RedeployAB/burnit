package app

import (
	"github.com/RedeployAB/burnit/common/middleware"
)

func (s *Server) routes() {
	// Routes.
	s.router.HandleFunc("/api/generate", s.generateSecret).Methods("GET")
	s.router.PathPrefix("/").HandlerFunc(s.notFound)
	s.router.Use(middleware.Logger)
}
