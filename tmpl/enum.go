package tmpl

import (
	"bytes"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"io/ioutil"
	"path/filepath"
)

func WriteEnumFile(enumDefs []def.Interface, dir string) error {
	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(_notes + "\n")
	// package
	pknName := enumDefs[0].Package
	buffer.WriteString(fmt.Sprintf(`package %s` + "\n", pknName))
	buffer.WriteString("\n")
	// imports
	imports := make([]string, 0, 1)
	for _, enumDef := range enumDefs {
		for _, importName := range enumDef.Imports {
			imports = append(imports, importName)
		}
	}
	buffer.WriteString(`import (` + "\n")
	buffer.WriteString( "\t" + `"database/sql/driver"` + "\n")
	buffer.WriteString( "\t" + `"errors"` + "\n")
	buffer.WriteString(importsCode(imports))
	buffer.WriteString(`)` + "\n")
	buffer.WriteString("\n")

	for _, enumDef := range enumDefs {
		buffer.WriteString(fmt.Sprintf(`func New%s(v %s) %s {` + "\n", toCamel(enumDef.MapName, true), enumDef.MapType, toCamel(enumDef.MapName, true)))
		buffer.WriteString( "\t" + `ok := false` + "\n")
		buffer.WriteString( "\t" + `switch v {` + "\n")
		quotaMark := ""
		if enumDef.MapType == "string" {
			quotaMark = `"`
		}
		valueQuotaMask := ""
		if enumDef.Type == "string" {
			valueQuotaMask = `"`
		}
		hasDefault := false
		var defaultOption def.EnumOption
		for _, option := range enumDef.Options {
			buffer.WriteString(fmt.Sprintf( "\t" + `case %s%s%s:` + "\n", quotaMark, option.MapValue, quotaMark))
			buffer.WriteString( "\t" +  "\t" + `ok = true` + "\n")
			if option.Default {
				hasDefault = true
				defaultOption = option
			}
		}
		buffer.WriteString( "\t" + `if !ok {` + "\n")
		if hasDefault {
			buffer.WriteString(fmt.Sprintf( "\t" +  "\t" + `v = %s%s%s` + "\n", quotaMark, defaultOption.MapValue, quotaMark))
		} else {
			buffer.WriteString(fmt.Sprintf( "\t" +  "\t" + `panic(errors.New("dal: new %s failed, value is invalid"))` + "\n", toCamel(enumDef.MapName, true)))
		}
		buffer.WriteString( "\t" + `}` + "\n")
		buffer.WriteString( "\t" + `}` + "\n")
		buffer.WriteString(fmt.Sprintf( "\t" + `return %s{v, true}` + "\n", toCamel(enumDef.MapName, true)))
		buffer.WriteString( "\t" + `}` + "\n")
		buffer.WriteString("\n")

		// struct
		buffer.WriteString(fmt.Sprintf(`type %s struct {` + "\n", toCamel(enumDef.MapName, true)))
		buffer.WriteString(fmt.Sprintf( "\t" + ` Data %s ` + "\n", toType(enumDef.MapType)))
		buffer.WriteString( "\t" + `Valid bool` + "\n")
		buffer.WriteString(`}` + "\n")
		buffer.WriteString("\n")
		// scan
		buffer.WriteString(fmt.Sprintf(`func (n *%s) Scan(value interface{}) error {` + "\n", toCamel(enumDef.MapName, true)))
		buffer.WriteString( "\t" + `if value == nil {` + "\n")
		buffer.WriteString( "\t" + `n.Valid = false` + "\n")
		buffer.WriteString( "\t" + `return nil` + "\n")
		buffer.WriteString( "\t" + `}` + "\n")

		buffer.WriteString(fmt.Sprintf( "\t" + ` vv, ok := value.(%s) ` + "\n", enumDef.Type))
		buffer.WriteString( "\t" + ` if !ok { ` + "\n")
		buffer.WriteString(fmt.Sprintf( "\t" +  "\t" + ` return errors.New("dal: %s scan value failed, value type is not %s") ` + "\n", toCamel(enumDef.MapName, true), enumDef.Type))
		buffer.WriteString( "\t" + `}` + "\n")

		buffer.WriteString( "\t" + `switch vv {` + "\n")
		for _, option := range enumDef.Options {
			buffer.WriteString(fmt.Sprintf( "\t" +  "\t" + `case %s%s%s:` + "\n", valueQuotaMask, option.Value, valueQuotaMask))
			buffer.WriteString(fmt.Sprintf( "\t" +  "\t" +  "\t" + `n.Data = %s%s%s` + "\n", quotaMark, option.MapValue, quotaMark))
		}
		if hasDefault {
			buffer.WriteString( "\t" +  "\t" + `default:` + "\n")
			buffer.WriteString(fmt.Sprintf( "\t" +  "\t" +  "\t" + `n.Data = %s%s%s` + "\n", quotaMark, defaultOption.MapValue, quotaMark))
		} else {
			buffer.WriteString(fmt.Sprintf( "\t" +  "\t" + ` default: \n \t\t\t return errors.New("dal: %s scan value failed, value is out of range") ` + "\n", toCamel(enumDef.MapName, true)))
		}
		buffer.WriteString( "\t" + `}` + "\n")
		buffer.WriteString( "\t" + `n.Valid = true ` + "\n")
		buffer.WriteString( "\t" + ` return nil ` + "\n")
		buffer.WriteString(`}` + "\n")
		buffer.WriteString("\n")
		// value
		buffer.WriteString(fmt.Sprintf(`func (n %s) Value() (driver.Value, error) {` + "\n", toCamel(enumDef.MapName, true)))
		buffer.WriteString( "\t" + ` if !n.Valid { ` + "\n")
		buffer.WriteString( "\t" +  "\t" + ` return nil, nil ` + "\n")
		buffer.WriteString( "\t" + `}` + "\n")

		buffer.WriteString( "\t" + ` switch n.Data { ` + "\n")
		for _, option := range enumDef.Options {
			buffer.WriteString(fmt.Sprintf( "\t" + ` case %s%s%s: ` + "\n", quotaMark, option.MapValue, quotaMark))
			buffer.WriteString(fmt.Sprintf( "\t" +  "\t" + ` return %s%s%s, nil ` + "\n", valueQuotaMask, option.Value, valueQuotaMask))
		}
		if hasDefault {
			buffer.WriteString( "\t" + ` default: ` + "\n")
			buffer.WriteString(fmt.Sprintf( "\t" +  "\t" + ` return %s%s%s, nil  ` + "\n", valueQuotaMask, defaultOption.Value, valueQuotaMask))
		}
		buffer.WriteString( "\t" + `}` + "\n")
		if !hasDefault {
			buffer.WriteString(fmt.Sprintf( "\t" + ` return nil, errors.New("dal: %s value is invalid") ` + "\n", toCamel(enumDef.MapName, true)))
		}
		buffer.WriteString(`}` + "\n")

		buffer.WriteString("\n")
	}

	writeFileErr := ioutil.WriteFile(filepath.Join(dir, "extra_enum.go"), buffer.Bytes(), 0666)
	if writeFileErr != nil {
		return writeFileErr
	}
	return nil
}
