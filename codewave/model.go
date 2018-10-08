package codewave

import (
	"fmt"
	"github.com/pharosnet/dalg/def"
	"strings"
)

func waveModel(w Writer, model *def.Interface)  {
	// struct
	w.WriteString(fmt.Sprintf(`type %s struct {`, toCamel(model.MapName, true)))
	w.WriteString("\n")
	for _, col := range model.Columns {
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
	w.WriteString(fmt.Sprintf(`			%s(%s) { `, toCamel(model.MapName, true), strings.TrimSpace(strings.ToUpper(model.Name))))
	for i, col := range model.Columns {
		if i > 0 {
			w.WriteString(fmt.Sprintf(`, `))
		}
		w.WriteString(fmt.Sprintf(`%s(%s): `, toCamel(col.MapName, true), strings.TrimSpace(strings.ToUpper(col.Name))))
		w.WriteString(`%v`)
	}
	w.WriteString(` }, `)
	w.WriteString("\n")

	for _, col := range model.Columns {
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
	w.WriteString(fmt.Sprintf(`func %sScan(rows *sql.Rows) (*%s, error) {`, toCamel(model.MapName, false), toCamel(model.MapName, true)))
	w.WriteString("\n")
	w.WriteString(fmt.Sprintf(`	row := &%s{}`, toCamel(model.MapName, true)))
	w.WriteString("\n")
	w.WriteString(`	err := rows.Scan(`)
	for _, col := range model.Columns {
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
	w.WriteString(fmt.Sprintf(`type %sQueryCallbackFunc func(ctx context.Context, rows *%s, rowErr error) (err error)`, toCamel(model.MapName, true), toCamel(model.MapName, true)))
	w.WriteString("\n")
	w.WriteString("\n")
	// map to interfaces
	w.WriteString(fmt.Sprintf(`func %sArrayMapToInterfacs(rows []*%s) []interface{} {`, toCamel(model.MapName, false),toCamel(model.MapName, true)))
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