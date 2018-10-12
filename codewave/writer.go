package codewave

import (
	"bytes"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"github.com/pharosnet/dalg/logger"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func NewWriter() Writer {
	return bytes.NewBuffer([]byte{})
}

type Writer interface {
	WriteString(s string) (n int, err error)
	WriteTo(w io.Writer) (n int64, err error)
	Len() int
}

func WriteToFile(w Writer, interfaceDef def.Interface) (err error) {
	filename := filepath.Join(codeFileFolder, fmt.Sprintf("%s_%s.go", strings.TrimSpace(strings.ToLower(interfaceDef.Class)), strings.TrimSpace(strings.ToLower(interfaceDef.Name))))
	f, openErr := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if openErr != nil {
		err = openErr
		logger.Log().Println(err)
		return
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
			logger.Log().Println(err)
		}
	}()
	n, wToErr := w.WriteTo(f)
	if wToErr == nil && n < int64(w.Len()) {
		err = io.ErrShortWrite
		logger.Log().Println(err)
		return
	}
	return
}
