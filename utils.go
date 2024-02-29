package caddydiscord

import (
	"crypto/rand"
	"fmt"
)

func randomness(length uint) []byte {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return randomBytes
}

type CookieNamer func(string) string

func CookieName(executionKey string) CookieNamer {
	return func(realm string) string {
		return fmt.Sprintf("%s_%s_%s", cookieName, realm, executionKey)
	}
}
