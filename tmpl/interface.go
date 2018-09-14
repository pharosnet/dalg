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
