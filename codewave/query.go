package codewave

import (
	"fmt"
	"github.com/pharosnet/dalg/def"
	"strings"
)

func waveQuery(w Writer, definition def.Interface) {
	for _, query := range definition.Queries {
		queryResult := strings.ToLower(strings.TrimSpace(query.Result))
		query.Result = queryResult
		if queryResult == "one" {
			waveQueryOne(w, definition, query)
		} else if queryResult == "list" || queryResult == "" {
			waveQueryList(w, definition, query)
		} else if queryResult == "int64" || query.Result == "float62" || query.Result == "string" || query.Result == "bool" {
			waveQueryBuiltin(w, definition, query)
		}
	}
}

func flatSQL(ql string) string {
	querySql := ""
	sqlLines := strings.Split(ql, "\n")
	for _, line := range sqlLines {
		querySql = querySql + " " + strings.TrimSpace(strings.Replace(line, "\t", " ", -1))
	}
	return strings.TrimSpace(querySql)
}

func waveQueryOne(w Writer, d def.Interface, q def.Query) {
	w.WriteString(fmt.Sprintf("const %s%sSQL = `%s` \n", toCamel(d.MapName, false), toCamel(q.MapName, true), flatSQL(q.Sql.Value)))
	w.WriteString("\n")
	queryArgs := ""
	queryArgsVar := ""
	for _, arg := range q.Args {
		queryArgs = queryArgs + fmt.Sprintf(`, %s %s`, toCamel(arg.MapName, false), arg.MapType)
		queryArgsVar = queryArgsVar + ", " + toCamel(arg.MapName, false)
	}
	w.WriteString(fmt.Sprintf("func %s%s(ctx context.Context%s) (row *%s, err error) { \n",
		toCamel(d.MapName, true), toCamel(q.MapName, true), queryArgs, toCamel(d.MapName, true)))
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
	w.WriteString(fmt.Sprintf(`		%s, scanErr := %sScan(rows)`, toCamel(d.MapName, false), toCamel(d.MapName, false)))
	w.WriteString("\n")
	w.WriteString(`		if scanErr != nil {`)
	w.WriteString("\n")
	w.WriteString(`			err = scanErr`)
	w.WriteString("\n")
	w.WriteString(`			return`)
	w.WriteString("\n")
	w.WriteString(`		}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		row = %s`, toCamel(d.MapName, false)))
	w.WriteString("\n")
	w.WriteString(`		return`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	if err = dalc.Query(ctx, %s%sSQL, queryFn%s); err != nil {`, toCamel(d.MapName, false), toCamel(q.MapName, true), queryArgsVar))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		err = fmt.Errorf("query: %s%s failed, %s", err)`, toCamel(d.MapName, true), toCamel(q.MapName, true), "%v"))
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	return`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

func waveQueryList(w Writer, d def.Interface, q def.Query) {
	w.WriteString(fmt.Sprintf("const %s%sSQL = `%s` \n", toCamel(d.MapName, false), toCamel(q.MapName, true), flatSQL(q.Sql.Value)))
	w.WriteString("\n")
	queryArgs := ""
	queryArgsVar := ""
	for _, arg := range q.Args {
		queryArgs = queryArgs + fmt.Sprintf(`, %s %s`, toCamel(arg.MapName, false), arg.MapType)
		queryArgsVar = queryArgsVar + ", " + toCamel(arg.MapName, false)
	}
	w.WriteString(fmt.Sprintf(`func %s%s(ctx context.Context, fn %sQueryCallbackFunc%s) (err error) {`,
		toCamel(d.MapName, true), toCamel(q.MapName, true), toCamel(d.MapName, true), queryArgs))
	w.WriteString("\n")
	w.WriteString(`	queryFn := func(ctx context.Context, rows *sql.Rows, rowErr error) (err error) {`)
	w.WriteString("\n")
	w.WriteString(`		if rowErr != nil {`)
	w.WriteString("\n")
	w.WriteString(`			err = fn(ctx, nil, rowErr)`)
	w.WriteString("\n")
	w.WriteString(`			return`)
	w.WriteString("\n")
	w.WriteString(`		}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		%s, scanErr := %sScan(rows)`, toCamel(d.MapName, false), toCamel(d.MapName, false)))
	w.WriteString("\n")
	w.WriteString(`		if scanErr != nil {`)
	w.WriteString("\n")
	w.WriteString(`			err = fn(ctx, nil, scanErr)`)
	w.WriteString("\n")
	w.WriteString(`			return`)
	w.WriteString("\n")
	w.WriteString(`		}`)
	w.WriteString("\n")
	w.WriteString(`		fn(ctx, user, nil)`)
	w.WriteString("\n")
	w.WriteString(`		return`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	if err = dalc.Query(ctx, %s%sSQL, queryFn%s); err != nil {`, toCamel(d.MapName, false), toCamel(q.MapName, true), queryArgsVar))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		err = fmt.Errorf("query: %s%s failed, %s", err)`, toCamel(d.MapName, true), toCamel(q.MapName, true), "%v"))
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	return`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

// int64 float64 string bool -> sql.nullType
func waveQueryBuiltin(w Writer, d def.Interface, q def.Query) {
	w.WriteString(fmt.Sprintf("const %s%sSQL = `%s` \n", toCamel(d.MapName, false), toCamel(q.MapName, true), flatSQL(q.Sql.Value)))
	w.WriteString("\n")
	queryArgs := ""
	queryArgsVar := ""
	for _, arg := range q.Args {
		queryArgs = queryArgs + fmt.Sprintf(`, %s %s`, toCamel(arg.MapName, false), arg.MapType)
		queryArgsVar = queryArgsVar + ", " + toCamel(arg.MapName, false)
	}

	w.WriteString(fmt.Sprintf("func %s%s(ctx context.Context%s) (result %s, err error) { \n",
		toCamel(d.MapName, true), toCamel(q.MapName, true), queryArgs, q.Result))
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

	w.WriteString(fmt.Sprintf(`		nullResult := sql.Null%s{}`, toCamel(q.Result, true)))
	w.WriteString("\n")
	w.WriteString(`		if scanErr := rows.Scan(&result); scanErr != nil {`)
	w.WriteString("\n")
	w.WriteString(`			err = scanErr`)
	w.WriteString("\n")
	w.WriteString(`			return`)
	w.WriteString("\n")
	w.WriteString(`		}`)
	w.WriteString("\n")
	w.WriteString(`		if nullResult.Valid {`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`			result = nullResult.%s`, toCamel(q.Result, true)))
	w.WriteString("\n")
	w.WriteString(`		}`)
	w.WriteString("\n")
	w.WriteString(`		return`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	if err = dalc.Query(ctx, %s%sSQL, queryFn%s); err != nil {`, toCamel(d.MapName, false), toCamel(q.MapName, true), queryArgsVar))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`		err = fmt.Errorf("query: %s%s failed, %s", err)`, toCamel(d.MapName, true), toCamel(q.MapName, true), "%v"))
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	return`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}
