package security

import (
	"encoding/base64"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestDecodeBase64(t *testing.T) {
	var tests = []struct {
		name    string
		input   string
		want    []byte
		wantErr error
	}{
		{
			name:  "decode base64 standard encoding",
			input: base64.StdEncoding.EncodeToString([]byte("key")),
			want:  []byte("key"),
		},
		{
			name:  "decode base64 raw standard encoding",
			input: base64.RawStdEncoding.EncodeToString([]byte("key")),
			want:  []byte("key"),
		},
		{
			name:  "decode base64 URL encoding",
			input: base64.URLEncoding.EncodeToString([]byte("key")),
			want:  []byte("key"),
		},
		{
			name:  "decode base64 raw URL encoding",
			input: base64.RawURLEncoding.EncodeToString([]byte("key")),
			want:  []byte("key"),
		},
		{
			name:    "invalid input",
			input:   "$invalid",
			wantErr: ErrInvalidBase64,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := DecodeBase64(test.input)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("DecodeBase64() = unexpected result (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("DecodeBase64() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}
