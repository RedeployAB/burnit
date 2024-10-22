package server

import (
	"net/http"

	"github.com/RedeployAB/burnit/internal/frontend"
)

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
			withRateLimiterTTL(s.rateLimiter.TTL),
			withRateLimiterCleanupInterval(s.rateLimiter.CleanupInterval),
		)(handler)
		s.shutdownFuncs = append(s.shutdownFuncs, closeRateLimiter)
	}

	s.httpServer.Handler = requestLogger(handler, s.log)

	s.router.Handle("GET /secret", s.generateSecret())
	s.router.Handle("GET /secrets/", s.getSecret())
	s.router.Handle("POST /secrets", s.createSecret())
	s.router.Handle("DELETE /secrets/", s.deleteSecret())

	if s.ui == nil {
		return
	}

	s.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(s.ui.Static()))))
	s.router.Handle("/", frontend.CreateSecret(s.ui, s.secrets))
	s.router.Handle("/ui/secrets/", frontend.GetSecret(s.ui, s.secrets))

	s.router.Handle("/ui/handlers/secret/get", frontend.HTMXHandler(frontend.HandlerGetSecret(s.ui, s.secrets)))
	s.router.Handle("/ui/handlers/secret/create", frontend.HTMXHandler(frontend.HandlerCreateSecret(s.ui, s.secrets)))
}
