package logs

import (
	"log"
	"os"
)

var logger *log.Logger = nil

func NewLog(verbose bool) {
	if verbose {
		logger = log.New(os.Stdout, "dalc: ", log.LstdFlags|log.Lshortfile)
	} else {
		newMuteLog()
	}
}

func newMuteLog() {
	devNull, openNullErr := os.Open(os.DevNull)
	if openNullErr != nil {
		logger = log.New(os.Stdout, "dalc: ", log.LstdFlags|log.Lshortfile)
	} else {

	}
	logger = log.New(devNull, "dalc: ", log.LstdFlags|log.Lshortfile)
}

func Log() *log.Logger {
	return logger
}

func init() {
	newMuteLog()
}
