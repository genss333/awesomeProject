package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func EnSha256Hash(data string) string {
	hashes := sha256.New()

	hashes.Write([]byte(data))

	hashBytes := hashes.Sum(nil)

	hashHex := hex.EncodeToString(hashBytes)

	return hashHex
}

func MatchSha256Hash(data string, hash string) bool {
	return EnSha256Hash(data) == hash
}
