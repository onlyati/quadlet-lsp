package terminal

import (
	"strings"
)

const IndentSpaces = 2

var Indent = strings.Repeat(" ", IndentSpaces)

func IndentString(indent int) string {
	return strings.Repeat(Indent, indent)
}

func PrintIndent(indent int) {
	Print(IndentString(indent))
}
