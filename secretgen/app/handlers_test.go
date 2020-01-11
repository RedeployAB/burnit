package app

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"unicode/utf8"

	"github.com/gorilla/mux"
)

func SetupServer() Server {
	srv := Server{
		router: mux.NewRouter(),
	}
	srv.routes()
	return srv
}

func TestGenerateSecret(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/generate", nil)
	res := httptest.NewRecorder()
	SetupServer().router.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("Status code incorrect, got: %d, want: 200", res.Code)
	}

	var rb secretResponse
	b, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(b, &rb); err != nil {
		t.Errorf("Unmarshal failed.")
	}

	if rb.Data.Secret == "" {
		t.Errorf("Response incorrect, got: empty string, want: %s", rb.Data.Secret)
	}

	runeCount := utf8.RuneCountInString(rb.Data.Secret)
	if runeCount != 16 {
		t.Errorf("Response secret length incorrect, got: %d, want: 16", runeCount)
	}
}

func TestGenerateSecretHandlerParams(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v0/generate?length=22&specialshars=true", nil)
	res := httptest.NewRecorder()
	SetupServer().router.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("Status code incorrect, got: %d, want: 200", res.Code)
	}

	var rb secretResponse
	b, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(b, &rb); err != nil {
		t.Errorf("Unmarshal failed.")
	}

	if rb.Data.Secret == "" {
		t.Errorf("Response incorrect, got: empty string, want: %s", rb.Data.Secret)
	}

	runeCount := utf8.RuneCountInString(rb.Data.Secret)
	if runeCount != 22 {
		t.Errorf("Response secret length incorrect, got: %d, want: 22", runeCount)
	}
}

func TestNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()
	SetupServer().router.ServeHTTP(res, req)

	if res.Code != 404 {
		t.Errorf("Status code incorrect, got: %d, want: 404", res.Code)
	}
}

func TestHandleGenerateSecretQuery(t *testing.T) {
	query1 := url.Values{}
	query2 := url.Values{}
	query2.Set("length", "22")
	query2.Set("specialchars", "true")

	var tests = []struct {
		in           url.Values
		length       int
		specialChars bool
	}{
		{query1, 16, false},
		{query2, 22, true},
	}

	for _, test := range tests {
		l, sc := parseGenerateSecretQuery(test.in)
		if l != test.length {
			t.Errorf("got: %v, want: %v", l, test.length)
		}
		if sc != test.specialChars {
			t.Errorf("got %v, want: %v", sc, test.specialChars)
		}
	}
}
