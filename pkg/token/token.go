package token

import (
	"context"
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

const OutputCtxKey = "TokenOutputCtxKey"
const ClientCtxKey = "TokenClientCtxKey"

func GetOutputFromCtxOrPanic(ctx context.Context) *Output {
	decoded, ok := ctx.Value(OutputCtxKey).(*Output)
	if !ok {
		panic("token output not found in context")
	}
	return decoded
}

func GetClientFromCtxOrPanic(ctx context.Context) Client {
	client, ok := ctx.Value(ClientCtxKey).(Client)
	if !ok {
		panic("token client not found in context")
	}
	return client
}
