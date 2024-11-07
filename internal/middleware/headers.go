package middleware

import "net/http"

// HeadersOptions represents the options for the Headers middleware.
type HeadersOptions struct {
	CacheControl          bool
	ContentSecurityPolicy bool
	TLS                   bool
}

// HeadersOption is a function that sets an option for the Headers middleware.
type HeadersOption func(o *HeadersOptions)

// Headers is a middleware that sets the security headers.
func Headers(options ...HeadersOption) func(next http.Handler) http.Handler {
	opts := HeadersOptions{}
	for _, option := range options {
		option(&opts)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")

			if opts.CacheControl {
				w.Header().Set("Cache-Control", "no-store")
			}
			if opts.ContentSecurityPolicy {
				w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self';")
			}
			if opts.TLS {
				w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
			}
		})
	}
}
