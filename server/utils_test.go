package server

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestEncode(t *testing.T) {
	type testResponse struct {
		Name string `json:"name"`
	}

	var tests = []struct {
		name  string
		input struct {
			status int
			v      testResponse
		}
		want struct {
			status      int
			body        string
			contentType string
		}
		wantErr error
	}{
		{
			name: "encode",
			input: struct {
				status int
				v      testResponse
			}{
				status: http.StatusOK,
				v:      testResponse{Name: "test"},
			},
			want: struct {
				status      int
				body        string
				contentType string
			}{
				status:      http.StatusOK,
				body:        "{\"name\":\"test\"}\n",
				contentType: contentTypeJSON,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			gotErr := encode(rr, test.input.status, test.input.v)
			gotStatus := rr.Code
			gotBody := rr.Body.String()
			contentType := rr.Header().Get("Content-Type")

			if test.want.status != gotStatus {
				t.Errorf("encode() = unexpected status, want: %d, got: %d\n", test.want.status, gotStatus)
			}

			if test.want.body != gotBody {
				t.Errorf("encode() = unexpected body, want: %s, got: %s\n", test.want.body, gotBody)
			}

			if test.want.contentType != contentType {
				t.Errorf("encode() = unexpected content type, want: %s, got: %s\n", test.want.contentType, contentType)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("encode() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	var tests = []struct {
		name    string
		input   *http.Request
		want    testRequest
		wantErr error
	}{
		{
			name:  "decode",
			input: httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{\"name\":\"test\"}"))),
			want:  testRequest{Name: "test"},
		},
		{
			name:    "decode invalid request",
			input:   httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{}"))),
			wantErr: ErrInvalidRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := decode[testRequest](test.input)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("decode() = unexpected value (-want +got)\n%s\n", diff)
			}

			if diff := cmp.Diff(test.wantErr, gotErr, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("decode() = unexpected error (-want +got)\n%s\n", diff)
			}
		})
	}
}

type testRequest struct {
	Name string `json:"name"`
}

func (t testRequest) Valid(ctx context.Context) map[string]string {
	if len(t.Name) == 0 {
		return map[string]string{"name": "name is required"}
	}
	return nil
}
