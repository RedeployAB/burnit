package security

import (
	"encoding/base64"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestDecodeBase64(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			s        []byte
			encoding *base64.Encoding
		}
		want    string
		wantErr error
	}{
		{
			name: "decode base64 standard encoding",
			input: struct {
				s        []byte
				encoding *base64.Encoding
			}{
				s:        []byte("key"),
				encoding: base64.StdEncoding,
			},
			want: "key",
		},
		{
			name: "decode base64 raw standard encoding",
			input: struct {
				s        []byte
				encoding *base64.Encoding
			}{
				s:        []byte("key"),
				encoding: base64.RawStdEncoding,
			},
			want: "key",
		},
		{
			name: "decode base64 URL encoding",
			input: struct {
				s        []byte
				encoding *base64.Encoding
			}{
				s:        []byte("key"),
				encoding: base64.URLEncoding,
			},
			want: "key",
		},
		{
			name: "decode base64 raw URL encoding SHA-256",
			input: struct {
				s        []byte
				encoding *base64.Encoding
			}{
				s:        []byte("key"),
				encoding: base64.RawURLEncoding,
			},
			want: "key",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash := test.input.encoding.EncodeToString(test.input.s)

			got, gotErr := DecodeBase64(hash)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("decodeBase64SHA256() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("decodeBase64SHA256() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}
