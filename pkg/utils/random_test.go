package utils

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

const size = 32

func TestRandomBytes(t *testing.T) {
	b, err := RandomBytes(size)
	assert.Nil(t, err)
	assert.Equal(t, size, len(b))
}

func TestRandomBytesToBase64(t *testing.T) {
	s, err := RandomBytesToBase64(size)
	assert.Nil(t, err)
	decoded, err := base64.URLEncoding.DecodeString(s)
	assert.Nil(t, err)
	assert.Equal(t, size, len(decoded))
}

func TestRandomBytesToHex(t *testing.T) {
	s, err := RandomBytesToHex(size)
	assert.Nil(t, err)
	decoded, err := hex.DecodeString(s)
	assert.Nil(t, err)
	assert.Equal(t, size, len(decoded))
}

func TestRandomInt(t *testing.T) {
	minN, maxN := int64(1), int64(10)
	n, err := RandomInt(minN, maxN)
	assert.Nil(t, err)
	assert.True(t, n >= minN && n <= maxN)
}

func TestRandomString(t *testing.T) {
	n := 10
	s, err := RandomString(n)
	assert.Nil(t, err)
	assert.Equal(t, n, len(s))
}

func TestRandomStringWithEqualAlphabet(t *testing.T) {
	n := 10
	s, err := RandomString(n, []byte("a")...)
	assert.Nil(t, err)
	assert.Equal(t, "aaaaaaaaaa", s)
}
