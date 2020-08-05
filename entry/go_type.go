package entry

import (
	"fmt"
	"path/filepath"
	"strings"
)

func NewGoType(v string) *GoType {
	if !strings.Contains(v, ".") {
		return &GoType{
			Package: "",
			Name:    strings.TrimSpace(v),
		}
	}
	lastPointIdx := strings.LastIndex(v, ".")
	left := v[:lastPointIdx]
	right := v[lastPointIdx+1:]
	name := fmt.Sprintf("%s.%s", filepath.Base(left), strings.TrimSpace(right))
	return &GoType{
		Package: left,
		Name:    name,
	}
}

type GoType struct {
	Package string
	Name    string
}
