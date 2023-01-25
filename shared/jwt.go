package shared

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	secretKey       = EnvRequiredByName("JWT_SECRET_KEY")
	expiresInSecond = EnvGetByName("JWT_EXPIRES_IN_SECONDS", "0")
)

func getExpiresAt() int64 {
	expiresInSecond, err := strconv.Atoi(expiresInSecond)
	if err != nil || expiresInSecond <= 0 {
		expiresInSecond = int(time.Hour) * 24
	}
	return time.Now().Add(time.Duration(expiresInSecond)).Unix()
}

func JwtGenerateBySubject(subject string) (string, error) {
	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		Subject:   subject,
		ExpiresAt: getExpiresAt(),
		Issuer:    "internal",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.SignedString([]byte(secretKey))
	return signedString, err
}
