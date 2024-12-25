package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLogger(t *testing.T) {
	newUUID = func() string {
		return "test"
	}

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
			want: []string{"Request received.", "type", "request", "status", "200", "path", "/", "method", "GET", "requestId", "test", "sourceIp", "192.168.1.1"},
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
			want: []string{"Request received.", "type", "request", "status", "200", "path", "/", "method", "GET", "requestId", "test", "sourceIp", "192.168.1.1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logs := []string{}
			log := &stubLogger{
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
			Logger(log)(handler).ServeHTTP(rr, req)

			if diff := cmp.Diff(test.want, logs); diff != "" {
				t.Errorf("requestLogger() = unexpected result, (-want, +got):\n%s\n", diff)
			}
		})
	}
}

type stubLogger struct {
	logs *[]string
}

func (l *stubLogger) Info(msg string, args ...any) {
	if l.logs == nil {
		l.logs = &[]string{}
	}

	messages := []string{msg}
	for _, v := range args {
		var val string
		switch v := v.(type) {
		case string:
			val = v
		case int:
			val = strconv.Itoa(v)
		}
		messages = append(messages, val)
	}
	*l.logs = append(*l.logs, messages...)
}

func (l *stubLogger) Error(msg string, args ...any) {
	if l.logs == nil {
		l.logs = &[]string{}
	}

	messages := []string{msg}
	for _, v := range args {
		var val string
		switch v := v.(type) {
		case string:
			val = v
		case int:
			val = strconv.Itoa(v)
		}
		messages = append(messages, val)
	}
	*l.logs = append(*l.logs, messages...)
}

func (l *stubLogger) Debug(msg string, args ...any) {}

func (l *stubLogger) Warn(msg string, args ...any) {}
