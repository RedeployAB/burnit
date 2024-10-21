package server

import (
	"net/http"

	"github.com/RedeployAB/burnit/internal/server/ui"
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

	if s.backendOnly {
		return
	}

	s.router.Handle("/static/", http.FileServer(http.FS(ui.StaticFS)))
	s.router.Handle("/", ui.CreateSecret(s.secrets))
	s.router.Handle("/ui/secrets/", ui.GetSecret(s.secrets))

	// Handlers for htmx requests.
	s.router.Handle("/ui/handlers/secret/get", ui.HTMXHandler(ui.HandlerGetSecret(s.secrets)))
	s.router.Handle("/ui/handlers/secret/create", ui.HTMXHandler(ui.HandlerCreateSecret(s.secrets)))
}
