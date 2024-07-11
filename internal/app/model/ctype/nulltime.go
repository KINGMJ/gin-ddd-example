package ctype

// 自定义的类型
import (
	"database/sql"
	"encoding/json"
	"gin-ddd-example/internal/app/constants"
	"time"
)

type NullTime struct {
	sql.NullTime
}

func NewNullTime(t time.Time) NullTime {
	return NullTime{NullTime: sql.NullTime{Time: t, Valid: true}}
}

func (n NullTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
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
