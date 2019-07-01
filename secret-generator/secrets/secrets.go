package secrets

import (
	"math/rand"
	"strings"
	"time"
)

var stdLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUWXYZ_-"

// Generates a secret.
func GenerateSecret(l int, sc bool) string {
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
