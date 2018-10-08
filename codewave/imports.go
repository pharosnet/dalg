package codewave

import (
	"fmt"
	"strings"
)

func waveImports(w Writer, imports []string) error {
	if imports == nil || len(imports) == 0 {
		return nil
	}
	importsMap := make(map[string]int)
	for _, importName := range imports {
		importsMap[importName] = 1
	}
	w.WriteString(`import (`)
	w.WriteString("\n")
	for _, importPkg := range imports {
		w.WriteString(fmt.Sprintf(`	"%s" `, strings.TrimSpace(importPkg)))
		w.WriteString("\n")
	}
	w.WriteString(`)`)
	w.WriteString("\n")
	w.WriteString("\n")
	return nil
}
