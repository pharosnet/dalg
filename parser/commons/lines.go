package commons

import (
	"bytes"
	"strings"
)

func NewLines(content string) *Lines {
	v := strings.Split(content, "\n")
	return &Lines{
		values: v,
		idx:    -1,
		size:   len(v),
	}
}

type Lines struct {
	values []string
	idx    int
	size   int
}

func (lines *Lines) NextLine() string {
	lines.idx++
	return lines.values[lines.idx]
}

func (lines *Lines) Prev() {
	lines.idx--
}

func (lines *Lines) Reset() {
	lines.idx = -1
}

func (lines *Lines) HasNext() bool {
	return lines.idx+1 < lines.size
}

func (lines *Lines) CurrentLineWords() []string {
	p := []byte(lines.values[lines.idx])
	words := ReadWords(p)
	return words
}

func (lines *Lines) Remain() string {
	buf := bytes.NewBufferString("")
	for i := lines.idx + 1; i < lines.size; i++ {
		line := strings.ReplaceAll(lines.values[i], "\t", " ")
		line = strings.TrimSpace(line)
		buf.WriteString(line)
		buf.WriteString(" ")
	}
	return buf.String()
}
