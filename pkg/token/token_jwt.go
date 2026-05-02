package token

import (
	"errors"
	"strings"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	expiresIn time.Duration
	secretKey []byte
}

type jwtClaims struct {
	Meta map[string]any `json:"meta,omitempty"`
	jwt.RegisteredClaims
}

const minJWTSecretKeyLength = 32

func NewJwt(secretKey []byte) *Jwt {
	return &Jwt{secretKey: secretKey, expiresIn: time.Hour * 24}
}

func JwtFromEnv() *Jwt {
	secretKey := strings.TrimSpace(env.Required("JWT_SECRET_KEY"))
	expiresIn := env.GetAsInt("JWT_EXPIRES_IN_SECONDS", "86400")

	if len(secretKey) < minJWTSecretKeyLength {
		panic(errors.New("environment \"JWT_SECRET_KEY\" must have at least 32 characters"))
	}

	return &Jwt{
		secretKey: []byte(secretKey),
		expiresIn: time.Duration(expiresIn) * time.Second,
	}
}

func (j *Jwt) Encode(input *Input) (*Output, error) {
	if input.Subject == "" {
		return nil, errors.New("sub needs to be filled to create a token")
	}

	if input.IssuedAt.IsZero() {
		input.IssuedAt = time.Now()
	}

	if input.ExpiresAt.IsZero() {
		input.ExpiresAt = time.Now().Add(j.expiresIn)
	}

	if input.Issuer == "" {
		input.Issuer = "go"
	}

	claims := &jwtClaims{
		Meta: input.Meta,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   input.Subject,
			IssuedAt:  jwt.NewNumericDate(input.IssuedAt),
			ExpiresAt: jwt.NewNumericDate(input.ExpiresAt),
			Issuer:    input.Issuer,
		},
	}

	if input.Audience != "" {
		claims.Audience = jwt.ClaimStrings{input.Audience}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString(j.secretKey)

	return &Output{Input: *input, Token: signedString}, err
}

func (j *Jwt) Decode(token string) (*Output, error) {
	claims := new(jwtClaims)

	jwtToken, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return j.secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, errors.New("invalid jwt token")
	}

	if claims.Subject == "" {
		return nil, errors.New("missing jwt subject")
	}

	if claims.IssuedAt == nil {
		return nil, errors.New("missing jwt issued at")
	}

	if claims.ExpiresAt == nil {
		return nil, errors.New("missing jwt expires at")
	}

	audience := ""
	if len(claims.Audience) > 0 {
		audience = claims.Audience[0]
	}

	output := &Output{
		Token: token,
		Input: Input{
			IssuedAt:  claims.IssuedAt.Time,
			ExpiresAt: claims.ExpiresAt.Time,
			Subject:   claims.Subject,
			Audience:  audience,
			Issuer:    claims.Issuer,
			Meta:      claims.Meta,
		},
	}

	return output, nil
}
