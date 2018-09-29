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
	buffer.WriteString(fmt.Sprintf(`type %s struct {\n `, toCamel(tableDef.MapName, true)))
	for _, col := range tableDef.Columns {
		buffer.WriteString(fmt.Sprintf(`\t%s %s \n`, toCamel(col.MapName, true), strings.TrimSpace(col.MapType)))
	}
	buffer.WriteString(`}\n`)
	buffer.WriteByte('\n')
	buffer.WriteString(fmt.Sprintf(`func (row *%s) Format(s fmt.State, verb rune) { \n`, toCamel(tableDef.MapName, true)))
	buffer.WriteString(`\t switch verb { \n`)
	buffer.WriteString(`\t case 'v': \n`)
	buffer.WriteString(`\t\t switch { \n`)
	buffer.WriteString(`\t\t case s.Flag('+'): \n`)
	buffer.WriteString(`\t\t\t fmt.Fprintf(s, "(`)
	for i, col := range tableDef.Columns {
		if i > 0 {
			buffer.WriteString(`, `)
		}
		buffer.WriteString(fmt.Sprintf(`%s: `, toCamel(col.MapName, true)))
		buffer.WriteString(`%v`)

	}
	buffer.WriteString(`)",\n \t\t\t\t`)
	for i, col := range tableDef.Columns {
		if i > 0 {
			buffer.WriteString(`, `)
		}
		buffer.WriteString(fmt.Sprintf(`row.%s`, toCamel(col.MapName, true)))

	}
	buffer.WriteString(`)\n`)
	buffer.WriteString(`\t\t default:\n`)
	buffer.WriteString(`\t\t\t fmt.Fprintf(s, "&{`)
	for i, _ := range tableDef.Columns {
		if i > 0 {
			buffer.WriteString(`, `)
		}
		buffer.WriteString(`%v`)
	}

	buffer.WriteString(`}",\n`)
	buffer.WriteString(`\t\t\t\t`)
	for i, col := range tableDef.Columns {
		if i > 0 {
			buffer.WriteString(`, `)
		}
		buffer.WriteString(fmt.Sprintf(`row.%s`, toCamel(col.MapName, true)))

	}
	buffer.WriteString(`)\n`)
	buffer.WriteString(`\t\t\t } \n`)
	buffer.WriteString(`\t\t } \n`)
	buffer.WriteString(`} \n`)

	buffer.WriteByte('\n')

	buffer.WriteString(fmt.Sprintf(`func scan%s(sa scanner) (row *%s, err error) {\n`, toCamel(tableDef.MapName, true), toCamel(tableDef.MapName, true)))
	buffer.WriteString(fmt.Sprintf(`\trow = &%s{}\n`, toCamel(tableDef.MapName, true)))
	buffer.WriteString(`\t scanErr := sa.Scan( \n`)
	for _, col := range tableDef.Columns {
		buffer.WriteString(fmt.Sprintf(`\t\t\t &row.%s ,\n`, toCamel(col.MapName, true)))
	}
	buffer.WriteString(`\t\t)\n`)
	buffer.WriteString(`\t if scanErr != nil { \n`)
	buffer.WriteString(`\t\t err = fmt.Errorf("dal: scan failed. reason: %v", scanErr) \n`)
	buffer.WriteString(`\t\t return \n`)
	buffer.WriteString(`\t } \n`)
	buffer.WriteString(`\t return \n`)
	buffer.WriteString(`}\n`)

	buffer.WriteByte('\n')

	buffer.WriteString(fmt.Sprintf(`type %sRangeFn func(ctx context.Context, row *%s, err error) error \n`, toCamel(tableDef.MapName, true), toCamel(tableDef.MapName, true)))

	buffer.WriteByte('\n')

	// crud todo
	if tableDef.Class == "table" {
		// insert
		buffer.WriteString(fmt.Sprintf(`func Insert%s(ctx context.Context, rows ...*%s) (affected int64, err error) { \n`, toCamel(tableDef.MapName, true), toCamel(tableDef.MapName, true)))
		buffer.WriteString(`\t if ctx == nil { \n`)
		buffer.WriteString(fmt.Sprintf(`\t\t err = errors.New("dal: insert %s failed, context is empty") \n`, toCamel(tableDef.MapName, true)))
		buffer.WriteString(`\t\t return \n`)
		buffer.WriteString(`\t\ } \n`)
		buffer.WriteString(`\t if rows == nil || len(rows) == 0 { \n`)
		buffer.WriteString(fmt.Sprintf(`\t\t err = errors.New("dal: insert %s failed, row is empty") \n`, toCamel(tableDef.MapName, true)))
		buffer.WriteString(`\t\t return \n`)
		buffer.WriteString(`\t\ } \n`)
		buffer.WriteString(fmt.Sprintf(`\t stmt, prepareErr := prepare(ctx).PrepareContext(ctx, %sInsertSql) \n`, toCamel(tableDef.MapName, false)))
		buffer.WriteString(`\t if prepareErr != nil { \n`)
		buffer.WriteString(fmt.Sprintf(`\t\t err = fmt.Errorf("dal: insert %s failed, prepared statment failed. reason: %s", prepareErr) \n`, toCamel(tableDef.MapName, true), "%v"))
		buffer.WriteString(`\t\t return \n`)
		buffer.WriteString(`\t\ } \n`)

		buffer.WriteString(`\t defer func() { \n`)
		buffer.WriteString(`\t\t stmtCloseErr := stmt.Close() \n`)
		buffer.WriteString(`\t\t if stmtCloseErr != nil { \n`)
		buffer.WriteString(fmt.Sprintf(`\t\t\t err = fmt.Errorf("dal: insert %s failed, close prepare statment failed. reason: %s", stmtCloseErr) \n`, toCamel(tableDef.MapName, true), "%v"))
		buffer.WriteString(`\t\t\t return \n`)
		buffer.WriteString(`\t\t } \n`)
		buffer.WriteString(`\t }() \n`)

		buffer.WriteString(`\t for _, row := range rows { \n`)
		buffer.WriteString(`\t\t result, execErr :=  stmt.ExecContext(ctx, `)
		for i, col := range tableDef.Columns {
			if i > 0 {
				buffer.WriteString(`, `)
			}
			buffer.WriteString(fmt.Sprintf(`row.%s`, toCamel(col.MapName, true)))
		}
		buffer.WriteString(`)\n`)

		buffer.WriteString(`\t\t if execErr != nil { \n`)
		buffer.WriteString(fmt.Sprintf(`\t\t\t err = fmt.Errorf("dal: insert %s failed, execute statment failed. reason: %s", execErr)\n`, toCamel(tableDef.MapName, true), "%v"))
		buffer.WriteString(`\t\t\t return \n`)
		buffer.WriteString(`\t\t } \n`)

		buffer.WriteString(`\t\t affectedRows, affectedErr :=  result.RowsAffected() \n`)
		buffer.WriteString(`\t\t if affectedErr != nil { \n`)
		buffer.WriteString(fmt.Sprintf(`\t\t\t err = fmt.Errorf("dal: insert %s failed, get rows affected failed. reason: %s", affectedErr) \n`, toCamel(tableDef.MapName, true), "%v"))
		buffer.WriteString(`\t\t\t return \n`)
		buffer.WriteString(`\t\t } \n`)

		buffer.WriteString(`\t\t if affectedRows == 0 { \n`)
		buffer.WriteString(fmt.Sprintf(`\t\t\t err = errors.New("dal: insert %s failed, no rows affected") \n`, toCamel(tableDef.MapName, true)))
		buffer.WriteString(`\t\t\t return \n`)
		buffer.WriteString(`\t\t } \n`)

		buffer.WriteString(`\t\t affected = affected + affectedRows \n`)
		if len(tableDef.Pks) == 1 && tableDef.Pks[0].Increment {
			buffer.WriteString(fmt.Sprintf(`\t\t %s, get%sErr := result.LastInsertId() \n`, toCamel(tableDef.Pks[0].MapName, false), toCamel(tableDef.Pks[0].MapName, true)))
			buffer.WriteString(fmt.Sprintf(`\t\t if get%sErr != nil { \n`, toCamel(tableDef.Pks[0].MapName, true)))
			buffer.WriteString(fmt.Sprintf(`\t\t\t err = fmt.Errorf("dal: insert %s failed, get last insert pk failed. reason: %s", get%sErr) \n`,
				toCamel(tableDef.MapName, true), "%v", toCamel(tableDef.Pks[0].MapName, true)))
			buffer.WriteString(`\t\t\t return \n`)
			buffer.WriteString(`\t\t } \n`)
			buffer.WriteString(fmt.Sprintf(`\t\t if %s < 0 { \n`, toCamel(tableDef.Pks[0].MapName, false)))
			buffer.WriteString(fmt.Sprintf(`\t\t\t err = errors.New("dal: insert %s failed, get last insert pk failed. pk is invalid") \n`, toCamel(tableDef.MapName, true)))
			buffer.WriteString(`\t\t\t return \n`)
			buffer.WriteString(`\t\t } \n`)
			buffer.WriteString(fmt.Sprintf(`\r\r row.%s = %s \n `, toCamel(tableDef.Pks[0].MapName, true), toCamel(tableDef.Pks[0].MapName, false)))
			buffer.WriteByte('\n')
		}

		buffer.WriteString(`\t\t if hasLog() { \n`)
		buffer.WriteString(fmt.Sprintf(`\t\t\t logf("dal: insert %s success, sql : %s, row : %s\n", %sInsertSql, row) \n`,
			toCamel(tableDef.MapName, true), "%+v", toCamel(tableDef.MapName, false)))
		buffer.WriteString(`\t\t } \n`)
		buffer.WriteString(`\t} \n`)
		buffer.WriteString(`\t return \n`)
		buffer.WriteString(`} \n`)

		buffer.WriteByte('\n')

		// update

		buffer.WriteByte('\n')

		// delete

		buffer.WriteByte('\n')

		// get one

		buffer.WriteByte('\n')
	}

	// query todo

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



