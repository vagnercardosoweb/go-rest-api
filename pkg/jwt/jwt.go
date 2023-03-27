package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"

	"github.com/golang-jwt/jwt"
)

type Payload struct {
	Sub string
	Iat float64
	Exp float64
	Iss string
}

func expiresAt() int64 {
	expiresInSecond, err := strconv.Atoi(env.Get("JWT_EXPIRES_IN_SECONDS", "0"))
	if err != nil || expiresInSecond <= 0 {
		expiresInSecond = int(time.Hour) * 24
	}
	return time.Now().Add(time.Second * time.Duration(expiresInSecond)).Unix()
}

func secretKey() []byte {
	return []byte(env.Required("JWT_SECRET_KEY"))
}

func Encode(subject string) (string, error) {
	claims := jwt.MapClaims{
		"iat": time.Now().Unix(),
		"sub": subject,
		"exp": expiresAt(),
		"iss": "internal",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.SignedString(secretKey())
	return signedString, err
}

func Decode(externalToken string) (*Payload, error) {
	var payload = &Payload{}

	token, err := jwt.Parse(externalToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secretKey(), nil
	})

	if err != nil {
		return payload, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return payload, errors.New("parse jwt claims error")
	}

	err = claims.Valid()
	if err != nil {
		return payload, err
	}

	payload.Sub = claims["sub"].(string)
	payload.Iat = claims["iat"].(float64)
	payload.Exp = claims["exp"].(float64)
	payload.Iss = claims["iss"].(string)

	return payload, nil
}
