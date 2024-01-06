package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"math/big"
)

const alphaNum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// RandomBytes returns securely generated random bytes. It will return
// an error if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
// Taken from https://stackoverflow.com/questions/35781197/generating-a-random-fixed-length-byte-array-in-go
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	// Note that err == nil only if we read len(b) bytes.
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}

// RandomBytesToBase64 returns a URL-safe, base64 encoded, securely generated, random string.
// It will return an error if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue. This should be
// used when there are concerns about security and need something cryptographically
// secure.
func RandomBytesToBase64(n int) (string, error) {
	b, err := RandomBytes(n)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func RandomBytesToHex(n int) (string, error) {
	b, err := RandomBytes(n)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) (int64, error) {
	bg := big.NewInt(max - min + 1)
	n, err := rand.Int(rand.Reader, bg)
	return n.Int64() + min, err
}

// RandomString generates a random alphanumeric string of the specified length,
// optionally using only specified characters
func RandomString(n int, alphabets ...byte) (string, error) {
	chars := alphaNum
	if len(alphabets) > 0 {
		chars = string(alphabets)
	}

	cnt := len(chars)
	maxBytes := 255 / cnt * cnt

	bytes := make([]byte, n)

	randRead := n * 5 / 4
	randBytes := make([]byte, randRead)

	for i := 0; i < n; {
		if _, err := rand.Read(randBytes); err != nil {
			return "", err
		}

		for j := 0; i < n && j < randRead; j++ {
			b := int(randBytes[j])
			if b >= maxBytes {
				continue
			}

			bytes[i] = chars[b%cnt]
			i++
		}
	}

	return string(bytes), nil
}
