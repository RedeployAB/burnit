package api

import (
	"net/url"
	"testing"
)

func TestHandleGenerateSecretQuery(t *testing.T) {
	query1 := url.Values{}
	query2 := url.Values{}
	query2.Set("length", "22")
	query2.Set("specialchars", "true")

	var tests = []struct {
		in  url.Values
		out secretParams
	}{
		{query1, secretParams{Length: 16, SpecialCharacters: false}},
		{query2, secretParams{Length: 22, SpecialCharacters: true}},
	}

	for _, test := range tests {
		params := handleGenerateSecretQuery(test.in)
		if params != test.out {
			t.Errorf("got: %v, want: %v", params, test.out)
		}
	}
}
