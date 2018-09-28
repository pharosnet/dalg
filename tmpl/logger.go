package tmpl

import (
	"fmt"
	"github.com/pharosnet/dalg/def"
	"io/ioutil"
	"path/filepath"
)

var _logTpl = `
package %s

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
	code := fmt.Sprintf(_logTpl, dbDef.Package)
	return ioutil.WriteFile(filepath.Join(dir, "logger.go"), []byte(code), 0666)
}
