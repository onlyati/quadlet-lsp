package utils

import "strings"

type QuadletLine struct {
	LineNumber uint32
	Length     uint32
	Property   string
	Value      string
}

// This function scanning the passed text and
// looking for property in specific section.
func FindItems(text, section, property string) []QuadletLine {
	var findings []QuadletLine

	section = "[" + section + "]"
	inSection := false

	for i, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)

		if inSection && strings.Contains(line, "=") {
			tmp := strings.SplitN(line, "=", 2)
			if len(tmp) > 1 {
				if tmp[0] == property {
					findings = append(findings, QuadletLine{
						LineNumber: uint32(i),
						Length:     uint32(len(line)),
						Property:   tmp[0],
						Value:      tmp[1],
					})
				}
			}
			continue
		}

		if strings.HasPrefix(line, "[") && line != section {
			inSection = false
			continue
		}

		if line == section {
			inSection = true
			continue
		}
	}

	return findings
}
