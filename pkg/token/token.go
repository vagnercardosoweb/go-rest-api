package token

import "time"

type Input struct {
	IssuedAt  time.Time
	Subject   string
	ExpiresAt time.Time
	Audience  string
	Meta      map[string]any
	Issuer    string
}

type Decoded struct {
	Input
	Token string
}

type Token interface {
	Encode(input *Input) (string, error)
	Decode(token string) (*Decoded, error)
}
