package dal

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