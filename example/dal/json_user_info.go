// This file is generated by pharosnet/dalg, please don't change it by hand.
package dal

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pharosnet/dalc"
)

type UserInfo struct {
	dalc.NullJson
	Id  string `json:""`
	Age int64  `json:""`
}

func (e UserInfo) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		fmt.Fprintf(
			s,
			"{Id: %v, Age: %v}",
			e.Id, e.Age,
		)
	}
}

func (e *UserInfo) Scan(value interface{}) error {
	if value == nil {
		e.NullJson.Valid = false
		return nil
	}
	val, ok := value.([]byte)
	if !ok {
		return errors.New("UserInfo: scan failed, column type is not []byte")
	}
	if err := json.Unmarshal(val, e); err != nil {
		return fmt.Errorf("UserInfo: scan failed, unmarshal json failed, %v", err)
	}
	e.NullJson.Valid = true
	e.NullJson.Bytes = val
	return nil
}

func (e UserInfo) Value() (driver.Value, error) {
	p, err := json.Marshal(&e)
	if err != nil {
		return nil, fmt.Errorf("UserInfo: value failed, marshal json faild, %v", err)
	}
	e.NullJson.Bytes = p
	return p, nil
}
