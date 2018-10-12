package codewave

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"strings"
)

func waveTableDelete(w Writer, table def.Interface) {
	ql, qlErr := buildDeleteSql(table)
	if qlErr != nil {
		panic(qlErr)
	}
	w.WriteString(fmt.Sprintf("const %sDeleteSQL = `%s` \n", toCamel(table.MapName, false), ql))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`func %sDelete(ctx context.Context, rows ...*%s) (affected int64, err error) {`, toCamel(table.MapName, true), toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	fn := func(ctx context.Context, stmt *sql.Stmt, arg interface{}) (result sql.Result, err error) {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		row, ok := arg.(*%s)`, toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`		if !ok {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`			err = errors.New("%s: delete failed, invalid type")`, toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`			return`)
	w.WriteString("\n")
	w.WriteString(`		}`)
	w.WriteString("\n")
	w.WriteString(`		result, err = stmt.Exec(`)
	w.WriteString("\n")
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
	w.WriteString(fmt.Sprintf(`	affected, err = dalc.Execute(ctx, %sDeleteSQL, fn, %sArrayMapToInterfacs(rows)...)`, toCamel(table.MapName, false), toCamel(table.MapName, false)))
	w.WriteString("\n")
	w.WriteString(`	return`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

func buildDeleteSql(table def.Interface) (ql string, err error) {
	switch table.Dialect {
	case "postgres":
		ql, err = buildPostgresDeleteSql(table)
	case "mysql":
		ql, err = buildMysqlDeleteSql(table)
	case "oracle":
		ql, err = buildOracleDeleteSql(table)
	default:
		err = errors.New("build sql failed, unsupported dialect")
	}
	return
}

func buildPostgresDeleteSql(table def.Interface) (ql string, err error) {
	bb := bytes.NewBuffer([]byte{})
	if table.Schema != "" {
		bb.WriteString(fmt.Sprintf(`DELETE FROM "%s"."%s" WHERE `, table.Schema, table.Name))
	} else {
		bb.WriteString(fmt.Sprintf(`DELETE FROM "%s" WHERE `, table.Name))
	}
	i := 1
	for _, pk := range table.Pks {
		if i > 1 {
			bb.WriteString(" AND ")
		}
		bb.WriteString(fmt.Sprintf(`"%s" = $%d`, strings.TrimSpace(pk.Name), i))
		i++
	}
	if table.Version.MapName != "" {
		bb.WriteString(fmt.Sprintf(` AND "%s" = $%d`, strings.TrimSpace(table.Version.Name), i))
	}
	ql = strings.TrimSpace(bb.String())
	return
}

func buildMysqlDeleteSql(table def.Interface) (ql string, err error) {

	return
}

func buildOracleDeleteSql(table def.Interface) (ql string, err error) {

	return
}
