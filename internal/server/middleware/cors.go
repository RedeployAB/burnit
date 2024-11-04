package middleware

import (
	"net/http"
)

var (
	// corsAllowMethods is the allowed methods for CORS.
	corsAllowMethods = "GET, POST, DELETE"
	// corsAllowHeaders is the allowed headers for CORS.
	corsAllowHeaders = "Content-Type, Passphrase"
)

// CORS returns a CORS middleware http.Handler.
func CORS(origin string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", corsAllowMethods)
			w.Header().Set("Access-Control-Allow-Headers", corsAllowHeaders)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
