// Package format
//
// This package contains the format request related actions.
package format

import (
	"sort"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
)

// Define the preferred order of categories
var categoryPriority = map[data.FormatGroup]int{
	data.FormatGroupBase:        1,
	data.FormatGroupEnvironment: 2,
	data.FormatGroupNetwork:     3,
	data.FormatGroupStorage:     4,
	data.FormatGroupLabel:       5,
	data.FormatGroupSecret:      6,
	data.FormatGroupHealth:      7,
	data.FormatGroupOther:       8,
}

func FormatDocument(q *parser.QuadletNode) string {
	reorderAssignments(q)

	var sb strings.Builder

	for _, doc := range q.Documents {
		sb.WriteString(doc.String())
	}
	if len(q.Documents) > 0 {
		sb.WriteRune('\n')
	}

	for i, section := range q.Sections {
		if i > 0 {
			sb.WriteRune('\n')
		}

		for _, doc := range section.Documents {
			sb.WriteString(doc.String())
		}
		if section.Text != nil {
			sb.WriteString(ensureBrackets(*section.Text) + "\n")
		}

		sectionName := strings.Trim(getSafeStr(section.Text), "[] ")
		var lastGroup data.FormatGroup

		for j, assign := range section.Assignments {
			if assign == nil {
				continue
			}

			currentGroup := getGroupForAssignment(sectionName, assign.Name)

			if j > 0 && currentGroup != lastGroup {
				sb.WriteRune('\n')
			}

			sb.WriteString(assign.String())
			lastGroup = currentGroup
		}
	}

	return sb.String()
}

func reorderAssignments(q *parser.QuadletNode) {
	for _, section := range q.Sections {
		if section.Text == nil || len(section.Assignments) == 0 {
			continue
		}

		sectionName := strings.Trim(*section.Text, "[]")
		props, hasProps := data.PropertiesMap[sectionName]

		groupLookup := make(map[string]data.FormatGroup)
		if hasProps {
			for _, item := range props {
				group := item.FormatGroup
				if group == "" {
					group = data.FormatGroupOther
				}
				groupLookup[item.Label] = group
			}
		}

		sort.SliceStable(section.Assignments, func(i, j int) bool {
			a := section.Assignments[i]
			b := section.Assignments[j]

			nameA := ""
			if a.Name != nil {
				nameA = *a.Name
			}
			nameB := ""
			if b.Name != nil {
				nameB = *b.Name
			}

			catA := groupLookup[nameA]
			if catA == "" {
				catA = data.FormatGroupOther
			}
			catB := groupLookup[nameB]
			if catB == "" {
				catB = data.FormatGroupOther
			}

			if catA != catB {
				return categoryPriority[catA] < categoryPriority[catB]
			}

			// If names are different, sort by Name
			if nameA != nameB {
				return nameA < nameB
			}

			// Names are same, sort by Value
			valA := ""
			if a.Value != nil && a.Value.Value != nil {
				valA = *a.Value.Value
			}
			valB := ""
			if b.Value != nil && b.Value.Value != nil {
				valB = *b.Value.Value
			}

			return valA < valB
		})
	}
}

func getGroupForAssignment(sectionName string, name *string) data.FormatGroup {
	if name == nil {
		return data.FormatGroupOther
	}

	props, ok := data.PropertiesMap[sectionName]
	if !ok {
		return data.FormatGroupOther
	}

	for _, p := range props {
		if p.Label == *name {
			if p.FormatGroup == "" {
				return data.FormatGroupOther
			}
			return p.FormatGroup
		}
	}
	return data.FormatGroupOther
}

func ensureBrackets(h string) string {
	h = strings.TrimSpace(h)
	if !strings.HasPrefix(h, "[") {
		h = "[" + h
	}
	if !strings.HasSuffix(h, "]") {
		h = h + "]"
	}
	return h
}

func getSafeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
