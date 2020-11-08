package db

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/RedeployAB/burnit/burnitgw/services/request"
)

type mockClient struct {
	address string
	path    string
	mode    string
}

func (c mockClient) Request(o request.Options) ([]byte, error) {
	var err error
	var responseJSON []byte
	switch c.mode {
	case "db-get-success":
		responseJSON = []byte(`{"id":"1234","value":"secret"}`)
	case "db-get-fail":
		err = errors.New("call to api failed")
	case "db-get-malformed":
		responseJSON = []byte(`{"value":`)
	case "db-create-success":
		responseJSON = []byte(`{"id":"4321","value":"terces"}`)
	case "db-create-fail":
		err = errors.New("call to api failed")
	case "db-create-mailformed":
		responseJSON = []byte(`{"value":`)
	}

	return responseJSON, err
}

func TestGet(t *testing.T) {
	secret1 := &Secret{}
	secret1.ID = "1234"
	secret1.Value = "secret"

	secret2 := &Secret{}
	secret2.ID = "4321"
	secret2.Value = "terces"

	var tests = []struct {
		mode    string
		param   string
		want    *Secret
		wantErr error
	}{
		{mode: "db-get-success", param: "1234", want: secret1, wantErr: nil},
		{mode: "db-get-fail", param: "4321", want: nil, wantErr: errors.New("call to api failed")},
		{mode: "db-get-malformed", param: "1234", want: nil, wantErr: errors.New("unexpected end of JSON input")},
	}

	for _, test := range tests {
		u, _ := url.Parse("http://localhost:3001/secrets/" + test.param)
		r := &http.Request{URL: u}
		svc := NewService(&mockClient{mode: test.mode})
		got, err := svc.Get(r, map[string]string{"id": test.param})

		if got != nil && got.ID != test.want.ID {
			t.Errorf("incorrect value, got: %s, want: %s", got.ID, test.want.ID)
		}
		if got != nil && got.Value != test.want.Value {
			t.Errorf("incorrect value, got: %s, want: %s", got.Value, test.want.Value)
		}
		if got == nil && err != nil && err.Error() != test.wantErr.Error() {
			t.Errorf("incorrect value, got: %s, want: %s", err.Error(), test.wantErr.Error())
		}
	}
}

func TestCreate(t *testing.T) {
	secret1 := &Secret{}
	secret1.ID = "4321"
	secret1.Value = "terces"
	jsonStr := []byte(`{"value":"terces"}`)

	var tests = []struct {
		mode    string
		body    []byte
		want    *Secret
		wantErr error
	}{
		{mode: "db-create-success", body: jsonStr, want: secret1, wantErr: nil},
		{mode: "db-create-fail", body: jsonStr, want: nil, wantErr: errors.New("call to api failed")},
		{mode: "db-create-malformed", body: jsonStr, want: nil, wantErr: errors.New("unexpected end of JSON input")},
	}

	for _, test := range tests {
		u, _ := url.Parse("http://localhost:3001/secrets")
		r := &http.Request{URL: u, Body: ioutil.NopCloser(bytes.NewBuffer(test.body))}
		svc := NewService(&mockClient{mode: test.mode})
		got, err := svc.Create(r)

		if got != nil && got.ID != test.want.ID {
			t.Errorf("incorrect value, got: %s, want: %s", got.ID, test.want.ID)
		}
		if got != nil && got.Value != test.want.Value {
			t.Errorf("incorrect value, got: %s, want: %s", got.Value, test.want.Value)
		}
		if got == nil && err != nil && err.Error() != test.wantErr.Error() {
			t.Errorf("incorrect value, got: %s, want: %s", err.Error(), test.wantErr.Error())
		}
	}
}
