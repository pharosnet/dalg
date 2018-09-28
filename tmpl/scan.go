package tmpl

import (
	"fmt"
	"github.com/pharosnet/dalg/def"
	"io/ioutil"
	"path/filepath"
)

var _scanTpl = `
package %s

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type scanner interface {
	Scan(args ...interface{}) error
}

type NullTime struct {
	Time time.Time
	Valid bool
}

func nowTime() NullTime {
	return NullTime{time.Now(), true}
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
	Bytes []byte
	Valid bool
}

`

func WriteScanFile(dbDef *def.Db, dir string) error {
	code := fmt.Sprintf(_scanTpl, dbDef.Package)
	return ioutil.WriteFile(filepath.Join(dir, "scan.go"), []byte(code), 0666)
}
