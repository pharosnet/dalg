package codewave

import (
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"github.com/pharosnet/dalg/logger"
	"log"
	"strings"
)

var codeFileFolder = ""

func Wave(db *def.Db, dir string) (err error) {
	if db == nil {
		err = errors.New("wave code failed, non-db definition")
		logger.Log().Println(err)
		return
	}
	if strings.ToLower(strings.TrimSpace(db.Dialect)) == "" {
		err = fmt.Errorf("def file is invalid, dialect is empty")
		logger.Log().Println(err)
		return
	}
	if db.Dialect != "mysql" && db.Dialect != "postgres" && db.Dialect != "oracle" {
		err = fmt.Errorf("def file is invalid, dialect is not supported")
		logger.Log().Println(err)
		return
	}
	if strings.TrimSpace(db.Package) == "" {
		db.Package = "dal"
		log.Println("package attr in def file is undefined, use default package, named dal")
	}
	if db.Interfaces == nil || len(db.Interfaces) == 0 {
		err = fmt.Errorf("def file is invalid, no interfaces definended")
		logger.Log().Println(err)
		return
	}
	codeFileFolder = strings.TrimSpace(dir)
	tables := make([]def.Interface, 0, 8) // ddl
	views := make([]def.Interface, 0, 8)
	enums := make([]def.Interface, 0, 8)
	jsonObjects := make([]def.Interface, 0, 8)
	for _, interfaceDef := range db.Interfaces {
		interfaceDef.Package = db.Package
		class := strings.ToLower(strings.TrimSpace(interfaceDef.Class))
		if class == "" {
			err = errors.New("def file is invalid, interface class is empty")
			logger.Log().Println(err)
			return
		}
		interfaceDef.Imports = make([]string, 0, 1)
		if class == "table" {
			interfaceDef.Dialect = db.Dialect
			interfaceDef.Owner = db.Owner
			interfaceDef.Tablespace = db.Tablespace
			tables = append(tables, interfaceDef)
		} else if class == "view" {
			views = append(views, interfaceDef)
		} else if class == "enum" {
			enums = append(enums, interfaceDef)
		} else if class == "json" {
			jsonObjects = append(jsonObjects, interfaceDef)
		} else {
			err = fmt.Errorf("invalid interface class, %s", class)
			logger.Log().Println(err)
			return
		}
	}
	if waveErr := waveTables(tables); waveErr != nil {
		err = waveErr
		return err
	}
	if db.DDL {
		if waveErr := waveDDL(tables); waveErr != nil {
			err = waveErr
			return err
		}
	}
	if waveErr := waveViews(views); waveErr != nil {
		err = waveErr
		return err
	}
	if waveErr := waveEnums(enums); waveErr != nil {
		err = waveErr
		return err
	}
	if waveErr := waveJsonObjects(jsonObjects); waveErr != nil {
		err = waveErr
		return err
	}
	return nil
}

func toCamel(src string, aheadUp bool) string {
	slices := strings.Split(src, "_")
	ss := ""
	for i, slice := range slices {
		p := []byte(slice)
		if i == 0 {
			if aheadUp && p[0] >= 97 && p[0] <= 122 {
				p[0] = p[0] - 32
			} else if !aheadUp && p[0] >= 65 && p[0] <= 90 {
				p[0] = p[0] + 32
			}
		} else {
			if p[0] >= 97 && p[0] <= 122 {
				p[0] = p[0] - 32
			}
		}
		ss += string(p)
	}
	return ss
}

func toUnderScore(src string) string {
	pp := make([]byte, 0, 8)
	p := []byte(strings.TrimSpace(src))
	for i, char := range p {
		if char >= 65 && char <= 90 {
			if i > 0 {
				pp = append(pp, '_')
			}
			pp = append(pp, char + 32)
		} else {
			pp = append(pp, char)
		}
	}
	return string(pp)
}

func parseCustomizeType(typeName string) (pkg string, name string) {
	if pos := strings.LastIndexByte(typeName, '.'); pos > 0 {
		pkg = typeName[0:pos]
	}
	if pos := strings.LastIndexByte(typeName, '/'); pos > 0 {
		name = typeName[pos+1:]
	}
	if name == "" {
		name = typeName
	}
	return
}
