// Package format
//
// This package contains the format request related actions.
package format

import (
	"sort"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
)

type documentSchema struct {
	header   []string
	sections map[string]map[data.FormatGroup][]sectionElement
}

type sectionElement struct {
	property string
	value    string
}

type sectionElements []sectionElement

func (a sectionElements) Len() int      { return len(a) }
func (a sectionElements) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a sectionElements) Less(i, j int) bool {
	if a[i].property == a[j].property {
		return a[i].value < a[j].value
	} else {
		return a[i].property < a[j].property
	}
}

func FormatDocument(text string) string {
	inSection := ""
	lastLine := ""
	var lastFormatGroup data.FormatGroup
	document := documentSchema{}
	document.header = make([]string, 0)
	document.sections = make(map[string]map[data.FormatGroup][]sectionElement, 0)

	// Read the whole file
	for line := range strings.SplitSeq(text, "\n") {
		line := strings.TrimSpace(line)

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := strings.TrimPrefix(line, "[")
			section = strings.TrimSuffix(section, "]")
			inSection = section
			document.sections[inSection] = make(map[data.FormatGroup][]sectionElement)
			lastLine = ""
			continue
		}

		if inSection == "" {
			// Keep the first lines as is
			document.header = append(document.header, line)
		}

		if inSection != "" {
			// Parse the lines to the map except comment and empty lines
			if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
				continue
			}

			if strings.HasSuffix(lastLine, " \\") {
				// If this is a continuation then just append to the last item's value
				actualValue := document.sections[inSection][lastFormatGroup][len(document.sections[inSection][lastFormatGroup])-1].value
				actualValue = strings.TrimSuffix(actualValue, " \\")
				actualValue = strings.TrimSpace(actualValue)

				document.sections[inSection][lastFormatGroup][len(document.sections[inSection][lastFormatGroup])-1].value = actualValue + " " + line
			} else {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) != 2 {
					continue
				}
				// If this is a new line then just parse it
				formatterGroup := data.FormatGroupOther
				for _, p := range data.PropertiesMap[inSection] {
					if p.Label == parts[0] {
						if p.FormatGroup != "" {
							formatterGroup = p.FormatGroup
						}
						break
					}
				}

				lastFormatGroup = formatterGroup
				if len(document.sections[inSection][formatterGroup]) > 0 {
					document.sections[inSection][formatterGroup] = append(document.sections[inSection][formatterGroup], sectionElement{
						property: parts[0],
						value:    parts[1],
					})
				} else {
					document.sections[inSection][formatterGroup] = []sectionElement{
						{
							property: parts[0],
							value:    parts[1],
						},
					}
				}
			}
		}

		lastLine = line
	}

	// Sort the arrays
	for k, v := range document.sections {
		if k == "Install" || k == "Unit" || k == "Service" {
			// Remained untouced
			continue
		}

		for _, vv := range v {
			sort.Sort(sectionElements(vv))
		}
	}

	// Generate the new text
	newText := ""

	for _, l := range document.header {
		newText += l + "\n"
	}

	sectionSeq := []string{
		"Unit",
		"Image",
		"Container",
		"Volume",
		"Network",
		"Kube",
		"Pod",
		"Build",
		"Service",
		"Install",
	}

	printSeq := []data.FormatGroup{
		data.FormatGroupBase,
		data.FormatGroupLabel,
		data.FormatGroupStorage,
		data.FormatGroupNetwork,
		data.FormatGroupEnvironment,
		data.FormatGroupSecret,
		data.FormatGroupHealth,
		data.FormatGroupOther,
	}

	for _, k := range sectionSeq {
		if v, ok := document.sections[k]; ok {
			newText += "[" + k + "]\n"

			for _, p := range printSeq {
				if vv, ok := v[p]; ok {
					if k != "Install" && k != "Unit" && k != "Service" {
						newText += "# " + string(p) + " options\n"
					}
					for _, element := range vv {
						newLine := element.property + "=" + strings.TrimSpace(element.value)
						if len(newLine) <= 80 {
							// Short line, one line is enough
							newText += newLine + "\n"
						} else {
							// Long line, split to multiple ones
							newText += wrapLine(newLine, 80)
						}
					}
					newText += "\n"
				}
			}

		}
	}

	return newText
}

func wrapLine(s string, width int) string {
	offset := 2 // The ' \' continuation sign
	width -= offset
	lastPossibleCutPoint := 0
	lastCutPoint := 0
	o := ""
	cutUrgent := false

	for i, c := range s {
		if i == 0 {
			continue
		}

		if c == ' ' && s[i-1] != ' ' {
			if cutUrgent {
				if offset == 2 {
					offset = 4
					width -= 2
					o = s[lastPossibleCutPoint:i] + " \\\n"
					lastCutPoint = i
				} else {
					o += " " + s[lastPossibleCutPoint:i] + " \\\n"
				}
				lastCutPoint = i
				lastPossibleCutPoint = i
				continue
			}
			lastPossibleCutPoint = i
		}

		t := (i - lastCutPoint) % width
		if t == 0 {
			if c == ' ' {
				if offset == 2 {
					offset = 4
					width -= 2
					o = s[lastCutPoint:i] + " \\\n"
					lastCutPoint = i
				} else {
					o += " " + s[lastCutPoint:i] + " \\\n"
				}
				lastCutPoint = i
			} else {
				if lastCutPoint == lastPossibleCutPoint {
					cutUrgent = true
					continue
				}
				if offset == 2 {
					offset = 4
					width -= 2
					o = s[lastCutPoint:lastPossibleCutPoint] + " \\\n"
				} else {
					o += " " + s[lastCutPoint:lastPossibleCutPoint] + " \\\n"
				}
				lastCutPoint = lastPossibleCutPoint
			}
		}
	}

	if lastCutPoint != len(s) {
		if offset == 2 {
			o = s[lastCutPoint:] + "\n"
		} else {
			o += " " + s[lastCutPoint:] + "\n"
		}
	}

	return o
}
