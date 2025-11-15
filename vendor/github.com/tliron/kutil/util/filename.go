package util

import (
	"regexp"
)

// On *nix: /
// On Windows: <>:"/\|?*
//   See: https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file

var fileEscapeRe = regexp.MustCompile(`[<>:"/\\|\?\*]`)

func SanitizeFilename(name string) string {
	switch name {
	case ".":
		return "_"

	case "..":
		return "__"

	default:
		return fileEscapeRe.ReplaceAllString(name, "")
	}
}
