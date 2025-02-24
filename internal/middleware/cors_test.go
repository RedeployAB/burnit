package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCORS(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			origin string
			req    *http.Request
		}
		want struct {
			headers    http.Header
			statusCode int
		}
	}{
		{
			name: "CORS handler with valid origin and headers",
			input: struct {
				origin string
				req    *http.Request
			}{
				origin: "http://localhost:3000",
				req:    httptest.NewRequest("GET", "/", nil),
			},
			want: struct {
				headers    http.Header
				statusCode int
			}{
				statusCode: http.StatusOK,
				headers: http.Header{
					"Access-Control-Allow-Origin":  []string{"http://localhost:3000"},
					"Access-Control-Allow-Methods": []string{"GET, POST, DELETE"},
					"Access-Control-Allow-Headers": []string{"Content-Type, Passphrase"},
				},
			},
		},
		{
			name: "CORS handler with OPTIONS method",
			input: struct {
				origin string
				req    *http.Request
			}{
				origin: "http://localhost:3000",
				req:    httptest.NewRequest("OPTIONS", "/", nil),
			},
			want: struct {
				headers    http.Header
				statusCode int
			}{
				statusCode: http.StatusOK,
				headers: http.Header{
					"Access-Control-Allow-Origin":  []string{"http://localhost:3000"},
					"Access-Control-Allow-Methods": []string{"GET, POST, DELETE"},
					"Access-Control-Allow-Headers": []string{"Content-Type, Passphrase"},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			rr := httptest.NewRecorder()
			req := test.input.req

			CORS(test.input.origin)(handler).ServeHTTP(rr, req)

			gotCode := rr.Code
			gotHeaders := rr.Header()

			if diff := cmp.Diff(test.want.statusCode, gotCode); diff != "" {
				t.Errorf("newCORSHandler() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.want.headers, gotHeaders); diff != "" {
				t.Errorf("newCORSHandler() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
