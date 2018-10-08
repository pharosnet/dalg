package codewave

import (
	"fmt"
	"github.com/pharosnet/dalg/def"
	"github.com/pharosnet/dalg/logger"
	"strconv"
	"strings"
)

func waveEnums(enums []*def.Interface) error {
	for _, enum := range enums {
		if err := waveEnum(enum); err != nil {
			logger.Log().Println(err)
			return err
		}
	}
	return nil
}

func waveEnum(enum *def.Interface) error {
	w := NewWriter()
	// intro
	waveIntroduction(w)
	// package
	wavePackage(w, enum.Package)
	// imports
	waveEnumImports(w, enum.Imports)
	// vars
	waveEnumVars(w, enum)
	// struct
	waveEnumStruct(w, enum)
	// scan
	waveEnumScan(w, enum)
	// value
	waveEnumValue(w, enum)

	return WriteToFile(w, enum)
}

func waveEnumImports(w Writer, imports []string)  {
	imports = append(imports, "database/sql")
	imports = append(imports, "database/sql/driver")
	imports = append(imports, "errors")
	imports = append(imports, "fmt")
	waveImports(w, imports)
}

func waveEnumVars(w Writer, enum *def.Interface)  {
	dataMask := ""
	if enum.MapType == "string" {
		dataMask = `"`
	}
	oriMask := ""
	typeUpper := strings.ToUpper(enum.Type)
	oriTypeIsString := strings.Contains(typeUpper, "VARCHAR") || strings.Contains(typeUpper, "CHAR") || strings.Contains(typeUpper, "TINYTEXT")
	oriTypeIsString = oriTypeIsString || strings.Contains(typeUpper, "CHARACTER") || strings.Contains(typeUpper, "VARCHAR2")
	oriTypeIsString = oriTypeIsString || strings.Contains(typeUpper, "NVARCHAR") || strings.Contains(typeUpper, "NVARCHAR2")
	oriTypeIsString = oriTypeIsString || strings.Contains(typeUpper, "NCHAR") || strings.Contains(typeUpper, "NVARCHAR2")
	if oriTypeIsString {
		oriMask = `"`
	}
	w.WriteString(`var (`)
	w.WriteString("\n")
	for _, opt := range enum.Options {
		w.WriteString(fmt.Sprintf(`	%s%s = %s{Data: %s%s%s, Origin: %s%s%s, Valid: true} `,
			toCamel(enum.MapName, true), toCamel(opt.MapValue, true),
			dataMask, opt.MapValue, dataMask,
			oriMask, opt.Value, oriMask,
			))
		w.WriteString("\n")
	}
	w.WriteString(`)`)
	w.WriteString("\n")
	w.WriteString("\n")
}

// origin enum type only supports string, bool and int64
func waveEnumStruct(w Writer, enum *def.Interface)  {
	w.WriteString(fmt.Sprintf(`type %s struct {`, toCamel(enum.MapName, true)))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	Data %s`, enum.MapType))
	w.WriteString("\n")
	oriType := "string"
	if _, err := strconv.Atoi(enum.Options[0].Value); err == nil {
		oriType = "int64"
	}
	if _, err := strconv.ParseBool(enum.Options[0].Value); err == nil {
		oriType = "bool"
	}
	w.WriteString(fmt.Sprintf(`	Origin %s`, oriType))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	Valid %s`, "bool"))
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
	// fmt
	w.WriteString(fmt.Sprintf(`func (e %s) Format(s fmt.State, verb rune) {`, toCamel(enum.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	switch verb {`)
	w.WriteString("\n")
	w.WriteString(`	case 'v':`)
	w.WriteString("\n")
	w.WriteString(`		fmt.Fprintf(`)
	w.WriteString("\n")
	w.WriteString(`			s,`)
	w.WriteString("\n")
	w.WriteString(`			{%v %v},`)
	w.WriteString("\n")
	w.WriteString(`			e.Data, e.Origin,`)
	w.WriteString("\n")
	w.WriteString(`		)`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

func waveEnumScan(w Writer, enum *def.Interface)  {
	dataMask := ""
	if enum.MapType == "string" {
		dataMask = `"`
	}
	oriMask := `"`
	if _, err := strconv.Atoi(enum.Options[0].Value); err == nil {
		oriMask = ""
	}
	if _, err := strconv.ParseBool(enum.Options[0].Value); err == nil {
		oriMask = ""
	}
	w.WriteString(fmt.Sprintf(`func (e *%s) Scan(value interface{}) error {`, toCamel(enum.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	if value == nil {`)
	w.WriteString("\n")
	w.WriteString(`		e.Valid = false`)
	w.WriteString("\n")
	w.WriteString(`		return nil`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	val, ok := value.(%s)`, enum.Type))
	w.WriteString("\n")
	w.WriteString(`	if !ok {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		return errors.New("%s: scan failed, value type is invalid")`, toCamel(enum.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	e.Origin = val`)
	w.WriteString("\n")
	w.WriteString(`	switch val {`)
	w.WriteString("\n")
	for _, opt := range enum.Options {
		w.WriteString(fmt.Sprintf(`	case %s%s%s:`, oriMask, opt.Value, oriMask))
		w.WriteString("\n")
		w.WriteString(fmt.Sprintf(`		e.Data = %s%s%s`, dataMask, opt.MapValue, dataMask))
		w.WriteString("\n")
	}
	w.WriteString(`	default:`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		return errors.New("%s: scan failed, value is out of range")`, toCamel(enum.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	e.Valid = true`)
	w.WriteString("\n")
	w.WriteString(`	return nil`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

func waveEnumValue(w Writer, enum *def.Interface)  {
	dataMask := ""
	if enum.MapType == "string" {
		dataMask = `"`
	}
	oriMask := `"`
	if _, err := strconv.Atoi(enum.Options[0].Value); err == nil {
		oriMask = ""
	}
	if _, err := strconv.ParseBool(enum.Options[0].Value); err == nil {
		oriMask = ""
	}
	w.WriteString(fmt.Sprintf(`func (e %s) Value() (driver.Value, error) {`, toCamel(enum.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	if !e.Valid {`)
	w.WriteString("\n")
	w.WriteString(`		return nil, nil`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	switch e.Data {`)
	w.WriteString("\n")
	for _, opt := range enum.Options {
		w.WriteString(fmt.Sprintf(`	case %s%s%s:`, dataMask, opt.MapValue, dataMask))
		w.WriteString("\n")
		w.WriteString(fmt.Sprintf(`		e.Origin = %s%s%s`, oriMask, opt.Value, oriMask))
		w.WriteString("\n")
	}
	w.WriteString(`	default:`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		return errors.New("%s: to sql driver value failed, value is out of range")`, toCamel(enum.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	return e.Origin, nil`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}