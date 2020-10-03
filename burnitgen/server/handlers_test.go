package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"unicode/utf8"
)

func SetupServer() Server {
	srv := Server{
		router: http.NewServeMux(),
	}
	srv.routes()
	return srv
}

func TestGenerateSecret(t *testing.T) {
	req, _ := http.NewRequest("GET", "/secret", nil)
	res := httptest.NewRecorder()
	SetupServer().router.ServeHTTP(res, req)

	expectedCode := 200
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}

	var rb secretResponse
	b, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(b, &rb); err != nil {
		t.Errorf("error in test: %v", err)
	}

	if len(rb.Secret.Value) == 0 {
		t.Errorf("response incorrect, got: empty string, want: %s", rb.Secret.Value)
	}

	c := utf8.RuneCountInString(rb.Secret.Value)
	expected := 16
	if c != expected {
		t.Errorf("response secret length incorrect, got: %d, want: %d", c, expected)
	}
}

func TestGenerateSecretHandlerParams(t *testing.T) {
	req, _ := http.NewRequest("GET", "/secret?length=22&specialshars=true", nil)
	res := httptest.NewRecorder()
	SetupServer().router.ServeHTTP(res, req)

	expectedCode := 200
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
	}

	var rb secretResponse
	b, err := ioutil.ReadAll(res.Body)

	if err = json.Unmarshal(b, &rb); err != nil {
		t.Errorf("error in test: %v", err)
	}

	if len(rb.Secret.Value) == 0 {
		t.Errorf("response incorrect, got: empty string, want: %s", rb.Secret.Value)
	}

	c := utf8.RuneCountInString(rb.Secret.Value)
	expected := 22
	if c != 22 {
		t.Errorf("response secret length incorrect, got: %d, want: %d", c, expected)
	}
}

func TestNotFound(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	res := httptest.NewRecorder()
	SetupServer().router.ServeHTTP(res, req)

	expectedCode := 404
	if res.Code != expectedCode {
		t.Errorf("status code incorrect, got: %d, want: %d", res.Code, expectedCode)
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
