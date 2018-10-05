package main

import (
	"github.com/pharosnet/dalg/logger"
	"github.com/pharosnet/dalg/process"
	"os"
	"time"
)

func main() {
	log := logger.Log()
	log.Println("begin to generate dal code file")
	begTime := time.Now()

	args := os.Args
	if len(args) != 3 {
		helpInfp := `usage: dalg [db def file path] [dir path of generated code files]`
		log.Printf("missing args,\n\t%s\n", helpInfp)
		return
	}
	if genErr := process.Generate(args[1], args[2]); genErr != nil {
		log.Printf("failed, \n%v", genErr)
		return
	}
	log.Printf("finish generating, cost %f sec", time.Now().Sub(begTime).Seconds())
}
