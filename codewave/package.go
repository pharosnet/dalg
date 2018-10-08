package codewave

import (
	"fmt"
)

func wavePackage(w Writer, pkg string)  {
	w.WriteString(fmt.Sprintf(`package %s `, pkg))
	w.WriteString("\n")
}
