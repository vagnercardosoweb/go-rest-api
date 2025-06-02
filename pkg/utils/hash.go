package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

func HashSHA512(bytes []byte) string {
	hash := sha512.Sum512(bytes)
	return hex.EncodeToString(hash[:])
}

func HashSHA256(bytes []byte) string {
	hash := sha256.Sum256(bytes)
	return hex.EncodeToString(hash[:])
}
