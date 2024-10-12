package server

func (s *server) routes() {
	handler := s.httpServer.Handler

	if !s.cors.isEmpty() {
		handler = corsHandler(s.cors.Origin)(handler)
	}

	if !s.rateLimiter.isEmpty() {
		var closeRateLimiter func() error
		handler, closeRateLimiter = rateLimitHandler(
			withRateLimiterRate(s.rateLimiter.Rate),
			withRateLimiterBurst(s.rateLimiter.Burst),
			withRateLimiterCleanupInterval(s.rateLimiter.CleanupInterval),
			withRateLimiterTTL(s.rateLimiter.TTL),
		)(handler)
		s.shutdownFuncs = append(s.shutdownFuncs, closeRateLimiter)
	}

	s.httpServer.Handler = requestLogger(handler, s.log)

	s.router.Handle("GET /secret", s.generateSecret())
	s.router.Handle("GET /secrets/{id}", s.getSecret())
	s.router.Handle("GET /secrets/{id}/{passphrase}", s.getSecret())
	s.router.Handle("POST /secrets", s.createSecret())
}
