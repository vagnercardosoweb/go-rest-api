package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"

	"github.com/golang-jwt/jwt"
)

func getExpiresAt() int64 {
	expiresInSecond, err := strconv.Atoi(env.Get("JWT_EXPIRES_IN_SECONDS", "0"))
	if err != nil || expiresInSecond <= 0 {
		expiresInSecond = int(time.Hour) * 24
	}
	return time.Now().Add(time.Duration(expiresInSecond)).Unix()
}

func getSecretKey() []byte {
	return []byte(env.Required("JWT_SECRET_KEY"))
}

func New(subject string) (string, error) {
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		Subject:   subject,
		ExpiresAt: getExpiresAt(),
		Issuer:    "internal",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.SignedString(getSecretKey())
	return signedString, err
}

func Verify(externalToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(externalToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return getSecretKey(), nil
	})
	return token, err
}
