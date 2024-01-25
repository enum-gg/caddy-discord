package caddydiscord

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func hashString512(input string) string {
	hasher := sha512.New()
	hasher.Write([]byte(input))
	return hex.EncodeToString(hasher.Sum(nil))
}

func hashString256(input string, length int) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	fullHash := hex.EncodeToString(hasher.Sum(nil))

	if length > len(fullHash) {
		length = len(fullHash)
	}
	return fullHash[:length]
}
