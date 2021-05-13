package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSHandler(t *testing.T) {
	corsHandler := CORSHandler{
		Origin: "http://localhost",
		Headers: map[string][]string{
			"Access-Control-Allow-Headers": {"Location, Authorization"},
			"Access-Control-Allow-Methods": {"GET", "POST"},
		},
	}

	h := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}

	req, _ := http.NewRequest(http.MethodOptions, "/", nil)
	res := httptest.NewRecorder()

	corsHandler.Handle(h()).ServeHTTP(res, req)
	if res.Header().Get("Access-Control-Allow-Origin") != "http://localhost" {
		t.Errorf("incorrect value, got: %s, want: %s", res.Header().Get("Access-Control-Allow-Origin"), "http://localhost")
	}
	if res.Header().Get("Access-Control-Allow-Headers") != "Location, Authorization" {
		t.Errorf("incorrect value, got: %s, want: %s", res.Header().Get("Access-Control-Allow-Headers"), "Location, Authorization")
	}
	if res.Header().Get("Access-Control-Allow-Methods") != "GET, POST" {
		t.Errorf("incorrect value, got: %s, want: %s", res.Header().Get("Access-Control-Allow-Methods"), "GET, POST")
	}

	req, _ = http.NewRequest(http.MethodGet, "/", nil)
	res = httptest.NewRecorder()

	corsHandler.Handle(h()).ServeHTTP(res, req)
}
