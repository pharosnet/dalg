package main

import (
	"bytes"
	"fmt"
)

func main() {

	buffer := bytes.NewBuffer([]byte{})
	buffer.WriteString(`sss\\nxx'`)
	fmt.Println(buffer.String())
}
