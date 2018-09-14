package def

import (
	"encoding/xml"
	"fmt"
	"github.com/pharosnet/dalg/logger"
	"io/ioutil"
	"strings"
)

func Read(filePath string) (db *Db, err error) {
	log := logger.Log()
	bb, readFileErr := ioutil.ReadFile(filePath)
	if readFileErr != nil {
		err = fmt.Errorf("read def file failed, %v\n", readFileErr)
		return
	}
	db = new(Db)
	xmlUnmarshalErr :=  xml.Unmarshal(bb, db)
	if readFileErr != nil {
		err = fmt.Errorf("read def file failed, %v\n", xmlUnmarshalErr)
		return
	}
	if strings.ToLower(strings.TrimSpace(db.Dialect)) == "" {
		err = fmt.Errorf("def file is invalid, dialect is empty\n")
		return
	}
	if db.Dialect != "mysql" && db.Dialect != "postgres" && db.Dialect != "oracle" {
		err = fmt.Errorf("def file is invalid, dialect is not supported\n")
		return
	}
	if strings.TrimSpace(db.Package) == "" {
		db.Package = "dal"
		log.Println("package attr in def file is undefined, use default package, named dal")
	}
	if db.Interfaces == nil || len(db.Interfaces) == 0 {
		err = fmt.Errorf("def file is invalid, no interfaces definended\n")
		return
	}
	log.Println("read def file success")
	return
}
