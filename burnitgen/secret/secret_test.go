package secret

import (
	"testing"
	"unicode/utf8"
)

func TestGenerate(t *testing.T) {
	var tests = []struct {
		x int
		y bool
		n int
	}{
		{8, false, 8},
		{16, true, 16},
	}

	for _, test := range tests {
		secret := Generate(test.x, test.y)
		count := utf8.RuneCountInString(secret)
		if count != test.n {
			t.Errorf("number of characters in generated string incorrect, got: %d, want: %d", count, test.n)
		}
	}
}
