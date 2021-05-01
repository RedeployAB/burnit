package httperror

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	var tests = []struct {
		statusCode    int
		code          string
		message       string
		wantedCode    int
		wantedMessage string
	}{
		{statusCode: http.StatusBadRequest, code: "", message: "", wantedCode: 400, wantedMessage: `{"message":"malformed JSON","statusCode":400}`},
		{statusCode: http.StatusUnauthorized, code: "", message: "", wantedCode: 401, wantedMessage: `{"message":"unauthorized","statusCode":401}`},
		{statusCode: http.StatusForbidden, code: "", message: "", wantedCode: 403, wantedMessage: `{"message":"forbidden","statusCode":403}`},
		{statusCode: http.StatusNotFound, code: "", message: "", wantedCode: 404, wantedMessage: `{"message":"not found","statusCode":404}`},
		{statusCode: http.StatusInternalServerError, code: "", message: "", wantedCode: 500, wantedMessage: `{"message":"internal server error","statusCode":500}`},
		{statusCode: http.StatusNotImplemented, code: "", message: "", wantedCode: 501, wantedMessage: `{"message":"not implemented","statusCode":501}`},
		{statusCode: http.StatusBadRequest, code: "BadInput", message: "Bad Input", wantedCode: 400, wantedMessage: `{"message":"Bad Input","code":"BadInput","statusCode":400}`},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		Write(w, test.statusCode, test.code, test.message)

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
