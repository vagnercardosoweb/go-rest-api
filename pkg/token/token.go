package token

import (
	"context"
	"fmt"
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

const CtxClientKey = "TokenClientKey"
const CtxDecodedKey = "DecodedTokenKey"

func DecodedFromCtx(ctx context.Context) *Output {
	value, ok := ctx.Value(CtxDecodedKey).(*Output)

	if !ok {
		panic(fmt.Errorf(`context key "%s" does not exist`, CtxDecodedKey))
	}

	return value
}

func ClientFromCtx(ctx context.Context) Client {
	value, ok := ctx.Value(CtxClientKey).(Client)

	if !ok {
		panic(fmt.Errorf(`context key "%s" does not exist`, CtxClientKey))
	}

	return value
}
