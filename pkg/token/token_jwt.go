package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

type TokenJwt struct {
	secretKey []byte
}

func NewJwt(secretKey []byte) Token {
	if len(secretKey) == 0 {
		secretKey = []byte(env.Required("JWT_SECRET_KEY"))
	}

	return &TokenJwt{secretKey: secretKey}
}

func (j *TokenJwt) Encode(input Input) (string, error) {
	if input.Subject == "" {
		return "", errors.New("sub needs to be filled to create a token")
	}

	if input.IssuedAt.IsZero() {
		input.IssuedAt = time.Now()
	}

	if input.ExpiresAt.IsZero() {
		expiresInSecond, _ := strconv.Atoi(env.Get("JWT_EXPIRES_IN_SECONDS", "86400"))
		input.ExpiresAt = time.Now().Add(time.Second * time.Duration(expiresInSecond))
	}

	claims := jwt.MapClaims{
		"sub":  input.Subject,
		"iat":  input.IssuedAt.Unix(),
		"exp":  input.ExpiresAt.Unix(),
		"meta": input.Meta,
	}

	if input.Issuer != "" {
		claims["iss"] = input.Issuer
	}

	if input.Audience != "" {
		claims["aud"] = input.Audience
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.SignedString(j.secretKey)
	return signedString, err
}

func (j *TokenJwt) Decode(token string) (*Output, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("parse jwt claims error")
	}

	err = claims.Valid()
	if err != nil {
		return nil, err
	}

	output := &Output{
		Token: token,
		Input: Input{
			Subject:   claims["sub"].(string),
			IssuedAt:  time.Unix(int64(claims["iat"].(float64)), 0),
			ExpiresAt: time.Unix(int64(claims["exp"].(float64)), 0),
			Issuer:    claims["iss"].(string),
			Audience:  claims["aud"].(string),
		},
	}

	if _, ok := claims["meta"].(map[string]any); ok {
		output.Meta = claims["meta"].(map[string]any)
	}

	return output, nil
}
