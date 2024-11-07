package middleware

import "net/http"

// CacheControlOptions contains the options for the CacheControl middleware.
type CacheControlOptions struct {
	NoStore bool
}

// CacheControlOption is a function that sets an option on the CacheControl middleware.
type CacheControlOption func(o *CacheControlOptions)

// CacheControl is a middleware that sets the Cache-Control header.
func CacheControl(options ...CacheControlOption) func(next http.Handler) http.Handler {
	opts := CacheControlOptions{
		NoStore: true,
	}
	for _, option := range options {
		option(&opts)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if opts.NoStore {
				w.Header().Set("Cache-Control", "no-store")
			}
		})
	}
}
