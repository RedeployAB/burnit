package server

func (s server) routes() {
	handler := s.httpServer.Handler
	if len(s.cors.Origin) > 0 {
		handler = corsHandler(s.cors.Origin)(handler)
	}
	s.httpServer.Handler = requestLogger(s.log, handler)

	s.router.Handle("GET /secret", s.generateSecret())
	s.router.Handle("GET /secrets/{id}", s.getSecret())
	s.router.Handle("GET /secrets/{id}/{passphrase}", s.getSecret())
	s.router.Handle("POST /secrets", s.createSecret())
}
