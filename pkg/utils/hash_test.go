package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashWithSha512(t *testing.T) {
	const expected = "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff"
	got := HashSHA512([]byte("test"))
	assert.Equal(t, expected, got)
}

func TestHashWithSha256(t *testing.T) {
	const expected = "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"
	got := HashSHA256([]byte("test"))
	assert.Equal(t, expected, got)
}
