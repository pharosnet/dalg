package codewave

import (
	"github.com/pharosnet/dalg/def"
	"github.com/pharosnet/dalg/logger"
)

func waveViews(views []def.Interface) error {
	for _, view := range views {
		if view.Name == "" {
			view.Name = toUnderScore(view.MapName)
		}
		for i, col := range view.Columns {
			pkg, mapType := parseCustomizeType(col.MapType)
			col.MapType = mapType
			if pkg != "" {
				view.Imports = append(view.Imports, pkg)
			}
			view.Columns[i] = col
		}
		if err := waveView(view); err != nil {
			logger.Log().Println(err)
			return err
		}
	}
	return nil
}

func waveView(view def.Interface) error {
	w := NewWriter()
	// intro
	waveIntroduction(w)
	// package
	wavePackage(w, view.Package)
	// imports
	waveViewImports(w, view.Imports)
	// struct
	waveViewStruct(w, view)
	// queries
	waveViewQueries(w, view)
	return WriteToFile(w, view)
}

func waveViewImports(w Writer, imports []string) {
	imports = append(imports, "context")
	imports = append(imports, "database/sql")
	imports = append(imports, "fmt")
	imports = append(imports, "github.com/pharosnet/dalc")
	waveImports(w, imports)
}

func waveViewStruct(w Writer, view def.Interface) {
	waveModel(w, view)
}

func waveViewQueries(w Writer, view def.Interface) {
	waveQuery(w, view)
}
