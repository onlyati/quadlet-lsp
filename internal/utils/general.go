// Package utils
//
// This package contains generic functions that are used in other packages.
package utils

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

// ReturnAsStringPtr I did not find better solution, probably not nice but works
func ReturnAsStringPtr(s string) *string {
	return &s
}

func FirstCharacterToUpper(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// ListQuadletFiles List quadlet files from the current work directory based on extenstion
func ListQuadletFiles(ext string) ([]protocol.CompletionItem, error) {
	dirs := []protocol.CompletionItem{}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cwd = filepath.Join(cwd, ext)
	files, err := filepath.Glob(cwd)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		chunks := strings.Split(file, string(os.PathSeparator))
		dirs = append(dirs, protocol.CompletionItem{
			Label:         chunks[len(chunks)-1],
			Documentation: "From work directory: " + cwd,
		})
	}

	return dirs, nil
}

// ConvertTemplateNameToFile Convert template name like 'web@siteA.container' to 'web@.container'
func ConvertTemplateNameToFile(s string) string {
	atSign := strings.Index(s, "@")
	dotSign := strings.LastIndex(s, ".")

	return s[:atSign] + "@" + s[dotSign:]
}
