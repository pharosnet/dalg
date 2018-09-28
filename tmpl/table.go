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



func WriteTableOrViewFile(tableDef def.Interface, dir string) error {

	if tableDef.Class == "table" {
		tableDef.Pks = make([]def.Column, 0, 1)
		tableDef.CommonColumns = make([]def.Column, 0, 1)
		for _, col := range tableDef.Columns {
			if col.Pk {
				tableDef.Pks = append(tableDef.Pks, col)
			} else if col.Version {
				tableDef.Version = col
			} else {
				tableDef.CommonColumns = append(tableDef.CommonColumns, col)
			}
			if pos := strings.LastIndexByte(col.MapType, '.'); pos > 0 {
				tableDef.Imports = append(tableDef.Imports, col.MapType[0:pos])
			}
		}
		tableDef.PkNum = int64(len(tableDef.Pks))
	} else if tableDef.Class == "view" {
		for _, col := range tableDef.Columns {
			if pos := strings.LastIndexByte(col.MapType, '.'); pos > 0 {
				tableDef.Imports = append(tableDef.Imports, col.MapType[0:pos])
			}
		}
	}

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
	buffer.WriteString(importsCode(tableDef.Imports))
	buffer.WriteString(`)\n`)
	buffer.WriteByte('\n')

	if tableDef.Class == "table" {
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
		buffer.WriteString(fmt.Sprintf(`\t%sInsertSql = %s\n`, toCamel(tableDef.MapName, false), fmt.Sprintf("`%s`", tableDef.InsertSql)))
		buffer.WriteString(fmt.Sprintf(`\t%sUpdateSql = %s\n`, toCamel(tableDef.MapName, false), fmt.Sprintf("`%s`", tableDef.UpdateSql)))
		buffer.WriteString(fmt.Sprintf(`\t%sDeleteSql = %s\n`, toCamel(tableDef.MapName, false), fmt.Sprintf("`%s`", tableDef.DeleteSql)))
		buffer.WriteString(fmt.Sprintf(`\t%sGetOneSql = %s\n`, toCamel(tableDef.MapName, false), fmt.Sprintf("`%s`", tableDef.GetOneSql)))
		buffer.WriteString(`)\n`)
		buffer.WriteByte('\n')
	}

	// model

	buffer.WriteByte('\n')

	if tableDef.Class == "table" {
		// insert

		buffer.WriteByte('\n')

		// update

		buffer.WriteByte('\n')

		// delete

		buffer.WriteByte('\n')

		// get one

		buffer.WriteByte('\n')
	}

	// query

	buffer.WriteByte('\n')

	writeFileErr := ioutil.WriteFile(filepath.Join(dir, strings.TrimSpace(strings.ToLower(tableDef.Class)) + "_" + strings.TrimSpace(strings.ToLower(tableDef.Name)) + ".go"), buffer.Bytes(), 0666)
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



