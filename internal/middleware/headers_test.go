package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHeaders(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			options []HeadersOption
		}
		want http.Header
	}{
		{
			name: "default",
			input: struct {
				req     *http.Request
				options []HeadersOption
			}{
				req:     httptest.NewRequest("GET", "/", nil),
				options: nil,
			},
			want: http.Header{
				"Cache-Control":             []string{"no-store"},
				"Strict-Transport-Security": []string{"max-age=31536000"},
				"X-Content-Type-Options":    []string{"nosniff"},
				"X-Frame-Options":           []string{"DENY"},
			},
		},
		{
			name: "with options",
			input: struct {
				req     *http.Request
				options []HeadersOption
			}{
				req: httptest.NewRequest("GET", "/", nil),
				options: []HeadersOption{
					func(o *HeadersOptions) {
						o.CacheControl = "no-store"
						o.ContentSecurityPolicy = "default-src 'self';"
					},
				},
			},
			want: http.Header{
				"X-Content-Type-Options":    []string{"nosniff"},
				"X-Frame-Options":           []string{"DENY"},
				"Cache-Control":             []string{"no-store"},
				"Content-Security-Policy":   []string{"default-src 'self';"},
				"Strict-Transport-Security": []string{"max-age=31536000"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			rr := httptest.NewRecorder()
			req := test.input.req

			Headers(test.input.options...)(handler).ServeHTTP(rr, req)

			got := rr.Header()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Headers() = unexpected result (-want +got)\n%s\n", diff)
			}

		})
	}
}
