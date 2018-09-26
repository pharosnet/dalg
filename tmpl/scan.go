package tmpl

import (
	"bytes"
	"github.com/pharosnet/dalg/def"
	"io/ioutil"
	"path/filepath"
	"text/template"
)

var _scanTpl = `
package {{.Package}}

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Scanable interface {
	Scan(dest ...interface{}) error
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

`

func WriteScanFile(dbDef *def.Db, dir string) error {
	tpl, tplErr := template.New("_scanTpl").Parse(_scanTpl)
	if tplErr != nil {
		return tplErr
	}
	buffer := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buffer, dbDef); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dir, "scan.go"), buffer.Bytes(), 0666)
}
