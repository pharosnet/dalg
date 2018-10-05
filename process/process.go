package process

import (
	"fmt"
	"github.com/pharosnet/dalg/def"
	"github.com/pharosnet/dalg/logger"
	"github.com/pharosnet/dalg/tmpl"
	"os"
	"path/filepath"
)

func Generate(defPath string, destDirPath string) error {
	dbDef, readDefErr := def.Read(defPath)
	if readDefErr != nil {
		return readDefErr
	}
	return generate0(destDirPath, dbDef)
}

func generate0(destDirPath string, dbDef *def.Db) error {
	log := logger.Log()
	destDirPath = filepath.Join(destDirPath, dbDef.Package)
	mkdirErr := os.MkdirAll(destDirPath, os.ModePerm)
	if mkdirErr != nil {
		return fmt.Errorf("mkdir failed, [%v], %v\n", destDirPath, mkdirErr)
	}
	log.Printf("mkdir success, %s\n", destDirPath)
	if err := tmpl.WriteLoggerFile(dbDef, destDirPath); err != nil {
		return err
	}
	if err := tmpl.WriteContextFile(dbDef, destDirPath); err != nil {
		return err
	}
	if err := tmpl.WriteScanFile(dbDef, destDirPath); err != nil {
		return err
	}
	if err := tmpl.WriteInterfaces(dbDef, destDirPath); err != nil {
		return err
	}
	if err := formatAndVet(destDirPath); err != nil {
		return err
	}
	return nil
}
