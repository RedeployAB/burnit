package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCompress(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			req *http.Request
		}
		want struct {
			headers http.Header
			body    []byte
		}
	}{
		{
			name: "Do not compress response",
			input: struct{ req *http.Request }{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/", nil)
					return req
				}(),
			},
			want: struct {
				headers http.Header
				body    []byte
			}{
				headers: http.Header{
					"Content-Type": []string{"text/html"},
				},
				body: []byte("test"),
			},
		},
		{
			name: "Compress response",
			input: struct{ req *http.Request }{
				req: func() *http.Request {
					req, _ := http.NewRequest("GET", "/", nil)
					req.Header.Set("Accept-Encoding", "gzip")
					return req
				}(),
			},
			want: struct {
				headers http.Header
				body    []byte
			}{
				headers: http.Header{
					"Content-Encoding": []string{"gzip"},
					"Content-Type":     []string{"text/html"},
				},
				body: []byte("test"),
			},
		},
	}

	mux := http.NewServeMux()
	mux.Handle("/", Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("test"))
	})))

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte("test"))
			})

			rr := httptest.NewRecorder()
			req := test.input.req

			Compress()(handler).ServeHTTP(rr, req)

			if diff := cmp.Diff(test.want.headers, rr.Header()); diff != "" {
				t.Errorf("Compress() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.want.body, testCompressReadBody(t, rr)); diff != "" {
				t.Errorf("Compress() = unexpected result (-want +got)\n%s\n", diff)
			}
		})
	}
}

func testCompressReadBody(t *testing.T, rr *httptest.ResponseRecorder) []byte {
	var reader io.ReadCloser
	header := rr.Header().Get("Content-Encoding")
	if len(header) == 0 || header != "gzip" {
		reader = io.NopCloser(bytes.NewReader(rr.Body.Bytes()))
	} else {
		var err error
		reader, err = gzip.NewReader(bytes.NewReader(rr.Body.Bytes()))
		if err != nil {
			t.Fatalf("Failed to create gzip reader: %v", err)
			return nil
		}
	}

	b, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
		return nil
	}
	return b
}
