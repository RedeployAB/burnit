package middleware

import (
	"net/http"
	"strings"
)

// logger is an interface for logging.
type logger interface {
	Error(msg string, args ...any)
	Info(msg string, args ...any)
}

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

// LoggerOptions contains the options for the Logger middleware.
type LoggerOptions struct {
	Type string
}

// LoggerOption is a function that sets an option on the Logger middleware.
type LoggerOption func(o *LoggerOptions)

// Logger is a middleware that logs the incoming request.
func Logger(log logger, options ...LoggerOption) func(next http.Handler) http.Handler {
	opts := LoggerOptions{
		Type: "backend",
	}
	for _, option := range options {
		option(&opts)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			lw := &loggingResponseWriter{ResponseWriter: w}
			next.ServeHTTP(lw, r)
			log.Info("Request received.", "type", "request", "component", opts.Type, "status", lw.status, "path", maskSecretHash(r.URL.Path), "method", r.Method, "remoteIp", resolveIP(r))
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
