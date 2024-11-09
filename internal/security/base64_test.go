package security

import (
	"encoding/base64"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestDecodeBase64SHA256(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			sha256   []byte
			encoding *base64.Encoding
		}
		want    []byte
		wantErr error
	}{
		{
			name: "decode base64 standard encoding SHA-256",
			input: struct {
				sha256   []byte
				encoding *base64.Encoding
			}{
				sha256:   SHA256([]byte("key")),
				encoding: base64.StdEncoding,
			},
			want: SHA256([]byte("key")),
		},
		{
			name: "decode base64 raw standard encoding SHA-256",
			input: struct {
				sha256   []byte
				encoding *base64.Encoding
			}{
				sha256:   SHA256([]byte("key")),
				encoding: base64.RawStdEncoding,
			},
			want: SHA256([]byte("key")),
		},
		{
			name: "decode base64 URL encoding SHA-256",
			input: struct {
				sha256   []byte
				encoding *base64.Encoding
			}{
				sha256:   SHA256([]byte("key")),
				encoding: base64.URLEncoding,
			},
			want: SHA256([]byte("key")),
		},
		{
			name: "decode base64 raw URL encoding SHA-256",
			input: struct {
				sha256   []byte
				encoding *base64.Encoding
			}{
				sha256:   SHA256([]byte("key")),
				encoding: base64.RawURLEncoding,
			},
			want: SHA256([]byte("key")),
		},

		{
			name: "decode base64 - error",
			input: struct {
				sha256   []byte
				encoding *base64.Encoding
			}{
				sha256:   []byte("key"),
				encoding: base64.RawURLEncoding,
			},
			wantErr: ErrInvalidHashLength,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hash := test.input.encoding.EncodeToString(test.input.sha256)

			got, gotErr := DecodeBase64SHA256(hash)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("decodeBase64SHA256() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("decodeBase64SHA256() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}
