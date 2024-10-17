package server

import "net/http"

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

	s.router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	s.router.Handle("/", s.uiCreateSecret())
	s.router.Handle("/ui/secrets/", s.uiGetSecret())

	// HTMX routes and handlers.
	s.router.Handle("/handlers/secret/create", htmxHandler(s.handlerCreateSecret()))
}

// htmxHandler is a middleware that ensures the request is an htmx request.
// If it is not, the request is redirected to the root.
func htmxHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Hx-Request") != "true" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
