package tmpl

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"strings"
)

func WriteInterfaces(dbDef *def.Db, dir string) error {
	enums := make([]def.Interface, 0, 1)
	jsons := make([]def.Interface, 0, 1)
	for _, interfaceDef := range dbDef.Interfaces {
		class := strings.ToLower(strings.TrimSpace(interfaceDef.Class))
		if class == "" {
			return errors.New("def file is invalid, interface class is empty")
		}
		interfaceDef.Package = dbDef.Package
		interfaceDef.Imports = make([]string, 0, 1)
		if class == "table" {
			interfaceDef.Dialect = dbDef.Dialect
			interfaceDef.Pks = make([]def.Column, 0, 1)
			interfaceDef.CommonColumns = make([]def.Column, 0, 1)
			for _, col := range interfaceDef.Columns {
				if col.Pk {
					interfaceDef.Pks = append(interfaceDef.Pks, col)
				} else if col.Version {
					interfaceDef.Version = col
				} else {
					interfaceDef.CommonColumns = append(interfaceDef.CommonColumns, col)
				}
				if pos := strings.LastIndexByte(col.MapType, '.'); pos > 0 {
					interfaceDef.Imports = append(interfaceDef.Imports, col.MapType[0:pos])
				}
			}
			interfaceDef.PkNum = int64(len(interfaceDef.Pks))
			if err := WriteTableOrViewFile(interfaceDef, dir); err != nil {
				return err
			}
		} else if class == "view" {
			for _, col := range interfaceDef.Columns {
				if pos := strings.LastIndexByte(col.MapType, '.'); pos > 0 {
					interfaceDef.Imports = append(interfaceDef.Imports, col.MapType[0:pos])
				}
			}
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
		}
		return fmt.Errorf("invalid interface class, %s", class)
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
		if aheadUp && i == 0 && p[0] >= 97 && p[0] <= 122 {
			ss += string(p)
			continue
		}
		if p[0] >= 97 && p[0] <= 122 {
			p[0] = p[0] - 32
		}
		ss += strings.TrimSpace(string(p))
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
		if importPkg != "" {
			buffer.WriteString(fmt.Sprintf(`\t"%s"\n`, importPkg))
		}
	}
	return buffer.String()
}