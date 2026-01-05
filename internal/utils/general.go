// Package utils
//
// This package contains generic functions that are used in other packages.
package utils

import (
	"io/fs"
	"path/filepath"
	"strings"
	"unicode"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

// ReturnAsStringPtr I did not find better solution, probably not nice but works
func ReturnAsStringPtr(s string) *string {
	return &s
}

// Return pointer of any
func ReturnAsPtr[T any](s T) *T {
	return &s
}

func FirstCharacterToUpper(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// ListQuadletFiles List quadlet files from the current work directory based on extenstion
func ListQuadletFiles(ext, rootDir string) ([]protocol.CompletionItem, error) {
	dirs := []protocol.CompletionItem{}

	err := filepath.WalkDir(rootDir, func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(entry.Name(), "."+ext) {
			return nil
		}

		valueKind := protocol.CompletionItemKindValue
		dirs = append(dirs, protocol.CompletionItem{
			Label:         entry.Name(),
			Documentation: "From work directory: " + p,
			Kind:          &valueKind,
		})

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dirs, nil
}

// ConvertTemplateNameToFile Convert template name like 'web@siteA.container' to 'web@.container'
func ConvertTemplateNameToFile(s string) string {
	atSign := strings.Index(s, "@")
	dotSign := strings.LastIndex(s, ".")

	return s[:atSign] + "@" + s[dotSign:]
}
