package tmpl

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"io/ioutil"
	"path/filepath"
	"strings"
)



func WriteTableFile(tableDef def.Interface, dir string) error {

	buffer := bytes.NewBuffer([]byte{})
	// package
	buffer.WriteString(fmt.Sprintf(`package %s\n`, tableDef.Package))
	buffer.WriteByte('\n')
	// imports
	buffer.WriteString(`import (\n`)
	buffer.WriteString(`\t"context"\n`)
	buffer.WriteString(`\t"database/sql"\n`)
	buffer.WriteString(`\t"database/sql/driver"\n`)
	buffer.WriteString(`\t"errors"\n`)
	buffer.WriteString(`\t"fmt"\n`)
	for _, importPkg := range tableDef.Imports {
		importPkg = strings.TrimSpace(importPkg)
		if importPkg != "" {
			buffer.WriteString(fmt.Sprintf(`\t"%s"\n`, importPkg))
		}
	}
	buffer.WriteString(`)\n`)
	buffer.WriteByte('\n')
	// crud sql
	tableName := strings.ToLower(strings.TrimSpace(tableDef.Name))
	tableDef.Name = tableName
	// sql, insert
	if err := buildInsertSql(&tableDef); err != nil {
		return err
	}
	// sql, update
	if err := buildUpdateSql(&tableDef); err != nil {
		return err
	}
	// sql, delete
	if err := buildDeleteSql(&tableDef); err != nil {
		return err
	}
	// sql, get one
	if err := buildGetOneSql(&tableDef); err != nil {
		return err
	}
	buffer.WriteString(`const (\n`)
	buffer.WriteString(fmt.Sprintf(`\t%sInsertSql = %s\n`, strings.TrimSpace(strings.ToLower(tableDef.MapName)), fmt.Sprintf("`%s`", tableDef.InsertSql)))
	buffer.WriteString(fmt.Sprintf(`\t%sUpdateSql = %s\n`, strings.TrimSpace(strings.ToLower(tableDef.MapName)), fmt.Sprintf("`%s`", tableDef.UpdateSql)))
	buffer.WriteString(fmt.Sprintf(`\t%sDeleteSql = %s\n`, strings.TrimSpace(strings.ToLower(tableDef.MapName)), fmt.Sprintf("`%s`", tableDef.DeleteSql)))
	buffer.WriteString(fmt.Sprintf(`\t%sGetOneSql = %s\n`, strings.TrimSpace(strings.ToLower(tableDef.MapName)), fmt.Sprintf("`%s`", tableDef.GetOneSql)))
	buffer.WriteString(`)\n`)
	buffer.WriteByte('\n')

	// extra
	for _, extra := range tableDef.ExtraType.EnumInterfaces {
		buffer.WriteString(fmt.Sprintf(`func New%s(v %s) %s {\n`, toCamel(extra.Id, true), extra.MapType, toCamel(extra.Id, true)))
		buffer.WriteString(`\tok := false\n`)
		buffer.WriteString(`\tswitch v {\n`)
		quotaMark := ""
		if extra.MapType == "string" {
			quotaMark = `"`
		}
		valueQuotaMask := ""
		if extra.OptionType == "string" {
			valueQuotaMask = `"`
		}
		hasDefault := false
		var defaultOption def.EnumOption
		for _, option := range extra.Options {
			buffer.WriteString(fmt.Sprintf(`\tcase %s%s%s:\n`, quotaMark, option.MapValue, quotaMark))
			buffer.WriteString(`\t\tok = true\n`)
			if option.Default {
				hasDefault = true
				defaultOption = option
			}
		}
		buffer.WriteString(`\tif !ok {`)
		if hasDefault {
			buffer.WriteString(fmt.Sprintf(`\t\tv = %s%s%s\n`, quotaMark, defaultOption.MapValue, quotaMark))
		} else {
			buffer.WriteString(fmt.Sprintf(`\t\tpanic(fmt.Errorf("dal: new %s failed, value is invalid"))\n`, toCamel(extra.Id, true)))
		}
		buffer.WriteString(`\t}\n`)
		buffer.WriteString(fmt.Sprintf(`\treturn %s{v, true}\n`, toCamel(extra.Id, true)))
		buffer.WriteByte('\n')

		// struct
		buffer.WriteString(fmt.Sprintf(`type %s struct {\n`, toCamel(extra.Id, true)))
		buffer.WriteString(fmt.Sprintf(`\t Value %s \n`, toType(extra.MapType)))
		buffer.WriteString(`\tValid bool\n`)
		buffer.WriteString(`}\n`)
		buffer.WriteByte('\n')
		// scan
		buffer.WriteString(fmt.Sprintf(`func (n *%s) Scan(value interface{}) error {\n`, toCamel(extra.Id, true)))
		buffer.WriteString(`\tif value == nil {\n`)
		buffer.WriteString(`\tn.Valid = false\n`)
		buffer.WriteString(`\treturn nil\n`)
		buffer.WriteString(`\t}\n`)

		buffer.WriteString(fmt.Sprintf(`\t vv, ok := value.(%s) \n`, extra.OptionType))
		buffer.WriteString(`\t if !ok { \n`)
		buffer.WriteString(fmt.Sprintf(`\t\t return fmt.Errorf("dal: %s scan value failed, value type is not %s") \n`, toCamel(extra.Id, true), extra.OptionType))
		buffer.WriteString(`\t}\n`)

		buffer.WriteString(`\tswitch vv {\n`)
		for _, option := range extra.Options {
			buffer.WriteString(fmt.Sprintf(`\t\tcase %s%s%s:\n`, valueQuotaMask, option.Value, valueQuotaMask))
			buffer.WriteString(fmt.Sprintf(`\t\t\tn.Value = %s%s%s\n`, quotaMark, option.MapValue, quotaMark))
		}
		if hasDefault {
			buffer.WriteString(`\t\tdefault:\n`)
			buffer.WriteString(fmt.Sprintf(`\t\t\tn.Value = %s%s%s\n`, quotaMark, defaultOption.MapValue, quotaMark))
		} else {
			buffer.WriteString(fmt.Sprintf(`\t\t default: \n \t\t\t return fmt.Errorf("dal: %s scan value failed, value is out of range") \n`, toCamel(extra.Id, true)))
		}
		buffer.WriteString(`\t}\n`)
		buffer.WriteString(`\tn.Valid = true \n`)
		buffer.WriteString(`\t return nil \n`)
		buffer.WriteString(`}\n`)
		buffer.WriteByte('\n')
		// value
		buffer.WriteString(fmt.Sprintf(`func (n %s) Value() (driver.Value, error) {\n`, toCamel(extra.Id, true)))
		buffer.WriteString(`\t if !n.Valid { \n`)
		buffer.WriteString(`\t\t return nil, nil \n`)
		buffer.WriteString(`\t}\n`)

		buffer.WriteString(`\t switch n.Value { \n`)
		for _, option := range extra.Options {
			buffer.WriteString(fmt.Sprintf(`\t case %s%s%s: \n`, quotaMark, option.MapValue, quotaMark))
			buffer.WriteString(fmt.Sprintf(`\t\t return %s%s%s, nil \n`, valueQuotaMask, option.Value, valueQuotaMask))
		}
		if hasDefault {
			buffer.WriteString(`\t default: \n`)
			buffer.WriteString(fmt.Sprintf(`\t\t return %s%s%s, nil  \n`, valueQuotaMask, defaultOption.Value, valueQuotaMask))
		}
		buffer.WriteString(`\t}\n`)
		buffer.WriteString(fmt.Sprintf(`\t return nil, fmt.Errorf("dal: %s value is invalid") \n`, toCamel(extra.Id, true)))
		buffer.WriteString(`}\n`)

		buffer.WriteByte('\n')
	}

	for _, extra := range tableDef.ExtraType.ElementInterfaces {
		// TODO TREE MODEL -> MAKE FLAT MAP, THEN WRITING....
		
	}
	buffer.WriteByte('\n')

	// model

	buffer.WriteByte('\n')

	// insert

	buffer.WriteByte('\n')

	// update

	buffer.WriteByte('\n')

	// delete

	buffer.WriteByte('\n')

	// get one

	buffer.WriteByte('\n')

	// query

	buffer.WriteByte('\n')

	writeFileErr := ioutil.WriteFile(filepath.Join(dir, "table_" + tableName + ".go"), buffer.Bytes(), 0666)
	if writeFileErr != nil {
		return writeFileErr
	}
	return nil
}

