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
		{input: http.StatusBadRequest, wantedCode: 400, wantedMessage: `{"error":"malformed JSON","code":400}`},
		{input: http.StatusUnauthorized, wantedCode: 401, wantedMessage: `{"error":"unauthorized","code":401}`},
		{input: http.StatusForbidden, wantedCode: 403, wantedMessage: `{"error":"forbidden","code":403}`},
		{input: http.StatusNotFound, wantedCode: 404, wantedMessage: `{"error":"not found","code":404}`},
		{input: http.StatusInternalServerError, wantedCode: 500, wantedMessage: `{"error":"internal server error","code":500}`},
		{input: http.StatusNotImplemented, wantedCode: 501, wantedMessage: `{"error":"not implemented","code":501}`},
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
