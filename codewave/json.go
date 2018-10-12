package codewave

import (
	"fmt"
	"github.com/pharosnet/dalg/def"
	"github.com/pharosnet/dalg/logger"
)

func waveJsonObjects(jsonObjects []def.Interface) error {
	for _, jsonObject := range jsonObjects {
		if jsonObject.Name == "" {
			jsonObject.Name = toUnderScore(jsonObject.MapName)
		}
		for i, col := range jsonObject.Fields {
			pkg, mapType := parseCustomizeType(col.MapType)
			col.MapType = mapType
			if pkg != "" {
				jsonObject.Imports = append(jsonObject.Imports, pkg)
			}
			jsonObject.Fields[i] = col
		}
		if err := waveJsonObject(jsonObject); err != nil {
			logger.Log().Println(err)
			return err
		}
	}
	return nil
}

func waveJsonObject(jsonObject def.Interface) error {
	w := NewWriter()
	// intro
	waveIntroduction(w)
	// package
	wavePackage(w, jsonObject.Package)
	// imports
	waveJsonObjectImports(w, jsonObject.Imports)
	// struct
	waveJsonObjectStruct(w, jsonObject)
	// scan
	waveJsonObjectScan(w, jsonObject)
	// value
	waveJsonObjectValue(w, jsonObject)
	return WriteToFile(w, jsonObject)
}

func waveJsonObjectImports(w Writer, imports []string) {
	imports = append(imports, "database/sql/driver")
	imports = append(imports, "encoding/json")
	imports = append(imports, "errors")
	imports = append(imports, "fmt")
	imports = append(imports, "github.com/pharosnet/dalc")
	waveImports(w, imports)
}

func waveJsonObjectStruct(w Writer, jsonObject def.Interface) {
	w.WriteString(fmt.Sprintf(`type %s struct {`, toCamel(jsonObject.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`dalc.NullJson`)
	w.WriteString("\n")
	format := "{"
	formatArgs := ""
	for i, field := range jsonObject.Fields {
		w.WriteString(fmt.Sprintf(`	%s  %s `, toCamel(field.MapName, true), field.MapType))
		w.WriteString("`" + fmt.Sprintf(`json:"%s"`, field.Name) + "`")
		w.WriteString("\n")
		if i > 0 {
			format = format + ", "
		}
		format = format + fmt.Sprintf(`%s: %s`, toCamel(field.MapName, true), "%v")
		formatArgs = formatArgs + fmt.Sprintf(`e.%s, `, toCamel(field.MapName, true))
	}
	format = format + "}"
	w.WriteString("}")
	w.WriteString("\n")
	w.WriteString("\n")
	// fmt
	w.WriteString(fmt.Sprintf(`func (e %s) Format(s fmt.State, verb rune) {`, toCamel(jsonObject.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	switch verb {`)
	w.WriteString("\n")
	w.WriteString(`	case 'v':`)
	w.WriteString("\n")
	w.WriteString(`		fmt.Fprintf(`)
	w.WriteString("\n")
	w.WriteString(`			s,`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`			"%s",`, format))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`			%s`, formatArgs))
	w.WriteString("\n")
	w.WriteString(`		)`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

func waveJsonObjectScan(w Writer, jsonObject def.Interface) {
	w.WriteString(fmt.Sprintf(`func (e *%s) Scan(value interface{}) error {`, toCamel(jsonObject.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	if value == nil {`)
	w.WriteString("\n")
	w.WriteString(`		e.NullJson.Valid = false`)
	w.WriteString("\n")
	w.WriteString(`		return nil`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	val, ok := value.([]byte)`)
	w.WriteString("\n")
	w.WriteString(`	if !ok {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		return errors.New("%s: scan failed, column type is not []byte")`, toCamel(jsonObject.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	if err := json.Unmarshal(val, e); err != nil {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		return fmt.Errorf("%s: scan failed, unmarshal json failed, %s", err)`, toCamel(jsonObject.MapName, true), "%v"))
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	e.NullJson.Valid = true`)
	w.WriteString("\n")
	w.WriteString(`	e.NullJson.Bytes = val`)
	w.WriteString("\n")
	w.WriteString(`	return nil`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

func waveJsonObjectValue(w Writer, jsonObject def.Interface) {
	w.WriteString(fmt.Sprintf(`func (e %s) Value() (driver.Value, error) {`, toCamel(jsonObject.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	p, err := json.Marshal(&e)`)
	w.WriteString("\n")
	w.WriteString(`	if err != nil {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		return nil, fmt.Errorf("%s: value failed, marshal json faild, %s", err)`, toCamel(jsonObject.MapName, true), "%v"))
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	e.NullJson.Bytes = p`)
	w.WriteString("\n")
	w.WriteString(`	return p, nil`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}
