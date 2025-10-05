package parser

import "strings"

func isItComment(s string) bool {
	return strings.HasPrefix(s, "#") || strings.HasPrefix(s, ";")
}

func removeCommentSign(s string) string {
	s, found := strings.CutPrefix(s, "#")
	if !found {
		s, _ = strings.CutPrefix(s, ";")
	}
	s = strings.TrimSpace(s)
	return s
}

func gatherSectionName(s string) (string, bool) {
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return s[1 : len(s)-1], true
	}

	return s, false
}
