package internal

import (
	"math/rand"
	"strings"
	"time"
)

var stdLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUWXYZ_-"

// GenerateRandomString generates a random string. Argument
// l determines amount of characters in the
// resulting string. Argument sc determines if
// special characters should be used.
func GenerateRandomString(l int, sc bool) string {
	var strbld strings.Builder

	strbld.WriteString(stdLetters)
	if sc == true {
		strbld.WriteString("!?=()&%")
	}
	bltrs := []byte(strbld.String())

	rand.Seed(time.Now().UnixNano())
	b := make([]byte, l)
	for i := range b {
		b[i] = bltrs[rand.Intn(len(bltrs))]
	}

	return string(b)
}
