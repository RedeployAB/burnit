package generator

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/RedeployAB/burnit/burnitgw/service/request"
)

type mockClient struct {
	address string
	path    string
	mode    string
}

func (c mockClient) Request(o request.Options) (*request.Response, error) {
	var err error
	var code int
	var responseJSON []byte
	switch c.mode {
	case "gen-success":
		code = 200
		responseJSON = []byte(`{"value":"secret"}`)
	case "gen-fail":
		code = 500
		err = errors.New("call to api failed")
	case "gen-malformed":
		code = 400
		responseJSON = []byte(`{"value":`)
	}

	return &request.Response{StatusCode: code, Body: responseJSON}, err
}

func TestGenerate(t *testing.T) {
	secret1 := &Secret{}
	secret1.Value = "secret"

	var tests = []struct {
		mode    string
		want    *Secret
		wantErr error
	}{
		{mode: "gen-success", want: secret1, wantErr: nil},
		{mode: "gen-fail", want: nil, wantErr: errors.New("call to api failed")},
		{mode: "gen-malformed", want: nil, wantErr: errors.New("unexpected end of JSON input")},
	}

	for _, test := range tests {
		u, _ := url.Parse("http://localhost:3002/secret")
		r := &http.Request{URL: u}
		svc := NewService(&mockClient{mode: test.mode})
		got, err := svc.Generate(r)

		if got != nil && got.Value != test.want.Value {
			t.Errorf("incorrect value, got: %s, want: %s", got.Value, test.want.Value)
		}
		if got == nil && err != nil && err.Error() != test.wantErr.Error() {
			t.Errorf("incorrect value, got: %s, want: %s", err.Error(), test.wantErr.Error())
		}
	}
}
