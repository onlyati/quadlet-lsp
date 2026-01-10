// Package utils
//
// This package contains generic functions that are used in other packages.
package utils

import (
	"io/fs"
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

// ReturnAsPtr return with pointer of any type
func ReturnAsPtr[T any](s T) *T {
	return &s
}

// FirstCharacterToUpper makes the first character upper case
func FirstCharacterToUpper(s string) string {
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func QuadletWalkDir(root string, level int, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		tmp := strings.TrimPrefix(path, root)
		tmp = strings.TrimPrefix(tmp, string(os.PathSeparator))

		// Ignore hidden directories
		if d.IsDir() && strings.HasPrefix(tmp, ".") {
			return filepath.SkipDir
		}

		// Check how deep we are
		numberOfSlices := strings.Count(tmp, string(os.PathSeparator))
		if numberOfSlices > level {
			if d.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		// Call the original function
		return fn(path, d, err)
	})
}

// ListQuadletFiles List quadlet files from the current work directory based on extenstion
func ListQuadletFiles(ext, rootDir string, level int) ([]protocol.CompletionItem, error) {
	dirs := []protocol.CompletionItem{}

	err := QuadletWalkDir(rootDir, level, func(p string, entry fs.DirEntry, err error) error {
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
