package server

import (
	"net/http"
	"strings"
)

// loggingResponseWriter is a wrapper around an http.ResponseWriter that keeps
// track of the status code and length of the response.
type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	length int
}

// WriteHeader acts as an adapter for the ResponseWriter's WriteHeader method,
// and also keeps track of the status code.
func (w *loggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Write acts as an adapter for the ResponseWriter's Write method,
// and also keeps track of the status code and length of the response.
func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

// requestLogger is a middleware that logs the incoming request.
func requestLogger(log logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := &loggingResponseWriter{ResponseWriter: w}
			next.ServeHTTP(lw, r)
			log.Info("Request received.", "status", lw.status, "path", maskSecretHash(r.URL.Path), "method", r.Method, "remoteIp", resolveIP(r))
		})
	}
}

// maskSecretHash masks the hash in the path.
func maskSecretHash(path string) string {
	lastIndex := strings.LastIndex(path, "/")
	if strings.HasPrefix(path, "/secrets/") && lastIndex == 45 || strings.HasPrefix(path, "/ui/secrets/") && lastIndex == 48 {
		return path[:lastIndex] + "/***"
	}
	return path
}
