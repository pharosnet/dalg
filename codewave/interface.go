package codewave

import (
	"bytes"
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
	tables := make([]*def.Interface, 0, 1) // ddl
	views := make([]*def.Interface, 0, 1)
	enums := make([]*def.Interface, 0, 1)
	jsonObjects := make([]*def.Interface, 0, 1)
	for _, interfaceDef := range db.Interfaces {
		class := strings.ToLower(strings.TrimSpace(interfaceDef.Class))
		if class == "" {
			err = errors.New("def file is invalid, interface class is empty")
			logger.Log().Println(err)
			return
		}
		interfaceDef.Imports = make([]string, 0, 1)
		if class == "table" {
			interfaceDef.Dialect = db.Dialect
			tables = append(tables, &interfaceDef)
		} else if class == "view" {
			views = append(views, &interfaceDef)
		} else if class == "enum" {
			if pos := strings.LastIndexByte(interfaceDef.MapType, '.'); pos > 0 {
				interfaceDef.Imports = append(interfaceDef.Imports, interfaceDef.MapType[0:pos])
			}
			enums = append(enums, &interfaceDef)
		} else if class == "json" {
			for _, col := range interfaceDef.Fields {
				if pos := strings.LastIndexByte(col.MapType, '.'); pos > 0 {
					interfaceDef.Imports = append(interfaceDef.Imports, col.MapType[0:pos])
				}
			}
			jsonObjects = append(jsonObjects, &interfaceDef)
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

func toType(src string) string {
	ss := src
	if pos := strings.LastIndexByte(src, '/'); pos > 0 {
		ss = src[pos+1:]
	}
	return ss
}

func importsCode(imports []string) string {
	if imports == nil || len(imports) == 0 {
		return ""
	}
	importsMap := make(map[string]int)
	for _, importName := range imports {
		importsMap[importName] = 1
	}
	buffer := bytes.NewBuffer([]byte{})
	for importPkg := range importsMap {
		importPkg = strings.TrimSpace(importPkg)
		if importPkg != "" && importPkg != "sql" {
			buffer.WriteString(fmt.Sprintf(`	"%s" `, importPkg))
			buffer.WriteString("\n")
		}
	}
	return buffer.String()
}