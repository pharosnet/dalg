package logger

import (
	"log"
	"os"
	"sync"
)

var _once = sync.Once{}
var _logger *log.Logger

func Log() *log.Logger {
	_once.Do(func() {
		_logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	})
	return _logger
}
