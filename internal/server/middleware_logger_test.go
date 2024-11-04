package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRequestLogger(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			status int
			req    func() *http.Request
		}
		want []string
	}{
		{
			name: "log requests with status OK",
			input: struct {
				status int
				req    func() *http.Request
			}{
				status: http.StatusOK,
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req.Header.Set("Forwarded", "for=192.168.1.1:1234")
					return req
				},
			},
			want: []string{"Request received.", "status", "200", "path", "/", "method", "GET", "remoteIp", "192.168.1.1"},
		},
		{
			name: "log requests with status OK (no status)",
			input: struct {
				status int
				req    func() *http.Request
			}{
				status: 0,
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req.Header.Set("Forwarded", "for=192.168.1.1:1234")
					return req
				},
			},
			want: []string{"Request received.", "status", "200", "path", "/", "method", "GET", "remoteIp", "192.168.1.1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logs := []string{}
			log := &mockLogger{
				logs: &logs,
			}

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if test.input.status != 0 {
					w.WriteHeader(test.input.status)
				}
				w.Write([]byte("response"))
			})

			rr := httptest.NewRecorder()
			req := test.input.req()
			requestLogger(log)(handler).ServeHTTP(rr, req)

			if diff := cmp.Diff(test.want, logs); diff != "" {
				t.Errorf("requestLogger() = unexpected result, (-want, +got):\n%s\n", diff)
			}
		})
	}
}
