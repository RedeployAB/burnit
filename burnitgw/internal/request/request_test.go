package request

import (
	"encoding/json"
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
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mr := ResponseBody{
			Secret: mockResponse{
				Value: "secret",
			},
		}
		json.NewEncoder(w).Encode(&mr)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, "/path")
	opts := Options{Method: GET}
	res, err := client.Request(opts)
	if err != nil {
		t.Fatalf("Error in setting client call: %v", err)
	}

	// Update this to use type assertion later on.
	jsonRes, err := json.Marshal(res)
	if err != nil {
		t.Fatalf("Error in JSON marshaling: %v", err)
	}

	expected := `{"secret":{"value":"secret"}}`
	if string(jsonRes) != expected {
		t.Errorf(`Incorrect value, got: %s, want: %s`, string(jsonRes), expected)
	}
}

func TestBasicRequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(400)
		return
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, "/path")
	opts := Options{Method: GET}
	_, err := client.Request(opts)
	if err == nil {
		t.Fatal("Incorrect, should result in an error")
	}

	if err.(*Error).code != 400 {
		t.Errorf("Incorrect value, got: %v, want: 400", err.(*Error).code)
	}
	if err.(*Error).Error() != "code 400: 400 Bad Request" {
		t.Errorf("Incorrect value, got: %s, want: code 400: Bad Request", err.(*Error).Error())
	}
}

func TestPostRequestWithoutBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mr := ResponseBody{}
		json.NewEncoder(w).Encode(&mr)
	}))

	client := NewHTTPClient(srv.URL, "/path")
	opts := Options{Method: POST}
	_, err := client.Request(opts)

	if err == nil {
		t.Fatal("Incorrect, should result in an error")
	}

	if err.(*Error).code != 400 {
		t.Errorf("Incorrect value, got: %v, want: 400", err.(*Error).code)
	}
	if err.(*Error).Error() != "code 400: 400 Bad Request" {
		t.Errorf("Incorrect value, got: %s, want: code 400: Bad Request", err.(*Error).Error())
	}
}

func TestRequestWithParams(t *testing.T) {
	expectedParam := "1234"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.String(), "/1234") {
			w.WriteHeader(400)
			return
		}
		mr := ResponseBody{
			Secret: mockResponse{
				Value: "secret",
			},
		}
		json.NewEncoder(w).Encode(&mr)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, "/path")
	params := map[string]string{"id": expectedParam}
	opts := Options{Method: GET, Params: params}
	_, err := client.Request(opts)
	if err != nil {
		t.Errorf("Incorrect value, got: %v, want: <nil>", err)
	}
}

func TestBasicRequestWithQuery(t *testing.T) {
	expectedQuery := "length=10"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.String(), "?"+expectedQuery) {
			w.WriteHeader(400)
			return
		}
		mr := ResponseBody{
			Secret: mockResponse{
				Value: "secret",
			},
		}
		json.NewEncoder(w).Encode(&mr)
	}))
	defer srv.Close()

	client := NewHTTPClient(srv.URL, "/path")
	opts := Options{Method: GET, Query: expectedQuery}
	_, err := client.Request(opts)
	if err != nil {
		t.Errorf("Incorrect value, got: %v, want: <nil>", err)
	}
}

func TestBasicRequestWithHeaders(t *testing.T) {
	hdrName := "X-PASSPHRASE"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hdrVal := r.Header.Get(hdrName)
		if len(hdrVal) == 0 {
			t.Error("Incorrect, header should be set.")
			return
		}
		mr := ResponseBody{}
		json.NewEncoder(w).Encode(&mr)
	}))

	client := NewHTTPClient(srv.URL, "/path")

	hdr := http.Header{}
	hdr.Add(hdrName, "test")
	opts := Options{Method: GET, Header: hdr}
	_, _ = client.Request(opts)
}
