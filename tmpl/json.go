package tmpl

import (
	"bytes"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"io/ioutil"
	"path/filepath"
)

func WriteJsonFile(jsonDefs []def.Interface, dir string) error {
	buffer := bytes.NewBuffer([]byte{})
	// package
	pknName := jsonDefs[0].Package
	buffer.WriteString(fmt.Sprintf(`package %s\n`, pknName))
	buffer.WriteByte('\n')
	// imports
	imports := make([]string, 0, 1)
	for _, jsonDef := range jsonDefs {
		for _, importName := range jsonDef.Imports {
			imports = append(imports, importName)
		}
	}
	buffer.WriteString(`import (\n`)
	buffer.WriteString(`\t"database/sql/driver"\n`)
	buffer.WriteString(`\t"encoding/json"\n`)
	buffer.WriteString(`\t"errors"\n`)
	buffer.WriteString(importsCode(imports))
	buffer.WriteString(`)\n`)
	buffer.WriteByte('\n')

	for _, jsonDef := range jsonDefs {
		buffer.WriteString(fmt.Sprintf(`type %s struct {\n`, toCamel(jsonDef.MapName, true)))
		buffer.WriteString(`\tNullJson ` + "`" + `json:"-"` + "`\n")
		for _, field := range jsonDef.Fields {
			buffer.WriteString(fmt.Sprintf(`\t%s %s %sjson:"%s"%s\n`, toCamel(field.MapName, true), field.MapType, "`", field.Name, "`"))
		}
		buffer.WriteString(`}\n`)
		buffer.WriteString(`\n`)
		// scan
		buffer.WriteString(fmt.Sprintf(`func (n *%s) Scan(value interface{}) error {\n`, toCamel(jsonDef.MapName, true)))
		buffer.WriteString(`\tif value == nil { \n`)
		buffer.WriteString(`\t\tn.NullJson.Bytes, n.NullJson.Valid = nil, false \n`)
		buffer.WriteString(`\t\treturn nil \n`)
		buffer.WriteString(`\t} \n`)
		buffer.WriteString(`\tswitch value.(type) { \n`)
		buffer.WriteString(`\tcase []byte: \n`)
		buffer.WriteString(`\t\tn.NullJson.Bytes = value.([]byte) \n`)
		buffer.WriteString(`\tcase string: \n`)
		buffer.WriteString(`\t\tn.NullJson.Bytes = []byte(value.(string)) \n`)
		buffer.WriteString(`\tdefault: \n`)
		buffer.WriteString(fmt.Sprintf(`\t\t return errors.New("%s scan value failed, value type is invalid") \n`, toCamel(jsonDef.MapName, true)))
		buffer.WriteString(`\t} \n`)
		buffer.WriteString(`\tif err := json.Unmarshal(n.NullJson.Bytes, n); err != nil { \n`)
		buffer.WriteString(`\t\treturn nil \n`)
		buffer.WriteString(`\t} \n`)
		buffer.WriteString(`\tn.NullJson.Valid = true \n`)
		buffer.WriteString(`\treturn nil \n`)
		buffer.WriteString(`} \n`)
		buffer.WriteString(`\n`)
		// value
		buffer.WriteString(fmt.Sprintf(`func (n %s) Value() (driver.Value, error) {\n`, toCamel(jsonDef.MapName, true)))
		buffer.WriteString(`\t if !n.NullJson.Valid { \n`)
		buffer.WriteString(`\t\t return nil, nil \n`)
		buffer.WriteString(`\t } \n`)
		buffer.WriteString(`\t p, err := json.Marshal(&n) \n`)
		buffer.WriteString(`\tif err != nil { \n`)
		buffer.WriteString(`\t\t return nil, err \n`)
		buffer.WriteString(`\t} \n`)
		buffer.WriteString(`\t return p, nil \n`)
		buffer.WriteString(`} \n`)
		buffer.WriteString(`\n`)
	}

	writeFileErr := ioutil.WriteFile(filepath.Join(dir, "extra_json.go"), buffer.Bytes(), 0666)
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil
}
