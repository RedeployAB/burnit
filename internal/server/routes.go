package server

import (
	"net/http"

	"github.com/RedeployAB/burnit/internal/frontend"
	"github.com/RedeployAB/burnit/internal/middleware"
)

// routes sets up the routes for the server.
func (s *server) routes() {
	middlewares, shutdownFuncs := setupMiddlewares(s.log, s.rateLimiter, s.cors)
	s.shutdownFuncs = append(s.shutdownFuncs, shutdownFuncs...)

	// Secret router and handlers.
	secretRouter := http.NewServeMux()
	secretRouter.Handle("GET /secret", s.generateSecret())

	// Secrets router and handlers.
	secretsRouter := http.NewServeMux()
	secretsRouter.Handle("GET /secrets", s.badRequestSecrets())
	secretsRouter.Handle("GET /secrets/", s.getSecret())
	secretsRouter.Handle("POST /secrets", s.createSecret())
	secretsRouter.Handle("DELETE /secrets/", s.deleteSecret())

	secretHandler := middleware.Chain(secretRouter, middlewares...)
	secretsHandler := middleware.Chain(secretsRouter, middlewares...)

	s.router.Handle("/secret", secretHandler)
	s.router.Handle("/secrets", secretsHandler)

	if s.ui == nil {
		s.router.Handle("/{$}", s.index())
		s.router.Handle("/", s.notFound())
		return
	}

	frontendMiddlewares := setupFrontendMiddlewares(s.log)

	fer := http.NewServeMux()
	fer.Handle("/ui/secrets", frontend.CreateSecret(s.ui, s.secrets))
	fer.Handle("/ui/secrets/", frontend.GetSecret(s.ui, s.secrets))
	fer.Handle("/ui/handlers/secret/get", middleware.HTMX(frontend.HandlerGetSecret(s.ui, s.secrets)))
	fer.Handle("/ui/handlers/secret/create", middleware.HTMX(frontend.HandlerCreateSecret(s.ui, s.secrets)))
	fer.Handle("/ui/", frontend.NotFound(s.ui))

	frontendHandler := middleware.Chain(fer, frontendMiddlewares...)

	s.router.Handle("/ui/", frontendHandler)

	s.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(s.ui.Static()))))
	s.router.Handle("/{$}", s.index())
	s.router.Handle("/", s.notFound())
}

// setupMiddlewares sets up the middlewares for the server.
func setupMiddlewares(log logger, rl RateLimiter, c CORS) ([]middleware.Middleware, []func() error) {
	middlewares := []middleware.Middleware{middleware.Logger(log)}
	var shutdownFuncs []func() error
	if !rl.isEmpty() {
		mw, closeRateLimiter := middleware.RateLimiter(
			middleware.WithRateLimiterRate(rl.Rate),
			middleware.WithRateLimiterBurst(rl.Burst),
			middleware.WithRateLimiterTTL(rl.TTL),
			middleware.WithRateLimiterCleanupInterval(rl.CleanupInterval),
		)
		middlewares = append(middlewares, mw)
		shutdownFuncs = append(shutdownFuncs, closeRateLimiter)
	}
	if !c.isEmpty() {
		middlewares = append(middlewares, middleware.CORS(c.Origin))
	}
	return middlewares, shutdownFuncs
}

// setupFrontendMiddlewares sets up the middlewares for the frontend.
func setupFrontendMiddlewares(log logger) []middleware.Middleware {
	return []middleware.Middleware{middleware.Logger(log, func(o *middleware.LoggerOptions) {
		o.Type = "frontend"
	})}
}
