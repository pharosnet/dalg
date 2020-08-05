package generates

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pharosnet/dalg/logs"
)

type GenerateFile struct {
	Name    string
	Content []byte
}

func writeFiles(out string, fs []*GenerateFile) (err error) {

	for _, f := range fs {

		formatted, fErr := format.Source(f.Content)
		if fErr != nil {
			err = fmt.Errorf("go fmt %s failed, %v, code:\n%s", f.Name, fErr, string(f.Content))
			return
		}
		f.Content = formatted

	}

	fileInfo, fileErr := os.Lstat(out)
	if fileErr != nil {
		if os.IsNotExist(fileErr) {
			mkdirErr := os.MkdirAll(out, os.ModePerm)
			if mkdirErr != nil {
				err = fmt.Errorf("make dir %s failed, %v", out, mkdirErr)
				return
			}
		} else {
			err = fmt.Errorf("open dir %s failed, %v", out, fileErr)
			return
		}
	}
	if fileInfo != nil && !fileInfo.IsDir() {
		err = fmt.Errorf("open  %s failed, it is not a dir", out)
		return
	} else {
		files, readErr := ioutil.ReadDir(out)
		if readErr != nil {
			err = fmt.Errorf("open  %s failed, read dir, %v", out, readErr)
			return
		}
		for _, file := range files {
			name := file.Name()
			if strings.ToLower(name[len(name)-3:]) == ".go" {
				rmErr := os.Remove(filepath.Join(out, file.Name()))
				if rmErr != nil {
					err = fmt.Errorf("remove %s failed, %v", name, rmErr)
					return
				}
			}
		}
	}

	for _, f := range fs {

		fileName := filepath.Join(out, f.Name)

		wErr := ioutil.WriteFile(fileName, f.Content, os.ModePerm)
		if wErr != nil {
			err = fmt.Errorf("write file %s filed, %v", f.Name, wErr)
			return
		}
		logs.Log().Println("write file succeed", fileName)
	}

	return
}
