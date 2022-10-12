package secret

import (
	"math/rand"
	"strings"
	"time"
)

var stdLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUWXYZ"
var generatedMaxLength = 512

// generate generates a random string. Argument
// l determines amount of characters in the
// resulting string. Argument sc determines if
// special characters should be used.
func generate(l int, sc bool) string {
	if l > generatedMaxLength {
		l = generatedMaxLength
	}

	var strb strings.Builder
	strb.WriteString(stdLetters)
	if sc {
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
