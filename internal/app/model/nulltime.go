package model

import (
	"database/sql/driver"
	"time"
)

type NullTime struct {
	Time  time.Time
	Valid bool
}

func NewValidNullTime(t time.Time) NullTime {
	return NullTime{Time: t, Valid: true}
}

func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time.Format("2006-01-02 15:04:05"), nil
}

func (n *NullTime) Scan(value any) error {
	if value == nil {
		n.Time, n.Valid = time.Time{}, false
	} else {
		n.Time, n.Valid = value.(time.Time), true
	}
	return nil
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		b := make([]byte, 0, 21)
		b = append(b, '"')
		b = n.Time.AppendFormat(b, "2006-01-02 15:04:05")
		b = append(b, '"')
		return b, nil
	} else {
		return []byte("null"), nil
	}
}

func (n *NullTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}
	now, err := time.ParseInLocation(`"2006-01-02 15:04:05"`, string(data), time.Local)
	*n = NullTime{
		Time:  now,
		Valid: true,
	}
	return
}
