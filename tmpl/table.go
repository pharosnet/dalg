package tmpl

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/pharosnet/dalg/def"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)



func WriteTableFile(tableDef def.Interface, dir string) error {
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
	}
	tableDef.PkNum = int64(len(tableDef.Pks))
	tpl, tplErr := template.New("_table").Parse(_tableTpl)
	tpl.Funcs(map[string]interface{}{
		"fup": func(s string) string {
			bb := []byte(strings.TrimSpace(s))
			a := strings.ToUpper(string(bb[0]))
			bb[0] = []byte(a)[0]
			return string(bb)
		},
		"flow": func(s string) string {
			bb := []byte(strings.TrimSpace(s))
			a := strings.ToLower(string(bb[0]))
			bb[0] = []byte(a)[0]
			return string(bb)
		},
	})
	if tplErr != nil {
		return tplErr
	}

	buffer := bytes.NewBuffer([]byte{})

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

	// extra types


	if err := tpl.Execute(buffer, tableDef); err != nil {
		return err
	}

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


var _tableTpl = `
package {{.Package}}

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	{{range .ExtraType.Packages}}
	{{.}}
	{{end}}
)

const (
	_{{.Name}}InsertSql = ` + "`{{.InsertSql}}`" + `
	_{{.Name}}UpdateSql = ` + "`{{.UpdateSql}}`" + `
	_{{.Name}}DeleteSql = ` + "`{{.DeleteSql}}`" + `
	_{{.Name}}GetOneSql = ` + "`{{.GetOneSql}}`" + `
	//
)

{{range .ExtraType.EnumInterfaces}}

func New{{fup .Id}}(v {{.MapType}}) {{fup .Id}} {
	ok := false
	switch v {
	{{ if eq "string" {{.MapType}} }}
	{{range .Options}}
	case "{{.MapValue}}":
		ok = true
	{{end}}
	{{ else if eq "byte" {{.MapType}} }}
	{{range .Options}}
	case '{{.MapValue}}':
		ok = true
	{{end}}
	{{else}}
	{{range .Options}}
	case {{.MapValue}}:
		ok = true
	{{end}}
	{{end}}
	}
	if !ok {
		panic(fmt.Errorf("dal: new {{fup .Id}} failed, value is invalid"))
	}
	return {{fup .Id}}{v, true}
}

type {{fup .Id}} struct {
	Value {{.MapType}}
	Valid bool
}

func (n {{fup .Id}}) Scan(value interface{}) error {
	if value == nil {
		n.Valid = false
		return nil
	}
	switch value.(type) {
	case {{.MapType}}:
		vv, ok := value.({{.MapType}})
		if !ok {
			return fmt.Errorf("dal: call {{.Id}}.scan() failed, value type is not {{.MapType}}")
		}
		switch vv {
		{{if eq "string" {{.MapType}}}}
		{{range .Options}}
		case {{.Value}}:
			n.Value = "{{.MapValue}}"
		{{end}}
		{{ else if eq "byte" {{.MapType}} }}
		{{range .Options}}
		case {{.Value}}:
			n.Value = '{{.MapValue}}'
		{{end}}
		{{else}}
		{{range .Options}}
		case {{.Value}}:
			n.Value = {{.MapValue}}
		{{end}}
		{{end}}
		default:
			return fmt.Errorf("dal: call {{.Id}}.scan() failed, value is out of range")
		}
		n.Valid = true
	}
	return nil
}

func (n {{fup .Id}}) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	switch n.Value {
	{{if eq "string" {{.MapType}}}}
	{{range .Options}}
	case "{{.MapValue}}":
		return true, nil
	{{end}}
	{{ else if eq "byte" {{.MapType}} }}
	{{range .Options}}
	case '{{.MapValue}}':
		return true, nil
	{{end}}
	{{else}}
	{{range .Options}}
	case {{.MapValue}}:
		return true, nil
	{{end}}
	{{end}}
	}
	return nil, fmt.Errorf("dal: call {{.Id}}.value() failed, value is out of range")
}



{{end}}

func New{{fup .MapName}}({{range $i, $v := .Columns}} {{if gt $i 0}}, {{end}} {{flow $v.MapName}} {{$v.MapType}} {{end}}) *{{fup .MapName}} {
	now := nowTime()
	{{if eq true .EnableNil}}
	return &{{fup .MapName}}{
		{{range .Columns}}
		{{fup .MapName}}: {{.MapType}}{ {{flow .MapName}}, true },
		{{end}}
	}
	{{else}}
	return &{{fup .MapName}}{
		{{range .Columns}}
		{{fup .MapName}}: {{flow .MapName}},
		{{end}}
	}
	{{end}}
}


type {{fup .MapName}} struct {
	{{range .Columns}}
	{{fup .MapName}} {{.MapType}}
	{{end}}
}

func (row *{{fup .MapName}}) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			fmt.Fprintf(s, "({{range $i, $v := .Columns}} {{if gt $i 0}}, {{end}} {{fup $v.MapName}}: %v {{end}})",
				{{range $i, $v := .Columns}} {{if gt $i 0}}, {{end}} row.{{fup $v.MapName}} {{end}})
		default:
			fmt.Fprintf(s, "&{{{range $i, $v := .Columns}} {{if gt $i 0}}, {{end}} %v {{end}}}",
				{{range $i, $v := .Columns}} {{if gt $i 0}}, {{end}} row.{{fup $v.MapName}} {{end}})
		}
	}
}

func scan{{fup .MapName}}(sa Scanable) (row *{{fup .MapName}}, err error) {
	row = &{{fup .MapName}}{}
	scanErr := sa.Scan({{range $i, $v := .Columns}} {{if gt $i 0}}, {{end}} &row.{{fup $v.MapName}} {{end}})
	if scanErr != nil {
		err = fmt.Errorf("dal: scan failed. reason: %v", scanErr)
		return
	}
	return
}

type {{fup .MapName}}RangeFn func(ctx context.Context, row *{{fup .MapName}}, err error) error

func Insert{{fup .MapName}}(ctx context.Context, rows ...*{{fup .MapName}}) (affected int64, err error) {
	if ctx == nil {
		err = errors.New("dal: insert {{fup .MapName}} failed, context is empty")
		return
	}
	if rows == nil || len(rows) == 0 {
		err = errors.New("dal: insert {{fup .MapName}} failed, row is empty")
		return
	}
	stmt, prepareErr := prepare(ctx).PrepareContext(ctx, _{{.Name}}InsertSql)
	if prepareErr != nil {
		err = fmt.Errorf("dal: insert {{fup .MapName}} failed, prepared statement failed. reason: %v", prepareErr)
		return
	}
	defer func() {
		stmtCloseErr := stmt.Close()
		if stmtCloseErr != nil {
			err = fmt.Errorf("dal: insert {{fup .MapName}} failed, close prepare statement failed. reason: %v", stmtCloseErr)
			return
		}
	}()
	for _, row := range rows {
		result, execErr :=  stmt.ExecContext(ctx, {{range $i, $v := .Columns}} {{if gt $i 0}}, {{end}} row.{{fup $v.MapName}}{{end}})
		if execErr != nil {
			err = fmt.Errorf("dal: insert {{fup .MapName}} failed, execute statement failed. reason: %v", execErr)
			return
		}
		affectedRows, affectedErr :=  result.RowsAffected()
		if affectedErr != nil {
			err = fmt.Errorf("dal: insert {{fup .MapName}} failed, get rows affected failed. reason: %v", affectedErr)
			return
		}
		if affectedRows == 0 {
			err = errors.New("dal: insert {{fup .MapName}} failed, no rows affected")
			return
		}
		affected = affected + affectedRows
		{{if eq PkNum 0}}
		{{if eq true .Pks[0].DbIncrement}}
		id, getIdErr := result.LastInsertId()
		if getIdErr != nil {
			err = fmt.Errorf("dal: insert {{fup .MapName}} failed, get last insert id failed. reason: %v", getIdErr)
			return
		}
		if id < 0 {
			err = errors.New("dal: insert {{fup .MapName}} failed, get last insert id failed. id is invalid")
			return
		}
		row.{{fup .Pks[0].MapName}} = id
		{{end}}
		{{end}}
		if hasLog() {
			logf("dal: insert {{fup .MapName}} success, sql : %s, row : %+v\n", _{{.Name}}InsertSql, row)
		}
	}
	return
}


func Update{{fup .MapName}}(ctx context.Context, rows ...*{{fup .MapName}}) (affected int64, err error) {
	if ctx == nil {
		err = errors.New("dal: update {{fup .MapName}} failed, context is empty")
		return
	}
	if rows == nil || len(rows) == 0 {
		err = errors.New("dal: update {{fup .MapName}} failed, row is empty")
		return
	}
	stmt, prepareErr := prepare(ctx).PrepareContext(ctx, _{{.Name}}UpdateSql)
	if prepareErr != nil {
		err = fmt.Errorf("dal: update {{fup .MapName}} failed, prepared statement failed. reason: %v", prepareErr)
		return
	}
	defer func() {
		stmtCloseErr := stmt.Close()
		if stmtCloseErr != nil {
			err = fmt.Errorf("dal: update {{fup .MapName}} failed, close prepare statement failed. reason: %v", stmtCloseErr)
			return
		}
	}()
	now := nowTime()
	for _, row := range rows {
		result, execErr :=  stmt.ExecContext(ctx, {{range .CommonColumns}} row.{{fup .MapName}}, {{end}} {{range .Pks}} row.{{fup .MapName}}, {{end}} row.{{fup .Version}})
		if execErr != nil {
			err = fmt.Errorf("dal: update {{fup .MapName}} failed, execute statement failed. reason: %v", execErr)
			return
		}
		affectedRows, affectedErr :=  result.RowsAffected()
		if affectedErr != nil {
			err = fmt.Errorf("dal: update {{fup .MapName}} failed, get rows affected failed. reason: %v", affectedErr)
			return
		}
		if affectedRows == 0 {
			err = errors.New("dal: update {{fup .MapName}} failed, no rows affected")
			return
		}
		affected = affected + affectedRows
		if hasLog() {
			logf("dal: update {{fup .MapName}} success, sql : %s, row : %+v\n", _{{.Name}}UpdateSql, row)
		}
		row.Version.Int64 ++
	}
	return
}

func Delete{{fup .MapName}}(ctx context.Context, rows ...*{{fup .MapName}}) (affected int64, err error) {
	if ctx == nil {
		err = errors.New("dal: delete {{fup .MapName}} failed, context is empty")
		return
	}
	if rows == nil || len(rows) == 0 {
		err = errors.New("dal: delete {{fup .MapName}} failed, row is empty")
		return
	}
	stmt, prepareErr := prepare(ctx).PrepareContext(ctx, _{{.Name}}DeleteSql)
	if prepareErr != nil {
		err = fmt.Errorf("dal: delete {{fup .MapName}} failed, prepared statement failed. reason: %v", prepareErr)
		return
	}
	defer func() {
		stmtCloseErr := stmt.Close()
		if stmtCloseErr != nil {
			err = fmt.Errorf("dal: delete {{fup .MapName}} failed, close prepare statement failed. reason: %v", stmtCloseErr)
			return
		}
	}()
	for _, row := range rows {
		result, execErr :=  stmt.ExecContext(ctx, {{range .Pks}} row.{{fup .MapName}}, {{end}} row.{{fup .Version}})
		if execErr != nil {
			err = fmt.Errorf("dal: delete {{fup .MapName}} failed, execute statement failed. reason: %v", execErr)
			return
		}
		affectedRows, affectedErr :=  result.RowsAffected()
		if affectedErr != nil {
			err = fmt.Errorf("dal: delete {{fup .MapName}} failed, get rows affected failed. reason: %v", affectedErr)
			return
		}
		if affectedRows == 0 {
			err = errors.New("dal: delete {{fup .MapName}} failed, no rows affected")
			return
		}
		affected = affected + affectedRows
		if hasLog() {
			logf("dal: delete {{fup .MapName}} success, sql : %s, row : %+v\n", _{{.Name}}DeleteSql, row)
		}
	}
	return
}

func GetOne{{fup .MapName}}(ctx context.Context {{range .Pks}} , {{flow .MapName}} {{ .MapType}} {{end}}) (row *{{fup .MapName}}, err error) {
	if ctx == nil {
		err = errors.New("dal: load {{fup .MapName}} failed, context is empty")
		return
	}
	if id == "" {
		err = errors.New("dal: load {{fup .MapName}} failed, id is empty")
		return
	}
	stmt, prepareErr := prepare(ctx).PrepareContext(ctx, _{{.Name}}GetOneSql)
	if prepareErr != nil {
		err = fmt.Errorf("dal: load {{fup .MapName}} failed, prepared statement failed. reason: %v", prepareErr)
		return
	}
	defer func() {
		stmtCloseErr := stmt.Close()
		if stmtCloseErr != nil {
			err = fmt.Errorf("dal: load {{fup .MapName}} failed, close prepare statement failed. reason: %v", stmtCloseErr)
			return
		}
	}()
	oriRow := stmt.QueryRowContext(ctx{{range .Pks}} , &{{flow .MapName}} {{end}})
	row, err = scan{{fup .MapName}}(oriRow)
	if hasLog() {
		logf("dal: load {{fup .MapName}} success, sql : %s, id : %v, row : %+v\n", _{{.Name}}GetOneSql, id, userRow)
	}
	return
}


`