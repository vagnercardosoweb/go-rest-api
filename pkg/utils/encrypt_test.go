package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEncryptKey(t *testing.T) {
	bytes, err := NewEncryptKey()
	assert.Nil(t, err)
	assert.Equal(t, 32, len(bytes[:]))
}

func TestEncryptDecrypt(t *testing.T) {
	key, err := NewEncryptKey()
	assert.Nil(t, err)
	plaintext := []byte("Hello World")
	ciphertext, err := Encrypt(plaintext, key)
	assert.Nil(t, err)
	decrypted, err := Decrypt(ciphertext, key)
	assert.Nil(t, err)
	assert.Equal(t, plaintext, decrypted)
}

func TestEncryptDecryptBadKey(t *testing.T) {
	key, err := NewEncryptKey()
	assert.Nil(t, err)
	plaintext := []byte("Hello World")
	ciphertext, err := Encrypt(plaintext, key)
	assert.Nil(t, err)
	key, err = NewEncryptKey()
	assert.Nil(t, err)
	_, err = Decrypt(ciphertext, key)
	assert.NotNil(t, err)
}

func TestEncryptDecryptBadCiphertext(t *testing.T) {
	key, err := NewEncryptKey()
	assert.Nil(t, err)
	plaintext := []byte("Hello World")
	ciphertext, err := Encrypt(plaintext, key)
	assert.Nil(t, err)
	ciphertext[0] = ciphertext[0] + 1
	_, err = Decrypt(ciphertext, key)
	assert.NotNil(t, err)
}
