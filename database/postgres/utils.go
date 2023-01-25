package postgres

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
)

// NewNullString returns a null if s is empty, otherwise it returns
// the string which was input
func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

// NewNullTime returns a null if t is the zero value for time.Time,
// otherwise it returns the time which was input
func NewNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{
		Time:  t,
		Valid: true,
	}
}

// NewNullInt64 returns a null if i == 0, otherwise it returns
// the int64 which was input.
func NewNullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

// NewNullInt32 returns a null if i == 0, otherwise it returns
// the int32 which was input.
func NewNullInt32(i int32) sql.NullInt32 {
	if i == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: i,
		Valid: true,
	}
}

// NewNullUUID returns a null if i == uuid.Nil, otherwise it returns
// the int32 which was input.
func NewNullUUID(i uuid.UUID) uuid.NullUUID {
	if i == uuid.Nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{
		UUID:  i,
		Valid: true,
	}
}
