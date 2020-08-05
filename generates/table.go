package generates

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/pharosnet/dalg/entry"
	"github.com/pharosnet/dalg/parser/commons"
)

func generateTable(packageName string, jsonTag bool, tables []*entry.Table) (fs []*GenerateFile, err error) {

	dataList := makeGenerateTableData(packageName, jsonTag, tables)
	tmpl, templateErr := template.New("TABLE_TEMPLATE").Parse(tableTemplate)
	if templateErr != nil {
		err = templateErr
		return
	}

	fs = make([]*GenerateFile, 0, 1)
	for _, data := range dataList {
		buf := bytes.NewBufferString("")
		execErr := tmpl.Execute(buf, data)
		if execErr != nil {
			err = execErr
			return
		}

		file := &GenerateFile{
			Name:    fmt.Sprintf("model.%s.go", data.RawName),
			Content: buf.Bytes(),
		}
		fs = append(fs, file)
	}

	return
}

func makeGenerateTableData(packageName string, jsonTag bool, tables []*entry.Table) (dataList []*GenerateTableData) {
	dataList = make([]*GenerateTableData, 0, 1)

	for _, table := range tables {
		data := &GenerateTableData{}
		data.Package = packageName
		if table.Schema != "" {
			data.RawName = fmt.Sprintf("%s.%s", strings.ToLower(table.Schema), strings.ToLower(table.Name))
		} else {
			data.RawName = fmt.Sprintf("%s", strings.ToLower(table.Name))
		}
		data.Imports = make(map[string]string)
		data.LowName = strings.ToLower(table.GoName[0:1]) + table.GoName[1:]
		data.Name = table.GoName

		data.GetSQL = buildTableGetSQL(table)
		data.InsertSQL = buildTableInsertSQL(table)
		data.UpdateSQL = buildTableUpdateSQL(table)
		data.DeleteSQL = buildTableDeleteSQL(table)

		data.Fields = make([]*TableField, 0, 1)
		for _, column := range table.Columns {
			isPk := false
			for _, pk := range table.PKs {
				if strings.ToUpper(pk) == strings.ToUpper(column.Name) {
					isPk = true
					break
				}
			}
			data.Fields = append(data.Fields, &TableField{
				Pk:       isPk,
				AutoIncr: column.AutoIncrement,
				Name:     column.GoName,
				Type:     column.GoType.Name,
				Tags: func(jsonTag bool) string {
					if !jsonTag {
						return ""
					}
					tag := commons.CamelToSnakeLow(column.GoName)
					return fmt.Sprintf("`json:\"%s\"`", tag)
				}(jsonTag),
			})
			if !data.HasAutoIncrId {
				if column.AutoIncrement {
					data.HasAutoIncrId = true
				}
			}
			if column.GoType.Package != "" && column.GoType.Package != "sql" && column.GoType.Package != "database/sql" && column.GoType.Package != "github.com/pharosnet/dalc" && column.GoType.Package != "context" {
				data.Imports[column.GoType.Package] = column.GoType.Package
			}
		}

		dataList = append(dataList, data)
	}

	return
}

func buildTableGetSQL(table *entry.Table) (q string) {
	buf := bytes.NewBufferString("")

	buf.WriteString("SELECT ")
	for i, column := range table.Columns {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("`%s`", column.Name))
	}
	buf.WriteString("FROM ")
	if table.Schema != "" {
		buf.WriteString(fmt.Sprintf("`%s`.`%s`", table.Schema, table.Name))
	} else {
		buf.WriteString(fmt.Sprintf("`%s`", table.Name))
	}
	buf.WriteString(" WHERE ")
	for i, pk := range table.PKs {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(fmt.Sprintf("`%s` = ?", pk))
	}

	q = buf.String()
	return
}

func buildTableInsertSQL(table *entry.Table) (q string) {

	buf := bytes.NewBufferString("")
	buf.WriteString("INSERT INTO ")
	if table.Schema != "" {
		buf.WriteString(fmt.Sprintf("`%s`.`%s`", table.Schema, table.Name))
	} else {
		buf.WriteString(fmt.Sprintf("`%s`", table.Name))
	}
	buf.WriteString(" ( ")
	x := 0
	for _, column := range table.Columns {
		if column.AutoIncrement {
			continue
		}
		if x > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("`%s`", column.Name))
		x++
	}
	buf.WriteString(") VALUES ( ")
	x = 0
	for _, column := range table.Columns {
		if column.AutoIncrement {
			continue
		}
		if x > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("?")
		x++
	}
	buf.WriteString(")")

	q = buf.String()

	return
}

func buildTableUpdateSQL(table *entry.Table) (q string) {
	buf := bytes.NewBufferString("")

	buf.WriteString("UPDATE ")
	if table.Schema != "" {
		buf.WriteString(fmt.Sprintf("`%s`.`%s`", table.Schema, table.Name))
	} else {
		buf.WriteString(fmt.Sprintf("`%s`", table.Name))
	}
	buf.WriteString(" SET ")
	x := 0
	for _, column := range table.Columns {
		isPk := false
		for _, pk := range table.PKs {
			if strings.ToUpper(pk) == strings.ToUpper(column.Name) {
				isPk = true
				break
			}
		}
		if isPk {
			continue
		}
		if x > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("`%s` = ?", column.Name))
		x++
	}
	buf.WriteString(" WHERE ")
	for i, pk := range table.PKs {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(fmt.Sprintf("`%s` = ?", pk))
	}

	q = buf.String()
	return
}

