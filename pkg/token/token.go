package token

import (
	"time"
)

type Input struct {
	IssuedAt  time.Time
	Subject   string
	ExpiresAt time.Time
	Audience  string
	Meta      map[string]any
	Issuer    string
}

type Output struct {
	Input
	Token string
}

type Client interface {
	Encode(input *Input) (*Output, error)
	Decode(token string) (*Output, error)
}
