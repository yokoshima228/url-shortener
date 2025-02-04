package random

import (
	"math/rand/v2"
	"strings"
)

func NewRandomString(l int) string {
	var result strings.Builder
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	for i := 0; i < l; i++ {
		symbol := string(chars[rand.IntN(len(chars))])
		result.Write([]byte(symbol))
	}

	return result.String()
}
