package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogger(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	// Writing nothing.
	h1 := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if w.(*loggingResponseWriter).status != 0 {
				t.Errorf("incorrect value, got: %d, want: 0", w.(*loggingResponseWriter).status)
			}
		})
	}
	// Writing header.
	h2 := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.(*loggingResponseWriter).WriteHeader(200)
			if w.(*loggingResponseWriter).status != 200 {
				t.Errorf("incorrect value, got: %d, want: 200", w.(*loggingResponseWriter).status)
			}
		})
	}
	// Writing content.
	h3 := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.(*loggingResponseWriter).Write([]byte("content"))
			if w.(*loggingResponseWriter).status != 200 {
				t.Errorf("incorrect value, got: %d, want: 200", w.(*loggingResponseWriter).status)
			}
		})
	}

	Logger(h1()).ServeHTTP(res, req)
	Logger(h2()).ServeHTTP(res, req)
	Logger(h3()).ServeHTTP(res, req)
}
