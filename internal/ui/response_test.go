package ui

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFormatErrorMessage(t *testing.T) {
	var tests = []struct {
		name  string
		input error
		want  string
	}{
		{
			name:  "nil error",
			input: nil,
			want:  "",
		},
		{
			name:  "error lowercase",
			input: errors.New("an error occurred"),
			want:  "An error occurred.",
		},
		{
			name:  "error lowercase with period",
			input: errors.New("an error occurred."),
			want:  "An error occurred.",
		},
		{
			name:  "error with prefix and period",
			input: errors.New("failed to get secret: an error occurred."),
			want:  "An error occurred.",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := formatErrorMessage(test.input)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("formatErrorMessage(%q) = unexpected result (-want +got)\n%s\n", test.input, diff)
			}
		})
	}
}
