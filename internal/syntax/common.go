package syntax

import (
	"fmt"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

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
