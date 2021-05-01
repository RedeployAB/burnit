package server

import (
	"github.com/RedeployAB/burnit/common/middleware"
)

func (s *Server) routes() {
	s.router.Handle("/secret", middleware.Logger(s.generateSecret()))
}
