package middleware

import (
	"log"
	"net/http"
	"time"
)

// Inspiration for this way of intercepting and logging request was inspired/taken
// from a reply by nemith at https://www.reddit.com/r/golang/comments/7p35s4/how_do_i_get_the_response_status_for_my_middleware/
// and a reply by huangapple at https://stackoverflow.com/questions/53272536/how-do-i-get-response-statuscode-in-golang-middleware.
// So credit to those posts and posters.

// logginResponseWriter wraps around http.ResponseWriter
// with additional fields for status and length.
type loggingResponseWriter struct {
	http.ResponseWriter
	status int
	lenght int
}

// WriteHeader writes the status code to the status field on
// and to the headers of the writer.
func (w *loggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Write wraps around the ResponseWriters Write
// for additional handling. default status and
// setting the length field.
func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.lenght += n
	return n, err
}

// Logger is a middleware to handle logging for the server.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(&lw, r)
		log.Printf("%d\t%s\t%s\t%s", lw.status, r.Method, r.RequestURI, time.Since(start))
	})
}
