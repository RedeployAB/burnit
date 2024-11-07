package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCacheControl(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req     *http.Request
			options []CacheControlOption
		}
		want http.Header
	}{
		{
			name: "no store",
			input: struct {
				req     *http.Request
				options []CacheControlOption
			}{
				req: httptest.NewRequest("GET", "/", nil),
				options: []CacheControlOption{
					func(o *CacheControlOptions) {
						o.NoStore = true
					},
				},
			},
			want: http.Header{
				"Cache-Control": []string{"no-store"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			rr := httptest.NewRecorder()
			req := test.input.req

			CacheControl(test.input.options...)(handler).ServeHTTP(rr, req)

			got := rr.Header()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("CacheControl() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}
