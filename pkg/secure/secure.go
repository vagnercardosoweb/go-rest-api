package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// NewEncryptionKey generates a random 256-bit key. It will return an
// error if the system's secure random number generator fails to
// function correctly, in which case the caller should not continue.
// Taken from https://github.com/gtank/cryptopasta/blob/master/encrypt.go
func NewEncryptionKey() (*[32]byte, error) {
	key := [32]byte{}
	_, err := io.ReadFull(rand.Reader, key[:])
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// ParseEncryptionKey decodes the string representation of an encryption key
// and returns its bytes
func ParseEncryptionKey(s string) (*[32]byte, error) {
	// get hex encoded encryption key from cloud secret
	key, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, errors.New("encryption key byte length must be exactly 32 bytes")
	}
	// loop through each byte and add it to the 32 byte encryption key array (ek)
	ek := [32]byte{}
	copy(ek[:], key)
	return &ek, nil
}

// Encrypt encrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Output takes the
// form nonce|ciphertext|tag where '|' indicates concatenation.
// Taken from https://github.com/gtank/cryptopasta/blob/master/encrypt.go
func Encrypt(plaintext []byte, key *[32]byte) (ciphertext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
// Taken from https://github.com/gtank/cryptopasta/blob/master/encrypt.go
func Decrypt(cipherText []byte, key *[32]byte) (plaintext []byte, err error) {
	var block cipher.Block
	block, err = aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	var gcm cipher.AEAD
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("malformed cipherText")
	}

	plaintext, err = gcm.Open(nil,
		cipherText[:nonceSize],
		cipherText[nonceSize:],
		nil,
	)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
