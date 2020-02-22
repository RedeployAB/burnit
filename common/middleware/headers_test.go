package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStrip(t *testing.T) {
	exceptions := []string{"X-API-Key", "X-Passphrase"}
	hmw := HeaderStrip{Exceptions: exceptions}

	h := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	token := "abcdefg"
	passphrase := "gfedcba"
	inHeaders := map[string]string{
		"Forwarded":    "for=192.168.0.1",
		"X-Real-IP":    "192.168.0.1",
		"X-API-Key":    token,
		"X-Passphrase": passphrase,
	}

	for k, v := range inHeaders {
		req.Header.Set(k, v)
	}

	hmw.Strip(h()).ServeHTTP(res, req)
	if len(req.Header.Get("Forwarded")) > 0 {
		t.Errorf("incorrect value, got: %s, want: %s", req.Header.Get("Forwarded"), "")
	}
	if len(req.Header.Get("X-Real-IP")) > 0 {
		t.Errorf("incorrect value, got: %s, want: %s", req.Header.Get("X-Real-IP"), "")
	}
	if req.Header.Get("X-API-Key") != token {
		t.Errorf("incorrect value, got: %s, want: %s", req.Header.Get("X-API-Key"), token)
	}
	if req.Header.Get("X-Passphrase") != passphrase {
		t.Errorf("incorrect value, got: %s, want: %s", req.Header.Get("X-Passphrase"), passphrase)
	}
}
