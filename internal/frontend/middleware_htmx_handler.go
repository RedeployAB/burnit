package frontend

import "net/http"

// HTMXHandler is a middleware that ensures the request is an htmx request.
// If it is not, the request is redirected to the root.
func HTMXHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Hx-Request") != "true" {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