func buildInsertSql(tableDef *def.Interface) error {
	switch tableDef.Dialect {
	case "postgres":
		return buildPostgresInsertSql(tableDef)
	case "mysql":
		return buildMysqlInsertSql(tableDef)
	case "oracle":
		return buildOracleInsertSql(tableDef)
	default:
		return errors.New("build sql failed, unsupported dialect")
	}
	return nil
}

func buildPostgresInsertSql(tableDef *def.Interface) error {
	bb := bytes.NewBuffer([]byte{})
	bb.WriteString(`INSERT INTO "` + tableDef.Name + `" (`)
	for i, col := range tableDef.Columns {
		if i == 0 {
			bb.WriteString(`"` + strings.TrimSpace(col.Name) + `"`)
		} else {
			bb.WriteString(`, "` + strings.TrimSpace(col.Name) + `"`)
		}
	}
	bb.WriteString(`) VALUES (`)
	colLen := len(tableDef.Columns)
	for i := 1 ; i <= colLen ; i ++ {
		if i == 1 {
			bb.WriteString(fmt.Sprintf("$%d", i))
		} else {
			bb.WriteString(fmt.Sprintf(", $%d", i))
		}
	}
	bb.WriteString(`)`)
	tableDef.InsertSql = bb.String()
	return nil
}


