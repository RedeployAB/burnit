package request

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockResponse struct {
	Value string `json:"value"`
}

type mockFullResponse struct {
	Secret mockResponse
}

func TestBasicRequest(t *testing.T) {
	jsonByte := []byte(`{"secret":{"value":"secret"}}`)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(jsonByte)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "/path")
	opts := Options{Method: GET}
	res, err := client.Request(opts)
	if err != nil {
		t.Fatalf("Error in setting client call: %v", err)
	}

	expected := `{"secret":{"value":"secret"}}`
	if string(res.Body) != expected {
		t.Errorf(`Incorrect value, got: %s, want: %s`, string(res.Body), expected)
	}
}

func TestRequestWithParams(t *testing.T) {
	jsonByte := []byte(`{"secret":{"value":"secret"}}`)
	expectedParam := "1234"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.String(), "/1234") {
			w.WriteHeader(400)
			return
		}
		w.Write(jsonByte)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "/path")
	params := map[string]string{"id": expectedParam}
	opts := Options{Method: GET, Params: params}
	_, err := client.Request(opts)
	if err != nil {
		t.Errorf("Incorrect value, got: %v, want: <nil>", err)
	}
}

func TestBasicRequestWithQuery(t *testing.T) {
	jsonByte := []byte(`{"secret":{"value":"secret"}}`)
	expectedQuery := "length=10"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.String(), "?"+expectedQuery) {
			w.WriteHeader(400)
			return
		}
		w.Write(jsonByte)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, "/path")
	opts := Options{Method: GET, Query: expectedQuery}
	_, err := client.Request(opts)
	if err != nil {
		t.Errorf("Incorrect value, got: %v, want: <nil>", err)
	}
}

func TestBasicRequestWithHeaders(t *testing.T) {
	jsonByte := []byte(`{"secret":{"value":"secret"}}`)
	hdrName := "PASSPHRASE"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdrVal := r.Header.Get(hdrName)
		if len(hdrVal) == 0 {
			t.Error("Incorrect, header should be set.")
			return
		}
		w.Write(jsonByte)
	}))

	client := NewClient(srv.URL, "/path")

	hdr := http.Header{}
	hdr.Add(hdrName, "test")
	opts := Options{Method: GET, Header: hdr}
	_, _ = client.Request(opts)
}
