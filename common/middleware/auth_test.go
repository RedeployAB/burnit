package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedeployAB/burnit/common/auth"
)

func TestAuthenticate(t *testing.T) {
	token := "abcdefg"
	ts := auth.NewMemoryTokenStore()
	ts.Set(token, "userA")
	amw := Authentication{TokenStore: ts}

	h := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}

	var tests = []struct {
		input string
		want  int
	}{
		{input: token, want: 200},
		{input: "", want: 401},
		{input: "wrongkey", want: 403},
	}

	for _, test := range tests {
		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()
		req.Header.Set("API-Key", test.input)

		amw.Authenticate(h()).ServeHTTP(res, req)
		if res.Code != test.want {
			t.Errorf("incorrect value, got: %d, want: %d", res.Code, test.want)
		}
	}
}

func TestAddAuthHeader(t *testing.T) {
	token := "abcdefg"
	amw := AuthHeader{Token: token}

	h := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}

	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	amw.AddAuthHeader(h()).ServeHTTP(res, req)

	header := req.Header.Get("API-Key")
	if header != token {
		t.Errorf("incorrect value, got: %s, want: %s", header, token)
	}
}
