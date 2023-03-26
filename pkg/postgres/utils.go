package postgres

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type JSONToMap map[string]any

func (j *JSONToMap) Scan(value any) error {
	if value == nil {
		return nil
	}
	var data = value.([]byte)
	return json.Unmarshal(data, &j)
}

func (j *JSONToMap) Value() (driver.Value, error) {
	if len(*j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
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

func NewNullUUID(i uuid.UUID) uuid.NullUUID {
	if i == uuid.Nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{
		UUID:  i,
		Valid: true,
	}
}
