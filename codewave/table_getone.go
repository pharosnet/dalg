package codewave

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"strings"
)

func waveTableGetOne(w Writer, table *def.Interface) {
	ql, qlErr := buildGetOneSql(table)
	if qlErr != nil {
		panic(qlErr)
	}
	w.WriteString(fmt.Sprintf("const %sGetOneSQL = `%s` \n", toCamel(table.MapName, false), ql))
	pkArgs := ""
	for _, pk := range table.Pks {
		pkArgs = pkArgs + fmt.Sprintf(", %s %s", toCamel(pk.MapName, false), pk.MapType)
	}
	w.WriteString(fmt.Sprintf("func %sGetOne(ctx context.Context%s) (row *%s, err error) { \n",
		toCamel(table.MapName, true), pkArgs, toCamel(table.MapName, true)))
	w.WriteString(`	queryFn := func(ctx context.Context, rows *sql.Rows, rowErr error) (err error) {`)
	w.WriteString("\n")
	w.WriteString(`		if rowErr != nil {`)
	w.WriteString("\n")
	w.WriteString(`			err = rowErr`)
	w.WriteString("\n")
	w.WriteString(`			return`)
	w.WriteString("\n")
	w.WriteString(`		}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		%s, scanErr := %sScan(rows)`, toCamel(table.MapName, false), toCamel(table.MapName, false)))
	w.WriteString("\n")
	w.WriteString(`		if scanErr != nil {`)
	w.WriteString("\n")
	w.WriteString(`			err = scanErr`)
	w.WriteString("\n")
	w.WriteString(`			return`)
	w.WriteString("\n")
	w.WriteString(`		}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		row = %s`, toCamel(table.MapName, false)))
	w.WriteString("\n")
	w.WriteString(`		return`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	if err = dalc.Query(ctx, %sGetOneSQL, queryFn, id); err != nil {`, toCamel(table.MapName, false)))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		err = fmt.Errorf("%s: get one failed, %s", err)`, toCamel(table.MapName, true), "%v"))
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	return`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

func buildGetOneSql(table *def.Interface) (ql string, err error) {
	switch table.Dialect {
	case "postgres":
		ql, err = buildPostgresGetOneSql(table)
	case "mysql":
		ql, err = buildMysqlGetOneSql(table)
	case "oracle":
		ql, err = buildOracleGetOneSql(table)
	default:
		err = errors.New("build sql failed, unsupported dialect")
	}
	return
}

func buildPostgresGetOneSql(table *def.Interface) (ql string, err error) {
	bb := bytes.NewBuffer([]byte{})
	bb.WriteString(`SELECT `)
	for i, col := range table.Columns {
		if i > 0 {
			bb.WriteString(", ")
		}
		bb.WriteString(fmt.Sprintf(`"%s"`, strings.TrimSpace(col.Name)))
	}
	if table.Schema != "" {
		bb.WriteString(fmt.Sprintf(` FROM "%s"."%s" WHERE `, table.Schema, table.Name))
	} else {
		bb.WriteString(fmt.Sprintf(` FROM "%s" WHERE `, table.Name))
	}
	for pi, pk := range table.Pks {
		if pi > 0 {
			bb.WriteString(` AND `)
		}
		bb.WriteString(fmt.Sprintf(`"%s" = $%d`, strings.TrimSpace(pk.Name), pi+1))
	}
	ql = strings.TrimSpace(strings.ToUpper(bb.String()))
	return
}

func buildMysqlGetOneSql(table *def.Interface) (ql string, err error) {

	return
}

func buildOracleGetOneSql(table *def.Interface) (ql string, err error) {

	return
}
