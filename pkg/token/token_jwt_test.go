package token

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var secretKey = []byte("secret")
var jwtInstance Client

func init() {
	jwtInstance = NewJwt(secretKey)
}

func TestTokenEncodeWithSubject(t *testing.T) {
	output, err := jwtInstance.Encode(&Input{Subject: "any_sub"})
	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Len(t, strings.Split(output.Token, "."), 3)
}

func TestTokenEncodeWithoutSubject(t *testing.T) {
	_, err := jwtInstance.Encode(&Input{})
	assert.NotNil(t, err)
}

func TestTokenDecode(t *testing.T) {
	issuedAt := time.Now()
	expiresAt := issuedAt.Add(time.Hour * 2)

	output, err := jwtInstance.Encode(&Input{
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
	assert.NotNil(t, output)

	decode, err := jwtInstance.Decode(output.Token)
	assert.Nil(t, err)
	assert.Equal(t, decode.Token, output.Token)

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
