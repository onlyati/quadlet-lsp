package syntax

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

// Function checking what is the extenstion in the URI
// and return if it is on the allowed array list.
// This also check for drop-ins file like foo.contaner.d/10-ports.conf.
// Return with the extension (with high capitalized first character)
// if matches on list. If does not match return value is empty.
func canFileBeApplied(uri string, allowed []string) string {
	// First check for drop-ins
	if strings.HasSuffix(uri, ".conf") {
		tmp := strings.Split(uri, string(os.PathSeparator))
		if len(tmp) > 2 {
			parentDirectory := tmp[len(tmp)-2]
			for _, item := range allowed {
				if strings.HasSuffix(parentDirectory, item+".d") {
					return "[" + utils.FirstCharacterToUpper(item) + "]"
				}
			}
		}
	}

	// Check for actual file extension like foo.container
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
