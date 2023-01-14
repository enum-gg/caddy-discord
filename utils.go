package caddydiscord

import (
	"crypto/rand"
)

func randomness(length uint) []byte {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return randomBytes
}
