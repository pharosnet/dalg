package generates

import (
	"path/filepath"

	"github.com/pharosnet/dalg/entry"
)

func Generate(out string, jsonTag bool, tables []*entry.Table, queries []*entry.Query) (err error) {

	packageName := filepath.Base(out)

	fs := make([]*GenerateFile, 0, 1)
	fs0, tableErr := generateTable(packageName, jsonTag, tables)
	if tableErr != nil {
		err = tableErr
		return
	}
	fs = append(fs, fs0...)
	fs1, queryErr := generateQuery(packageName, jsonTag, queries)
	if queryErr != nil {
		err = queryErr
		return
	}
	fs = append(fs, fs1...)

	err = writeFiles(out, fs)

	return
}
