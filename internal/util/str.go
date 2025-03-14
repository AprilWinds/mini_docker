package util

import (
	"math/rand"
)

func GetRandomStr() string {
	bytes := make([]byte, 6)
	for i := range bytes {
		bytes[i] = "0123456789abcdefghijklmnopqrstuvwxyz"[rand.Intn(36)]
	}
	return string(bytes)
}
