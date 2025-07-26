package syntax

import (
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
		return utils.FirstCharacterToUpper(ext)
	}

	return ""
}
