package util

import (
	"math/rand"
	"strings"
)

func GetRandomStr() string {
	bytes := make([]byte, 6)
	for i := range bytes {
		bytes[i] = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"[rand.Intn(62)]
	}
	return string(bytes)
}

func ReverseStr(s string) string {
	builder := strings.Builder{}
	for i := len(s) - 1; i >= 0; i-- {
		builder.WriteByte(s[i])
	}
	return builder.String()
}
