package tmpl

import (
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"strings"
)

func WriteInterfaces(dbDef *def.Db, dir string) error {
	for _, interfaceDef := range dbDef.Interfaces {
		_type := strings.ToLower(strings.TrimSpace(interfaceDef.Type))
		if _type == "" {
			return errors.New("def file is invalid, interface type is empty")
		}
		interfaceDef.Package = dbDef.Package
		interfaceDef.Dialect = dbDef.Dialect
		interfaceDef.EnableNil = dbDef.EnableNil
		interfaceDef.Pks = make([]def.Column, 0, 1)
		interfaceDef.CommonColumns = make([]def.Column, 0, 1)
		interfaceDef.ByteElementColumns = make([]def.Column, 0, 1)
		interfaceDef.Imports = make([]string, 0, 1)
		for _, col := range interfaceDef.Columns {
			if col.Pk {
				interfaceDef.Pks = append(interfaceDef.Pks, col)
			} else if col.Version {
				interfaceDef.Version = col
			} else if col.Json || col.Xml {
				interfaceDef.ByteElementColumns = append(interfaceDef.ByteElementColumns, col)
				for i, extra := range interfaceDef.ExtraType.ElementInterfaces {
					if extra.Id == col.MapName {
						if col.Json {
							extra.Type = "json"
						} else if col.Xml {
							extra.Type = "xml"
						}
						interfaceDef.ExtraType.ElementInterfaces[i] = extra
					}
				}
			} else {
				interfaceDef.CommonColumns = append(interfaceDef.CommonColumns, col)
			}
			if pos := strings.LastIndexByte(col.MapType, '.'); pos > 0 {
				interfaceDef.Imports = append(interfaceDef.Imports, col.MapType[0:pos])
			}
		}
		interfaceDef.PkNum = int64(len(interfaceDef.Pks))
		for _, extra := range interfaceDef.ExtraType.ElementInterfaces {
			for _, field := range extra.Fields {
				if pos := strings.LastIndexByte(field.MapType, '.'); pos > 0 {
					interfaceDef.Imports = append(interfaceDef.Imports, field.MapType[0:pos])
				}
			}
		}
		if _type == "table" {
			if err := WriteTableFile(interfaceDef, dir); err != nil {
				return err
			}
		} else if _type == "view" {
			if err := WriteViewFile(interfaceDef, dir); err != nil {
				return err
			}
		}
		return fmt.Errorf("invalid interface type, %s", _type)
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