package tmpl

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"strings"
)

func WriteInterfaces(dbDef *def.Db, dir string) error {
	enums := make([]def.Interface, 0, 8)
	jsons := make([]def.Interface, 0, 8)
	for _, interfaceDef := range dbDef.Interfaces {
		class := strings.ToLower(strings.TrimSpace(interfaceDef.Class))
		if class == "" {
			return errors.New("def file is invalid, interface class is empty")
		}
		interfaceDef.Package = dbDef.Package
		interfaceDef.Imports = make([]string, 0, 1)
		if class == "table" {
			interfaceDef.Dialect = dbDef.Dialect
			if err := WriteTableOrViewFile(interfaceDef, dir); err != nil {
				return err
			}
			if dbDef.DDL {
				if err := WriteSqlDDL(interfaceDef, dir); err != nil {
					return err
				}
			}
		} else if class == "view" {
			if err := WriteTableOrViewFile(interfaceDef, dir); err != nil {
				return err
			}
		} else if class == "enum" {
			if pos := strings.LastIndexByte(interfaceDef.MapType, '.'); pos > 0 {
				interfaceDef.Imports = append(interfaceDef.Imports, interfaceDef.MapType[0:pos])
			}
			enums = append(enums, interfaceDef)
		} else if class == "json" {
			for _, col := range interfaceDef.Fields {
				if pos := strings.LastIndexByte(col.MapType, '.'); pos > 0 {
					interfaceDef.Imports = append(interfaceDef.Imports, col.MapType[0:pos])
				}
			}
			jsons = append(jsons, interfaceDef)
		} else {
			return fmt.Errorf("invalid interface class, %s", class)
		}
	}
	if len(enums) > 0 {
		if err := WriteEnumFile(enums, dir); err != nil {
			return err
		}
	}
	if len(jsons) > 0 {
		if err := WriteJsonFile(jsons, dir); err != nil {
			return err
		}
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
			buffer.WriteString(fmt.Sprintf("\t"+`"%s"`+"\n", importPkg))
		}
	}
	return buffer.String()
}
