package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	newUUID = func() string {
		return "test"
	}

	var tests = []struct {
		name  string
		input *http.Request
		want  string
	}{
		{
			name: "With request ID",
			input: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				return req
			}(),
			want: "test",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got := getRequestID(r.Context())
				if test.want != got {
					t.Errorf("RequestID() = unexpected result, want %s, got: %s", test.want, got)
				}
			})

			rr := httptest.NewRecorder()
			req := test.input
			RequestID()(handler).ServeHTTP(rr, req)
		})
	}
}
