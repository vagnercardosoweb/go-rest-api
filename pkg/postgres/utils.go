package postgres

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"time"
)

const CtxKey = "PgClientCtxKey"

func GetFromCtxOrPanic(c context.Context) *Client {
	value, exists := c.Value(CtxKey).(*Client)
	if !exists {
		panic("postgres client not found in context")
	}
	return value
}

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

func NewNullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NewNullInt32(i int32) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: i,
		Valid: true,
	}
}

func NewNullFloat64(f float64) sql.NullFloat64 {
	if f == 0 {
		return sql.NullFloat64{}
	}
	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

func NewNullUUID(i uuid.UUID) uuid.NullUUID {
	if i == uuid.Nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{
		UUID:  i,
		Valid: true,
	}
}

func NewNullBool(b bool) sql.NullBool {
	if !b {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}
