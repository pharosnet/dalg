package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type SQLFile struct {
	Name    string
	Content []byte
}

func ReadFiles(filePath string) (s []*SQLFile, err error) {

	f, openErr := os.Open(filePath)
	if openErr != nil {
		err = fmt.Errorf("open %s failed, %v", filePath, openErr)
		return
	}

	info, statErr := f.Stat()
	if statErr != nil {
		err = fmt.Errorf("get %s stat failed, %v", filePath, statErr)
		return
	}

	filePaths := make([]string, 0, 1)
	if info.IsDir() {
		fs, dirErr := ioutil.ReadDir(filePath)
		if dirErr != nil {
			err = fmt.Errorf("read %s dir failed, %v", filePath, dirErr)
			return
		}
		for _, fileInfo := range fs {
			if fileInfo.IsDir() {
				continue
			}
			filePaths = append(filePaths, filepath.Join(filePath, fileInfo.Name()))
		}
	} else {
		filePaths = append(filePaths, filePath)
	}

	if len(filePaths) == 0 {
		err = fmt.Errorf("read %s failed, no files", filePath)
		return
	}

	s = make([]*SQLFile, 0, 1)
	for _, path := range filePaths {
		data, readErr := ioutil.ReadFile(path)
		if readErr != nil {
			err = fmt.Errorf("read %s file failed, %v", path, readErr)
			return
		}
		s = append(s, &SQLFile{
			Name:    filepath.Base(path),
			Content: data,
		})
	}

	return
}
