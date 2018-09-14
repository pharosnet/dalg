package tmpl

import (
	"bytes"
	"github.com/pharosnet/dalg/def"
	"io/ioutil"
	"path/filepath"
	"text/template"
)

var _logTpl = `
package {{.Package}}

type Log interface {
	Printf(formatter string, args ...interface{})
}

func SetLog(logger Log)  {
	if logger == nil {
		return
	}
	_logger = logger
}

var _logger Log = nil

func hasLog() bool {
	return _logger != nil
}

func logf(formatter string, args ...interface{})  {
	_logger.Printf(formatter, args...)
}
`

func WriteLoggerFile(dbDef *def.Db, dir string) error {
	tpl, tplErr := template.New("_logger").Parse(_logTpl)
	if tplErr != nil {
		return tplErr
	}
	buffer := bytes.NewBuffer([]byte{})
	if err := tpl.Execute(buffer, dbDef); err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(dir, "logger.go"), buffer.Bytes(), 0666)
}
