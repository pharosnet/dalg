package def

import (
	"encoding/xml"
	"fmt"
	"github.com/pharosnet/dalg/logger"
	"io/ioutil"
)

func Read(filePath string) (db *Db, err error) {
	log := logger.Log()
	bb, readFileErr := ioutil.ReadFile(filePath)
	if readFileErr != nil {
		err = fmt.Errorf("read def file failed, %v\n", readFileErr)
		return
	}
	db = new(Db)
	xmlUnmarshalErr := xml.Unmarshal(bb, db)
	if readFileErr != nil {
		err = fmt.Errorf("read def file failed, %v\n", xmlUnmarshalErr)
		return
	}
	log.Println("load def file success")
	return
}
