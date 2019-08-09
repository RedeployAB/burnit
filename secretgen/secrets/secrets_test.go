package secrets

import (
	"testing"
	"unicode/utf8"
)

func TestGenerateSecret(t *testing.T) {
	var tests = []struct {
		x int
		y bool
		n int
	}{
		{8, false, 8},
		{16, true, 16},
	}

	for _, test := range tests {
		secret := GenerateSecret(test.x, test.y)
		count := utf8.RuneCountInString(secret)
		if count != test.n {
			t.Errorf("Number of characters in generated string incorrect, got: %d, want: %d", count, test.n)
		}
	}
}

func BenchmarkGenerateSecret(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateSecret(8, true)
	}
}
