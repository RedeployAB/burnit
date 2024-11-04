package server

import (
	"net/http"

	"github.com/RedeployAB/burnit/internal/frontend"
	"github.com/RedeployAB/burnit/internal/middleware"
)

// routes sets up the routes for the server.
func (s *server) routes() {
	r := http.NewServeMux()
	r.Handle("GET /secret", s.generateSecret())
	r.Handle("GET /secrets/", s.getSecret())
	r.Handle("POST /secrets", s.createSecret())
	r.Handle("DELETE /secrets/", s.deleteSecret())

	var middlewares []func(http.Handler) http.Handler
	if !s.rateLimiter.isEmpty() {
		mw, closeRateLimiter := rateLimitHandler(
			withRateLimiterRate(s.rateLimiter.Rate),
			withRateLimiterBurst(s.rateLimiter.Burst),
			withRateLimiterTTL(s.rateLimiter.TTL),
			withRateLimiterCleanupInterval(s.rateLimiter.CleanupInterval),
		)
		middlewares = append(middlewares, mw)
		s.shutdownFuncs = append(s.shutdownFuncs, closeRateLimiter)
	}
	if !s.cors.isEmpty() {
		middlewares = append(middlewares, corsHandler(s.cors.Origin))
	}
	middlewares = append(middlewares, requestLogger(s.log))

	handler := middleware.Chain(r, middlewares...)

	s.router.Handle("/secret", handler)
	s.router.Handle("/secrets", handler)
	s.router.Handle("/secrets/", handler)

	if s.ui == nil {
		return
	}

	s.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(s.ui.Static()))))

	fer := http.NewServeMux()
	fer.Handle("/", frontend.CreateSecret(s.ui, s.secrets))
	fer.Handle("/ui/secrets/", frontend.GetSecret(s.ui, s.secrets))
	fer.Handle("/ui/handlers/secret/get", frontend.HTMXHandler(frontend.HandlerGetSecret(s.ui, s.secrets)))
	fer.Handle("/ui/handlers/secret/create", frontend.HTMXHandler(frontend.HandlerCreateSecret(s.ui, s.secrets)))

	s.router.Handle("/", fer)
	s.router.Handle("/ui/", fer)
}
