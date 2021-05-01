package secret

import (
	"math/rand"
	"strings"
	"time"
)

var stdLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUWXYZ"
var maxLength = 512

// Generate generates a random string. Argument
// l determines amount of characters in the
// resulting string. Argument sc determines if
// special characters should be used.
func Generate(l int, sc bool) string {
	if l > maxLength {
		l = maxLength
	}

	var strb strings.Builder
	strb.WriteString(stdLetters)
	if sc == true {
		strb.WriteString("_-!?=()&%")
	}
	bltrs := []byte(strb.String())

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, l)
	for i := range b {
		b[i] = bltrs[rand.Intn(len(bltrs))]
	}

	return string(b)
}
