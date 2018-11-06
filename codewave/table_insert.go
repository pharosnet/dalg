package codewave

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"strings"
)

func waveTableInsert(w Writer, table def.Interface) {
	ql, qlErr := buildInsertSql(table)
	if qlErr != nil {
		panic(qlErr)
	}
	w.WriteString(fmt.Sprintf("const %sInsertSQL = `%s` \n", toCamel(table.MapName, false), ql))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`func %sInsert(ctx context.Context, rows ...*%s) (affected int64, err error) {`, toCamel(table.MapName, true), toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	fn := func(ctx context.Context, stmt *sql.Stmt, arg interface{}) (result sql.Result, err error) {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		row, ok := arg.(*%s)`, toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`		if !ok {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`			err = errors.New("%s: insert failed, invalid type")`, toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`			return`)
	w.WriteString("\n")
	w.WriteString(`		}`)
	w.WriteString("\n")
	w.WriteString(`		result, err = stmt.Exec(`)
	for _, col := range table.Columns {
		w.WriteString(fmt.Sprintf(`			&row.%s,`, toCamel(col.MapName, true)))
		w.WriteString("\n")
	}
	w.WriteString(`		)`)
	w.WriteString("\n")
	w.WriteString(`		return`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	affected, err = dalc.Execute(ctx, %sInsertSQL, fn, %sArrayMapToInterfacs(rows)...)`, toCamel(table.MapName, false), toCamel(table.MapName, false)))
	w.WriteString("\n")
	w.WriteString(`	return`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

func buildInsertSql(table def.Interface) (ql string, err error) {
	switch table.Dialect {
	case "postgres":
		ql, err = buildPostgresInsertSql(table)
	case "mysql":
		ql, err = buildMysqlInsertSql(table)
	case "oracle":
		ql, err = buildOracleInsertSql(table)
	default:
		err = errors.New("build sql failed, unsupported dialect")
	}
	return
}

func buildPostgresInsertSql(table def.Interface) (ql string, err error) {
	bb := bytes.NewBuffer([]byte{})
	if table.Schema != "" {
		bb.WriteString(fmt.Sprintf(`INSERT INTO "%s"."%s" (`, table.Schema, table.Name))
	} else {
		bb.WriteString(fmt.Sprintf(`INSERT INTO "%s" (`, table.Name))
	}
	for i, col := range table.Columns {
		if i == 0 {
			bb.WriteString(`"` + strings.TrimSpace(col.Name) + `"`)
		} else {
			bb.WriteString(`, "` + strings.TrimSpace(col.Name) + `"`)
		}
	}
	bb.WriteString(`) VALUES (`)
	colLen := len(table.Columns)
	for i := 1; i <= colLen; i++ {
		if i == 1 {
			bb.WriteString(fmt.Sprintf("$%d", i))
		} else {
			bb.WriteString(fmt.Sprintf(", $%d", i))
		}
	}
	bb.WriteString(`)`)
	ql = strings.TrimSpace(bb.String())
	return
}

func buildMysqlInsertSql(table def.Interface) (ql string, err error) {
	bb := bytes.NewBuffer([]byte{})
	bb.WriteString(fmt.Sprintf(`INSERT INTO %s (`, table.Name))
	for i, col := range table.Columns {
		if i == 0 {
			bb.WriteString(strings.TrimSpace(col.Name) )
		} else {
			bb.WriteString(`, ` + strings.TrimSpace(col.Name))
		}
	}
	bb.WriteString(`) VALUES (`)
	colLen := len(table.Columns)
	for i := 1; i <= colLen; i++ {
		if i == 1 {
			bb.WriteString("?")
		} else {
			bb.WriteString(", ?")
		}
	}
	bb.WriteString(`)`)
	ql = strings.TrimSpace(bb.String())
	return
}

func buildOracleInsertSql(table def.Interface) (ql string, err error) {
	bb := bytes.NewBuffer([]byte{})
	if table.Schema != "" {
		bb.WriteString(fmt.Sprintf(`INSERT INTO %s.%s (`, table.Schema, table.Name))
	} else {
		bb.WriteString(fmt.Sprintf(`INSERT INTO %s (`, table.Name))
	}
	for i, col := range table.Columns {
		if i == 0 {
			bb.WriteString(strings.TrimSpace(col.Name))
		} else {
			bb.WriteString(`, ` + strings.TrimSpace(col.Name))
		}
	}
	bb.WriteString(`) VALUES (`)
	colLen := len(table.Columns)
	for i := 1; i <= colLen; i++ {
		if i == 1 {
			bb.WriteString(fmt.Sprintf(":%d", i))
		} else {
			bb.WriteString(fmt.Sprintf(", :%d", i))
		}
	}
	bb.WriteString(`)`)
	ql = strings.TrimSpace(bb.String())
	return
}
