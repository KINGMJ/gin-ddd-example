package ctype_test

import (
	"encoding/json"
	"gin-ddd-example/internal/app/constants"
	"gin-ddd-example/internal/app/model/ctype"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewNullTime(t *testing.T) {
	now := time.Now()
	timeNow := ctype.NewNullTime(now)
	assert.True(t, timeNow.Valid)
	assert.Equal(t, now, timeNow.Time)
}

func TestNullTimeValue(t *testing.T) {
	now := time.Now()
	nullTime := ctype.NewNullTime(now)

	value, err := nullTime.Value()
	assert.NoError(t, err)
	assert.Equal(t, now.Format(constants.TIME_FORMAT), value)

	nullTime.Valid = false
	value, err = nullTime.Value()
	assert.NoError(t, err)
	assert.Nil(t, value)
}

func TestNullTimeScan(t *testing.T) {
	now := time.Now()
	var nullTime ctype.NullTime

	err := nullTime.Scan(now)
	assert.NoError(t, err)
	assert.True(t, nullTime.Valid)
	assert.Equal(t, now, nullTime.Time)

	err = nullTime.Scan(nil)
	assert.NoError(t, err)
	assert.False(t, nullTime.Valid)
	assert.Equal(t, time.Time{}, nullTime.Time)
}

func TestNullTimeMarshalJSON(t *testing.T) {
	now := time.Now()
	nullTime := ctype.NewNullTime(now)

	data, err := nullTime.MarshalJSON()
	assert.NoError(t, err)

	expectedData, err := json.Marshal(now)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)

	nullTime.Valid = false
	data, err = nullTime.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(constants.NULL_VALUE), data)
}

func TestNullTimeUnmarshalJSON(t *testing.T) {
	now := time.Now().Format(constants.TIME_FORMAT)
	var nullTime ctype.NullTime

	err := nullTime.UnmarshalJSON([]byte(now))
	assert.NoError(t, err)
	expectedTime, _ := time.ParseInLocation(constants.TIME_FORMAT, now, time.Local)
	assert.True(t, nullTime.Valid)
	assert.Equal(t, expectedTime, nullTime.Time)
}

func TestNullTimeUnmarshalJSON2(t *testing.T) {
	var nullTime ctype.NullTime
	err := nullTime.UnmarshalJSON([]byte(constants.NULL_VALUE))
	assert.NoError(t, err)
	assert.False(t, nullTime.Valid)
}
