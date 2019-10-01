package api

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

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/generate", generateSecret).Methods("GET")
	return router
}

func TestGenerateSecretHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/generate", nil)
	res := httptest.NewRecorder()
	Router().ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("Status code incorrect, got: %d, want: 200", res.Code)
	}

	var rb secretResponse
	b, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(b, &rb)
	if err != nil {
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
	req, _ := http.NewRequest("GET", "/generate?length=22&specialshars=true", nil)
	res := httptest.NewRecorder()
	Router().ServeHTTP(res, req)

	if res.Code != 200 {
		t.Errorf("Status code incorrect, got: %d, want: 200", res.Code)
	}

	var rb secretResponse
	b, err := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(b, &rb)
	if err != nil {
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
