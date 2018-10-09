package codewave

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"strings"
)

func waveTableUpdate(w Writer, table *def.Interface) {
	ql, qlErr := buildUpdateSql(table)
	if qlErr != nil {
		panic(qlErr)
	}
	w.WriteString(fmt.Sprintf("const %sUpdateSQL = `%s` \n", toCamel(table.MapName, false), ql))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`func %sUpdate(ctx context.Context, rows ...*%s) (affected int64, err error) {`, toCamel(table.MapName, true), toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	fn := func(ctx context.Context, stmt *sql.Stmt, arg interface{}) (result sql.Result, err error) {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		row, ok := arg.(*%s)`, toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`		if !ok {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`			err = errors.New("%s: update failed, invalid type")`, toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`			return`)
	w.WriteString("\n")
	w.WriteString(`		}`)
	w.WriteString("\n")
	w.WriteString(`		result, err = stmt.Exec(`)
	for _, col := range table.CommonColumns {
		w.WriteString(fmt.Sprintf(`			&row.%s,`, toCamel(col.MapName, true)))
		w.WriteString("\n")
	}
	for _, col := range table.Pks {
		w.WriteString(fmt.Sprintf(`			&row.%s,`, toCamel(col.MapName, true)))
		w.WriteString("\n")
	}
	if table.Version.MapName != "" {
		w.WriteString(fmt.Sprintf(`			&row.%s,`, toCamel(table.Version.MapName, true)))
		w.WriteString("\n")
	}
	w.WriteString(`		)`)
	w.WriteString("\n")
	w.WriteString(`		return`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	affected, err = dalc.Execute(ctx, %sUpdateSQL, fn, %sArrayMapToInterfacs(rows)...)`, toCamel(table.MapName, false), toCamel(table.MapName, false)))
	w.WriteString("\n")
	w.WriteString(`	return`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

func buildUpdateSql(table *def.Interface) (ql string, err error) {
	switch table.Dialect {
	case "postgres":
		ql, err = buildPostgresUpdateSql(table)
	case "mysql":
		ql, err = buildMysqlUpdateSql(table)
	case "oracle":
		ql, err = buildOracleUpdateSql(table)
	default:
		err = errors.New("build sql failed, unsupported dialect")
	}
	return
}

func buildPostgresUpdateSql(table *def.Interface) (ql string, err error) {
	bb := bytes.NewBuffer([]byte{})
	if table.Schema != "" {
		bb.WriteString(fmt.Sprintf(`UPDATE "%s"."%s" SET`, table.Schema, table.Name))
	} else {
		bb.WriteString(fmt.Sprintf(`UPDATE "%s" SET`, table.Name))
	}
	i := 1
	for _, col := range table.Columns {
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
		i++
	}
	bb.WriteString(` WHERE `)
	for pi, pk := range table.Pks {
		if pi > 0 {
			bb.WriteString(` AND `)
		}
		bb.WriteString(fmt.Sprintf(`"%s" = $%d`, strings.TrimSpace(pk.Name), i))
		i++
	}
	if table.Version.MapName != "" {
		bb.WriteString(fmt.Sprintf(` AND "%s" = $%d `, table.Version.Name, i))
	}
	ql = strings.TrimSpace(strings.ToUpper(bb.String()))
	return
}

func buildMysqlUpdateSql(table *def.Interface) (ql string, err error) {

	return
}

func buildOracleUpdateSql(table *def.Interface) (ql string, err error) {

	return
}
