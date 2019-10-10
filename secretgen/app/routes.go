package app

import (
	"github.com/RedeployAB/burnit/common/middleware"
)

var apiVer = "v0"

func (s *Server) routes() {
	// Routes.
	s.router.HandleFunc("/api/"+apiVer+"/generate", s.generateSecret).Methods("GET")
	s.router.PathPrefix("/").HandlerFunc(s.notFound)
	s.router.Use(middleware.Logger)
}
