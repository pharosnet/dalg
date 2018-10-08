package codewave

import (
	"fmt"
	"github.com/pharosnet/dalg/def"
	"github.com/pharosnet/dalg/logger"
	"strings"
)

func waveTables(tables []*def.Interface) error {
	for _, table := range tables {
		if err := waveTable(table); err != nil {
			logger.Log().Println(err)
			return err
		}
	}
	return nil
}

func waveTable(table *def.Interface) error {
	w := NewWriter()
	table.Pks = make([]def.Column, 0, 1)
	table.CommonColumns = make([]def.Column, 0, 1)
	for _, col := range table.Columns {
		if col.Pk {
			table.Pks = append(table.Pks, col)
		} else if col.Version {
			table.Version = col
		} else {
			table.CommonColumns = append(table.CommonColumns, col)
		}
		if pos := strings.LastIndexByte(col.MapType, '.'); pos > 0 {
			table.Imports = append(table.Imports, col.MapType[0:pos])
		}
	}
	table.PkNum = int64(len(table.Pks))
	// intro
	waveIntroduction(w)
	// package
	wavePackage(w, table.Package)
	// imports
	waveTableImports(w, table.Imports)
	// struct
	waveTableStruct(w, table)
	// CRUD
	waveTableCRUD(w, table)
	// queries
	waveTableQueries(w, table)
	return WriteToFile(w, table)
}

func waveTableImports(w Writer, imports []string)  {
	imports = append(imports, "context")
	imports = append(imports, "database/sql")
	imports = append(imports, "errors")
	imports = append(imports, "fmt")
	imports = append(imports, "github.com/pharosnet/dalc")
	waveImports(w, imports)
}

func waveTableStruct(w Writer, table *def.Interface)  {
	// struct
	w.WriteString(fmt.Sprintf(`type %s struct {`, toCamel(table.MapName, true)))
	w.WriteString("\n")
	for _, col := range table.Columns {
		w.WriteString(fmt.Sprintf(`	%s %s `, toCamel(col.MapName, true), col.MapType))
		w.WriteString("\n")
	}
	w.WriteString("}")
	w.WriteString("\n")
	w.WriteString("\n")
	// fmt
	w.WriteString(fmt.Sprintf(`func (row User) Format(s fmt.State, verb rune) {`))
	w.WriteString("\n")
	w.WriteString(`	switch verb {`)
	w.WriteString("\n")
	w.WriteString(`	case 'v':`)
	w.WriteString("\n")
	w.WriteString(`		fmt.Fprintf(`)
	w.WriteString("\n")
	w.WriteString(`			s,`)
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`			%s(%s) { `, toCamel(table.MapName, true), strings.TrimSpace(strings.ToUpper(table.Name))))
	for i, col := range table.Columns {
		if i > 0 {
			w.WriteString(fmt.Sprintf(`, `))
		}
		w.WriteString(fmt.Sprintf(`%s(%s): `, toCamel(col.MapName, true), strings.TrimSpace(strings.ToUpper(col.Name))))
		w.WriteString(`%v`)
	}
	w.WriteString(` }, `)
	w.WriteString("\n")

	for _, col := range table.Columns {
		w.WriteString(fmt.Sprintf(`			row.%s,`, toCamel(col.MapName, true)))
		w.WriteString("\n")
	}
	w.WriteString(`		)`)
	w.WriteString("\n")
	w.WriteString(`	)`)
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
	// scan fn
	w.WriteString(fmt.Sprintf(`func %sScan(rows *sql.Rows) (*%s, error) {`, toCamel(table.MapName, false), toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	row := &%s{}`, toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	err := rows.Scan(`)
	for _, col := range table.Columns {
		w.WriteString(fmt.Sprintf(`		&row.%s,`, toCamel(col.MapName, true)))
		w.WriteString("\n")
	}
	w.WriteString(`	)`)
	w.WriteString("\n")
	w.WriteString(`	if err != nil {`)
	w.WriteString("\n")
	w.WriteString(`		return nil, err`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	return row, nil`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
	// load fn
	w.WriteString(fmt.Sprintf(`type %sLoadOneFunc func(ctx context.Context, rows *%s, rowErr error) (err error)`, toCamel(table.MapName, true), toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString("\n")
	// map to interfaces
	w.WriteString(fmt.Sprintf(`func %sArrayMapToInterfacs(rows []*%s) []interface{} {`, toCamel(table.MapName, false),toCamel(table.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	array := make([]interface{}, len(rows))`)
	w.WriteString("\n")
	w.WriteString(`	for i, row := range rows {`)
	w.WriteString("\n")
	w.WriteString(`		array[i] = row`)
	w.WriteString("\n")
	w.WriteString(`	}`)
	w.WriteString("\n")
	w.WriteString(`	return array`)
	w.WriteString("\n")
	w.WriteString(`}`)
	w.WriteString("\n")
	w.WriteString("\n")
}

func waveTableCRUD(w Writer, table *def.Interface)  {
	waveTableGetOne(w, table)
	waveTableInsert(w, table)
	waveTableUpdate(w, table)
	waveTableDelete(w, table)
}

func waveTableQueries(w Writer, table *def.Interface)  {
	waveQuery(w, table)
}


