package main

import (
	"encoding/json"
	"github.com/pharosnet/dalg/sample/postgres/dal"
	"log"
	"os"
)

type SSS struct {
	dal.NullJson

		Id string
		Age int64
}

func (s *SSS) Element() interface{} {
	return s
}


func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags | log.Llongfile)
	s1 := SSS{Id: "1", Age: 1}
	bb, mErr := json.Marshal(&s1)
	if mErr != nil {
		logger.Println(mErr)
		return
	}
	logger.Println(string(bb))
	s2 := &SSS{}
	sErr := s2.Scan(bb)
	if sErr != nil {
		logger.Println(sErr)
		return
	}
	logger.Println(s2)

}
