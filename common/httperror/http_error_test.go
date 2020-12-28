package httperror

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestError(t *testing.T) {
	var tests = []struct {
		input         int
		wantedCode    int
		wantedMessage string
	}{
		{input: http.StatusBadRequest, wantedCode: 400, wantedMessage: `{"message":"malformed JSON","statusCode":400}`},
		{input: http.StatusUnauthorized, wantedCode: 401, wantedMessage: `{"message":"unauthorized","statusCode":401}`},
		{input: http.StatusForbidden, wantedCode: 403, wantedMessage: `{"message":"forbidden","statusCode":403}`},
		{input: http.StatusNotFound, wantedCode: 404, wantedMessage: `{"message":"not found","statusCode":404}`},
		{input: http.StatusInternalServerError, wantedCode: 500, wantedMessage: `{"message":"internal server error","statusCode":500}`},
		{input: http.StatusNotImplemented, wantedCode: 501, wantedMessage: `{"message":"not implemented","statusCode":501}`},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		Error(w, test.input)

		b, err := ioutil.ReadAll(w.Body)
		if err != nil {
			t.Errorf("error in test: %v", err)
		}

		if w.Code != test.wantedCode {
			t.Errorf("incorrect value, got: %d, want: %d", w.Code, test.wantedCode)
		}

		body := strings.TrimRight(strings.TrimRight(string(b), "\r\n"), "\n")
		if body != test.wantedMessage {
			t.Errorf("incorrect value, got: %s, want: %s", body, test.wantedMessage)
		}
	}
}
