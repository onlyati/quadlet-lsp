package syntax

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

// regex for like:
// key=value
// "key=value"
// 'key=value'
var keyValueCheck = regexp.MustCompile(`^(['"]?)([A-Za-z_][A-Za-z0-9_]*)=(.*)(['"]?)$`)

// From version 5.6, environment variables accept simple name, like Environemnt=MYVAR
var keyValueCheck56 = regexp.MustCompile(`^(['"]?)([A-Za-z_][A-Za-z0-9_]*)(['"]?)$`)

// regex for name convention, like ContainerName, PodName, VolumeName, NetworkName
var namingConvention = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_.-]*$`)

// Function checking what is the extenstion in the URI
// and return if it is on the allowed array list.
// Return with the extension (with high capitalized first character)
// if matches on list. If does not match return value is empty.
func canFileBeApplied(uri string, allowed []string) string {
	tmp := strings.Split(uri, ".")
	ext := tmp[len(tmp)-1]
	if slices.Contains(allowed, ext) {
		return "[" + utils.FirstCharacterToUpper(ext) + "]"
	}

	return ""
}

// splitQuoted splits a string by spaces but preserves quoted (single or double) sections
func splitQuoted(input string) ([]string, error) {
	var result []string
	var current strings.Builder
	var quoteChar rune // track current quote type (' or ")
	inQuotes := false

	for _, r := range input {
		switch r {
		case '\'', '"':
			if !inQuotes {
				inQuotes = true
				quoteChar = r
			} else if r == quoteChar {
				inQuotes = false
			}
			current.WriteRune(r)
		case ' ':
			if inQuotes {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				result = append(result, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	if inQuotes {
		return nil, fmt.Errorf("unclosed quote in input: %q", input)
	}

	return result, nil
}
