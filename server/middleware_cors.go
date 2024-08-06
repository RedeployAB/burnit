package server

import (
	"net/http"
	"strings"
)

// newCORSHandler creates a new CORS handler.
func newCORSHandler(origin string, headers http.Header) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			for k, v := range headers {
				w.Header().Set(k, strings.Join(v, ", "))
			}
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