func buildTableDeleteSQL(table *entry.Table) (q string) {
	buf := bytes.NewBufferString("")

	buf.WriteString("DELETE FROM ")
	if table.Schema != "" {
		buf.WriteString(fmt.Sprintf("`%s`.`%s`", table.Schema, table.Name))
	} else {
		buf.WriteString(fmt.Sprintf("`%s`", table.Name))
	}
	buf.WriteString(" WHERE ")
	for i, pk := range table.PKs {
		if i > 0 {
			buf.WriteString(" AND ")
		}
		buf.WriteString(fmt.Sprintf("`%s` = ?", pk))
	}

	q = buf.String()
	return
}

var tableTemplate = `// DO NOT EDIT THIS FILE, IT IS GENERATED BY DALC
package {{ .Package }}

import (
    "context"
    "database/sql"
    "github.com/pharosnet/dalc/v2"
    {{ range $key, $value := .Imports }}"{{ $key }}"{{ end }}
)

const (
    {{ .LowName }}RowGetByPkSQL = "{{ .GetSQL }}"
    {{ .LowName }}RowInsertSQL  = "{{ .InsertSQL }}"
    {{ .LowName }}RowUpdateSQL  = "{{ .UpdateSQL }}"
    {{ .LowName }}RowDeleteSQL  = "{{ .DeleteSQL }}"
)

type {{ .Name }}Row struct { {{ range $key, $field := .Fields}}
    {{ $field.Name }} {{ $field.Type }} {{ $field.Tags }}{{ end }}
}

func (row *{{ .Name }}Row) scanSQLRow(rows *sql.Rows) (err error) {
    err = rows.Scan( {{ range $key, $field := .Fields}}
        &row.{{ $field.Name }},{{ end }}
    )
    return
}

func (row *{{ .Name }}Row) conventToGetArgs() (args *dalc.Args) {

    args = dalc.NewArgs() {{ range $key, $field := .Fields}}
    {{ if eq $field.Pk true }}args.Arg(row.{{ $field.Name }}){{ end }}{{ end }}

    return
}

func (row *{{ .Name }}Row) Get(ctx dalc.PreparedContext) (err error) {
    err = dalc.Query(ctx, {{ .LowName }}RowGetByPkSQL, row.conventToGetArgs(), func(ctx context.Context, rows *sql.Rows, rowErr error) (err error) {
        if rowErr != nil {
            err = rowErr
            return
        }
        err = row.scanSQLRow(rows)
        return
    })
    return
}

func (row *{{ .Name }}Row) conventToInsertArgs() (args *dalc.Args) {

    args = dalc.NewArgs() {{ range $key, $field := .Fields}}
    {{ if eq $field.AutoIncr false }}args.Arg(row.{{ $field.Name }}){{ end }}{{ end }}

    return
}

func (row *{{ .Name }}Row) Insert(ctx dalc.PreparedContext) (err error) {
    {{ if eq .HasAutoIncrId true }}
    insertId, execErr := dalc.ExecuteReturnInsertId(ctx, {{ .LowName }}RowInsertSQL, row.conventToInsertArgs())
    if execErr != nil {
        err = execErr
        return
    }
    row.Id = insertId
    {{ else }}
    _, execErr := dalc.Execute(ctx, {{ .LowName }}RowInsertSQL, row.conventToInsertArgs())
    if execErr != nil {
        err = execErr
        return
    }
    {{ end }}
    return
}

func (row *{{ .Name }}Row) conventToUpdateArgs() (args *dalc.Args) {

    args = dalc.NewArgs() {{ range $key, $field := .Fields}}
    {{ if eq $field.Pk false }}args.Arg(row.{{ $field.Name }}){{ end }}{{ end }}
    {{ range $key, $field := .Fields}}
        {{ if eq $field.Pk true }}args.Arg(row.{{ $field.Name }}){{ end }}{{ end }}
    return
}

func (row *{{ .Name }}Row) Update(ctx dalc.PreparedContext) (err error) {
    _, execErr := dalc.Execute(ctx, {{ .LowName }}RowUpdateSQL, row.conventToUpdateArgs())
    if execErr != nil {
        err = execErr
        return
    }
    return
}

func (row *{{ .Name }}Row) conventToDeleteArgs() (args *dalc.Args) {

    args = dalc.NewArgs() {{ range $key, $field := .Fields}}
    {{ if eq $field.Pk true }}args.Arg(row.{{ $field.Name }}){{ end }}{{ end }}
    return
}

func (row *{{ .Name }}Row) Delete(ctx dalc.PreparedContext) (err error) {
    _, execErr := dalc.Execute(ctx, {{ .LowName }}RowDeleteSQL, row.conventToDeleteArgs())
    if execErr != nil {
        err = execErr
        return
    }
    return
}


`
