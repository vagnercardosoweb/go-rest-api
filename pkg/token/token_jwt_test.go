package token

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenJwt_Encode(t *testing.T) {
	jwt := NewJwt()
	token, err := jwt.Encode(Input{Subject: "any_sub"})
	assert.Nil(t, err)
	assert.NotNil(t, token)
	assert.Len(t, strings.Split(token, "."), 3)
}

func TestTokenJwt_Encode_WithoutSubject(t *testing.T) {
	jwt := NewJwt()
	_, err := jwt.Encode(Input{})
	assert.NotNil(t, err)
}

func TestTokenJwt_Decode(t *testing.T) {
	jwt := NewJwt()
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Hour * 2)

	token, err := jwt.Encode(Input{
		Subject:   "any_subject",
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Audience:  "any_audience",
		Issuer:    "any_issuer",
		Meta: map[string]any{
			"any_key": "any_value",
		},
	})

	assert.Nil(t, err)
	assert.NotNil(t, token)

	decode, err := jwt.Decode(token)
	assert.Nil(t, err)
	assert.Equal(t, decode.Token, token)

	assert.Equal(t, decode.Meta["any_key"], "any_value")
	assert.Equal(t, decode.Subject, "any_subject")
	assert.Equal(t, decode.Audience, "any_audience")
	assert.Equal(t, decode.Issuer, "any_issuer")

	assert.Equal(
		t,
		decode.IssuedAt.Format("2006-01-02 15:04:05"),
		issuedAt.Format("2006-01-02 15:04:05"),
	)

	assert.Equal(
		t,
		decode.ExpiresAt.Sub(decode.IssuedAt).Hours(),
		float64(2),
	)
}
