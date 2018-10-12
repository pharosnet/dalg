package codewave

import (
	"github.com/pharosnet/dalg/def"
	"github.com/pharosnet/dalg/logger"
)

func waveTables(tables []def.Interface) error {
	for _, table := range tables {
		if table.Name == "" {
			table.Name = toUnderScore(table.MapName)
		}
		table.Pks = make([]def.Column, 0, 1)
		table.CommonColumns = make([]def.Column, 0, 1)
		for i, col := range table.Columns {
			pkg, mapType := parseCustomizeType(col.MapType)
			col.MapType = mapType
			if pkg != "" {
				table.Imports = append(table.Imports, pkg)
			}
			if col.Pk {
				table.Pks = append(table.Pks, col)
			} else if col.Version {
				table.Version = col
			} else {
				table.CommonColumns = append(table.CommonColumns, col)
			}
			table.Columns[i] = col
		}
		if err := waveTable(table); err != nil {
			logger.Log().Println(err)
			return err
		}
	}
	return nil
}

func waveTable(table def.Interface) error {
	w := NewWriter()
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

func waveTableImports(w Writer, imports []string) {
	imports = append(imports, "context")
	imports = append(imports, "database/sql")
	imports = append(imports, "errors")
	imports = append(imports, "fmt")
	imports = append(imports, "github.com/pharosnet/dalc")
	waveImports(w, imports)
}

func waveTableStruct(w Writer, table def.Interface) {
	waveModel(w, table)
}

func waveTableCRUD(w Writer, table def.Interface) {
	waveTableGetOne(w, table)
	waveTableInsert(w, table)
	waveTableUpdate(w, table)
	waveTableDelete(w, table)
}

func waveTableQueries(w Writer, table def.Interface) {
	waveQuery(w, table)
}
