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
	w.WriteString(fmt.Sprintf(`type %s struct {`, toCamel(table.MapName, true)))
	w.WriteString("\n")
	for _, col := range table.Columns {
		w.WriteString(fmt.Sprintf(`	%s %s `, toCamel(col.MapName, true), col.MapType))
		w.WriteString("\n")
	}
	w.WriteString("}")
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


