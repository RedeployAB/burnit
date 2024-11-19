package middleware

import (
	"compress/gzip"
	"net/http"
	"strings"
	"sync"
)

// compressResponseWriter is a wrapper around an http.ResponseWriter that compresses
// the response using gzip.
type compressResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

// Write the response using gzip.
func (w compressResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

// Close the gzip writer.
func (w compressResponseWriter) Close() error {
	return w.writer.Close()
}

// CompressOptions contains the options for the Compress middleware.
type CompressOptions struct {
	Pool *sync.Pool
}

// CompressOption is a function that sets an option on the Compress middleware.
type CompressOption func(o *CompressOptions)

// Compress is a middleware that compresses the response using gzip.
func Compress(options ...CompressOption) func(next http.Handler) http.Handler {
	opts := CompressOptions{}
	for _, option := range options {
		option(&opts)
	}

	pool := opts.Pool
	if pool == nil {
		pool = NewGzipWriterPool()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			gz := pool.Get().(*gzip.Writer)
			defer pool.Put(gz)

			gz.Reset(w)
			defer gz.Close()

			w.Header().Set("Content-Encoding", "gzip")
			cw := compressResponseWriter{ResponseWriter: w, writer: gz}
			next.ServeHTTP(cw, r)
		})
	}
}

// NewGzipWriterPool creates a new sync.Pool for gzip.Writer.
func NewGzipWriterPool() *sync.Pool {
	return &sync.Pool{
		New: func() any {
			return gzip.NewWriter(nil)
		},
	}
}

// WithGzipWriterPool sets the sync.Pool for gzip.Writer.
func WithGzipWriterPool(pool *sync.Pool) CompressOption {
	return func(o *CompressOptions) {
		o.Pool = pool
	}
}
