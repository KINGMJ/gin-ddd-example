package ctype

// 自定义的类型
import (
	"database/sql"
	"database/sql/driver"
	"gin-ddd-example/internal/app/constants"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func NewNullTime(t time.Time) NullTime {
	return NullTime{NullTime: sql.NullTime{Time: t, Valid: true}}
}

func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	// 对时间进行格式化处理
	return n.Time.Format(constants.TIME_FORMAT), nil
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
		// 这里处理掉了时区信息，比如 2024-02-02 00:00:00 +0800 CS 这种格式
		b := make([]byte, 0, 21)
		b = append(b, '"')
		b = n.Time.AppendFormat(b, "2006-01-02 15:04:05")
		b = append(b, '"')
		return b, nil
	} else {
		// json.Marshal 对于 nil 值序列化后为 null
		// 这里我们自己实现，提升效率
		return []byte(constants.NULL_VALUE), nil
	}
}

func (n *NullTime) UnmarshalJSON(data []byte) (err error) {
	if string(data) == constants.NULL_VALUE {
		return nil
	}
	now, err := time.ParseInLocation(constants.TIME_FORMAT, string(data), time.Local)
	*n = NewNullTime(now)
	return
}
