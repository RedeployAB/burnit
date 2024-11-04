package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolveIP(t *testing.T) {
	var tests = []struct {
		name  string
		input func() *http.Request
		want  string
	}{
		{
			name: "With Forwarded header",
			input: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("Forwarded", "for=192.168.1.1:1234")
				return req
			},
			want: "192.168.1.1",
		},
		{
			name: "With X-Forwarded-For header",
			input: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("X-Forwarded-For", "192.168.1.1:1234")
				return req
			},
			want: "192.168.1.1",
		},
		{
			name: "With X-Real-IP header",
			input: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("X-Real-IP", "192.168.1.1:1234")
				return req
			},
			want: "192.168.1.1",
		},
		{
			name: "With RemoteAddr",
			input: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.RemoteAddr = "192.168.1.1:1234"
				return req
			},
			want: "192.168.1.1",
		},
		{
			name: "With invalid RemoteAddr",
			input: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.RemoteAddr = "1234"
				return req
			},
			want: "N/A",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := resolveIP(test.input())
			if test.want != got {
				t.Errorf("Resolve() = unexpected result, want %s, got: %s", test.want, got)
			}
		})
	}
}
