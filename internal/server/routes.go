package server

import (
	"net/http"

	"github.com/RedeployAB/burnit/internal/middleware"
	"github.com/RedeployAB/burnit/internal/ui"
)

const (
	// defaultContentSecurityPolicy is the default content security policy.
	defaultContentSecurityPolicy = "default-src 'self';"
)

// routes sets up the routes for the server.
func (s *server) routes() {
	s.httpServer.Handler = middleware.Chain(
		s.httpServer.Handler,
		middleware.RequestID(),
		middleware.SourceIP(),
		middleware.Logger(s.log),
	)

	middlewares, shutdownFuncs := setupMiddlewares(s.rateLimiter, s.cors)
	s.shutdownFuncs = append(s.shutdownFuncs, shutdownFuncs...)

	// Secret router and handlers.
	secretRouter := http.NewServeMux()
	secretRouter.Handle("GET /secret", generateSecret(s.secrets, s.log))

	// Secrets router and handlers.
	secretsRouter := http.NewServeMux()
	secretsRouter.Handle("GET /secrets/{id}", getSecret(s.secrets, s.log))
	secretsRouter.Handle("POST /secrets", createSecret(s.secrets, s.log))
	secretsRouter.Handle("DELETE /secrets/{id}", deleteSecret(s.secrets, s.log))

	secretHandler := middleware.Chain(secretRouter, middlewares...)
	secretsHandler := middleware.Chain(secretsRouter, middlewares...)

	s.router.Handle("/secret", secretHandler)
	s.router.Handle("/secrets", secretsHandler)

	if s.ui == nil {
		s.router.Handle("/{$}", index(nil, s.log))
		s.router.Handle("/", notFound(nil))
		return
	}

	uiMiddlewares := setupUIMiddlewares(s.ui.RuntimeRender())

	fer := http.NewServeMux()
	fer.Handle("/ui/secrets", ui.CreateSecret(s.ui, s.secrets, s.sessions))
	fer.Handle("/ui/secrets/", ui.GetSecret(s.ui, s.secrets, s.sessions, s.log))
	fer.Handle("/ui/handlers/secret/get", middleware.HTMX(ui.HandlerGetSecret(s.ui, s.secrets, s.sessions, s.log)))
	fer.Handle("/ui/handlers/secret/create", middleware.HTMX(ui.HandlerCreateSecret(s.ui, s.secrets, s.sessions, s.log)))
	fer.Handle("/ui/", ui.NotFound(s.ui))

	uiHandler := middleware.Chain(fer, uiMiddlewares...)
	s.router.Handle("/ui/", uiHandler)

	s.router.Handle("/static/", http.StripPrefix("/static/", ui.FileServer(s.ui.Static())))
	s.router.Handle("/{$}", index(s.ui, s.log))
	s.router.Handle("/", notFound(s.ui))
}

// setupMiddlewares sets up the middlewares for the server.
func setupMiddlewares(rl RateLimiter, c CORS) ([]middleware.Middleware, []func() error) {
	middlewares := []middleware.Middleware{}
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
	middlewares = append(middlewares, middleware.Headers())
	return middlewares, shutdownFuncs
}

// setupUIMiddlewares sets up the middlewares for the UI.
func setupUIMiddlewares(runtimeRender bool) []middleware.Middleware {
	middlewares := []middleware.Middleware{}
	var contentSecurityPolicy string
	if !runtimeRender {
		contentSecurityPolicy = defaultContentSecurityPolicy
	}
	middlewares = append(middlewares, middleware.Headers(func(o *middleware.HeadersOptions) {
		o.ContentSecurityPolicy = contentSecurityPolicy
	}))
	middlewares = append(middlewares, middleware.Compress())
	return middlewares
}
