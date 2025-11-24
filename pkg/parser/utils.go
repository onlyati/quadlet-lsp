package parser

import (
	"strings"
)

func isItComment(s string) bool {
	return strings.HasPrefix(s, "#") || strings.HasPrefix(s, ";")
}

func removeCommentSign(s string) string {
	s, found := strings.CutPrefix(s, "#")
	if !found {
		s, _ = strings.CutPrefix(s, ";")
	}
	s = strings.TrimRight(s, " ")
	if len(s) > 0 {
		if s[0] == ' ' {
			s = s[1:]
		}
	}
	return s
}

func gatherSectionName(s string) (string, bool) {
	if strings.HasPrefix(s, "[") && strings.HasSuffix(s, "]") {
		return s[1 : len(s)-1], true
	}

	return s, false
}

func isDropinsBelongsToQuadlet(possibleOwner, parentDir string) bool {
	// If this is an extension dropin and match with the possible owner
	ownerExt := possibleOwner[strings.LastIndex(possibleOwner, ".")+1:]
	if v := parentDir == ownerExt+".d"; v {
		return v
	}

	// If not contains "-" then just a simple comparison
	if v := possibleOwner+".d" == parentDir; v {
		return v
	}

	// Now check for the cases when there is "-" in the parentDir name.
	// If we have "foo-bar-app.container" it can have dropins like:
	// foo-.container.d, foo-bar-.container.d or foo-bar-app.container.d
	ownerName := possibleOwner[:strings.LastIndex(possibleOwner, ".")]
	nameTags := strings.Split(ownerName, "-")
	for i := range nameTags {
		k := len(nameTags) - i
		tryDirName := strings.Join(nameTags[0:k], "-")
		if i < len(nameTags) {
			tryDirName += "-"
		}
		tryDirName += "." + ownerExt + ".d"

		if v := tryDirName == parentDir; v {
			return v
		}
	}

	return false
}