func buildMysqlInsertSql(tableDef *def.Interface) error {

	return nil
}


func buildOracleInsertSql(tableDef *def.Interface) error {

	return nil
}

func buildUpdateSql(tableDef *def.Interface) error {
	switch tableDef.Dialect {
	case "postgres":
		return buildPostgresUpdateSql(tableDef)
	case "mysql":
		return buildMysqlUpdateSql(tableDef)
	case "oracle":
		return buildOracleUpdateSql(tableDef)
	default:
		return errors.New("build sql failed, unsupported dialect")
	}
	return nil
}

func buildPostgresUpdateSql(tableDef *def.Interface) error {
	bb := bytes.NewBuffer([]byte{})
	bb.WriteString(`UPDATE "` + tableDef.Name + `" SET `)
	i := 1
	for _, col := range tableDef.Columns {
		if col.Pk {
			continue
		}
		if i > 1 {
			bb.WriteString(", ")
		}
		if col.Version {
			if col.MapType == "int64" {
				bb.WriteString(fmt.Sprintf(`"%s" = "%s" + 1`, strings.TrimSpace(col.Name), strings.TrimSpace(col.Name)))
				continue
			}
		}
		bb.WriteString(fmt.Sprintf(`"%s" = $%d`, strings.TrimSpace(col.Name), i))
		i ++
	}
	bb.WriteString(` WHERE `)
	for pi, pk := range tableDef.Pks {
		if pi > 0 {
			bb.WriteString(` AND `)
		}
		bb.WriteString(fmt.Sprintf(`"%s" = $%d`, strings.TrimSpace(pk.Name), i))
		i ++
	}
	if tableDef.Version.MapName != "" {
		bb.WriteString(fmt.Sprintf(` AND "%s" = $%d `, tableDef.Version.Name, i))
	}
	tableDef.UpdateSql = bb.String()
	return nil
}


func buildMysqlUpdateSql(tableDef *def.Interface) error {

	return nil
}


func buildOracleUpdateSql(tableDef *def.Interface) error {

	return nil
}

func buildDeleteSql(tableDef *def.Interface) error {
	switch tableDef.Dialect {
	case "postgres":
		return buildPostgresDeleteSql(tableDef)
	case "mysql":
		return buildMysqlDeleteSql(tableDef)
	case "oracle":
		return buildOracleDeleteSql(tableDef)
	default:
		return errors.New("build sql failed, unsupported dialect")
	}
	return nil
}

func buildPostgresDeleteSql(tableDef *def.Interface) error {
	bb := bytes.NewBuffer([]byte{})
	bb.WriteString(`DELETE FROM "` + tableDef.Name + `" WHERE `)
	i := 1
	for _, col := range tableDef.Columns {
		if i > 1 {
			bb.WriteString(" AND ")
		}
		if col.Pk || col.Version {
			bb.WriteString(fmt.Sprintf(`"%s" = $%d`, strings.TrimSpace(col.Name), i))
		}
		i ++
	}
	tableDef.DeleteSql = bb.String()
	return nil
}

func buildMysqlDeleteSql(tableDef *def.Interface) error {

	return nil
}


func buildOracleDeleteSql(tableDef *def.Interface) error {

	return nil
}

func buildGetOneSql(tableDef *def.Interface) error {
	switch tableDef.Dialect {
	case "postgres":
		return buildPostgresGetOneSql(tableDef)
	case "mysql":
		return buildMysqlGetOneSql(tableDef)
	case "oracle":
		return buildOracleGetOneSql(tableDef)
	default:
		return errors.New("build sql failed, unsupported dialect")
	}
	return nil
}

func buildPostgresGetOneSql(tableDef *def.Interface) error {
	bb := bytes.NewBuffer([]byte{})
	bb.WriteString(`SELECT `)
	for i, col := range tableDef.Columns {
		if i > 0 {
			bb.WriteString(", ")
		}
		bb.WriteString(fmt.Sprintf(`"%s"`, strings.TrimSpace(col.Name)))
	}
	bb.WriteString(` FROM "` + tableDef.Name + `" WHERE `)
	for pi, pk := range tableDef.Pks {
		if pi > 0 {
			bb.WriteString(" AND ")
		}
		bb.WriteString(fmt.Sprintf(`"%s" = $%d`, strings.TrimSpace(pk.Name), pi + 1))
	}
	tableDef.GetOneSql = bb.String()
	return nil
}

func buildMysqlGetOneSql(tableDef *def.Interface) error {

	return nil
}


func buildOracleGetOneSql(tableDef *def.Interface) error {

	return nil
}

// extra type


