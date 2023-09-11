package secure

import (
	"crypto/sha512"
	"encoding/hex"
)

func NewSha512(bytes []byte) (string, error) {
	h := sha512.New()
	_, err := h.Write(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
