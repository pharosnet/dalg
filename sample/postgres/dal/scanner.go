package dal

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type scanner interface {
	Scan(args ...interface{}) error
}

func nowTime() NullTime {
	return NullTime{time.Now(), true}
}

type NullTime struct {
	Time time.Time
	Valid bool
}

func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		n.Time, n.Valid = time.Time{}, false
		return nil
	}
	switch value.(type) {
	case time.Time:
		n.Time = value.(time.Time)
		n.Valid = true
	}
	return nil
}

func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}


type NullJson struct {
	Bytes []byte `json:"-"`
	Valid bool `json:"-"`
}

func (n *NullJson) Scan(value interface{}) error {
	if value == nil {
		n.Bytes, n.Valid = nil, false
		return nil
	}
	switch value.(type) {
	case []byte:
		n.Bytes = value.([]byte)
		n.Valid = true
	case string:
		n.Bytes = []byte(value.(string))
		n.Valid = true
	}
	return nil
}

func (n NullJson) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bytes, nil
}

func jsonScannerMarshal(n *NullJson, v interface{}) (err error) {
	n.Bytes, err = json.Marshal(v)
	if err != nil {
		return
	}
	n.Valid = true
	return
}

func jsonScannerUnmarshal(n *NullJson, v interface{}) (err error) {
	if n.Bytes == nil || len(n.Bytes) == 0 {
		n.Valid = false
		return
	}
	if err = json.Unmarshal(n.Bytes, v); err != nil {
		return
	}
	n.Valid = true
	return
}