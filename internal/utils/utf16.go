package utils

// Utf16Len returns the number of UTF-16 code units in a string.
// This is what VS Code and LSP expect for 'length' and 'charPos'.
// The simple builtin len() does not work.
func Utf16Len(s string) uint32 {
	var count uint32
	for _, r := range s {
		if r <= 0xFFFF {
			count++
		} else {
			count += 2 // Surrogate pair for characters like emoji
		}
	}
	return count
}

func IsLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
