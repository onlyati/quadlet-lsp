package parser

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/pkg/quadlet/lexer"
)

type NodePosition lexer.TokenPosition

type Node interface {
	String() string
}

type FindTokenOutput struct {
	Node   Node
	RelPos NodePosition // Relative cursor position within found token
}

// QuadletNode represents the whole Quadlet file.
type QuadletNode struct {
	Documents []*CommentNode
	Sections  []*SectionNode
}

func (q *QuadletNode) String() string {
	strBuilder := strings.Builder{}

	for _, doc := range q.Documents {
		strBuilder.WriteString(doc.String())
	}

	for _, section := range q.Sections {
		strBuilder.WriteRune('\n')
		strBuilder.WriteString(section.String())
	}

	return strBuilder.String()
}

// FindToken method receives two parameter a line and a position number, then
// searches in the Quadlet and return with the token where the position is.
func (q *QuadletNode) FindToken(position NodePosition) FindTokenOutput {
	inLineFunc := func(position, startPos, endPos NodePosition) bool {
		betweenLine := position.LineNumber >= startPos.LineNumber && position.LineNumber <= endPos.LineNumber
		if !betweenLine {
			return false
		}

		if position.LineNumber == startPos.LineNumber && position.Position < startPos.Position {
			return false
		}

		if position.LineNumber == endPos.LineNumber && position.Position > endPos.Position {
			return false
		}
		return true
	}

	// Search in document's comment
	for _, node := range q.Documents {
		if inLineFunc(position, node.StartPos, node.EndPos) {
			return FindTokenOutput{
				Node: node,
				RelPos: NodePosition{
					LineNumber: position.LineNumber - node.StartPos.LineNumber,
					Position:   position.Position - node.StartPos.Position,
				},
			}
		}
	}

	for _, section := range q.Sections {
		// Search in section's comment
		for _, sectionDoc := range section.Documents {
			if inLineFunc(position, sectionDoc.StartPos, sectionDoc.EndPos) {
				return FindTokenOutput{
					Node: sectionDoc,
					RelPos: NodePosition{
						LineNumber: position.LineNumber - sectionDoc.StartPos.LineNumber,
						Position:   position.Position - sectionDoc.StartPos.Position,
					},
				}
			}
		}
		if inLineFunc(position, section.StartPos, section.EndPos) {
			return FindTokenOutput{
				Node: section,
				RelPos: NodePosition{
					LineNumber: position.LineNumber - section.StartPos.LineNumber,
					Position:   position.Position - section.StartPos.Position,
				},
			}
		}

		// Search in assigment
		for _, assigments := range section.Assignments {
			// Search in assingment's comment
			for _, assigmentDoc := range assigments.Documents {
				if inLineFunc(position, assigmentDoc.StartPos, assigmentDoc.EndPos) {
					return FindTokenOutput{
						Node: assigmentDoc,
						RelPos: NodePosition{
							LineNumber: position.LineNumber - assigmentDoc.StartPos.LineNumber,
							Position:   position.Position - assigmentDoc.StartPos.Position,
						},
					}
				}
			}
			// Search in assigment's value
			if inLineFunc(position, assigments.Value.StartPos, assigments.Value.EndPos) {
				return FindTokenOutput{
					Node: assigments.Value,
					RelPos: NodePosition{
						LineNumber: position.LineNumber - assigments.Value.StartPos.LineNumber,
						Position:   position.Position - assigments.Value.StartPos.Position,
					},
				}
			}

			// Search in assigment
			if inLineFunc(position, assigments.StartPos, assigments.EndPos) {
				return FindTokenOutput{
					Node: assigments,
					RelPos: NodePosition{
						LineNumber: position.LineNumber - assigments.StartPos.LineNumber,
						Position:   position.Position - assigments.StartPos.Position,
					},
				}
			}
		}
	}

	return FindTokenOutput{
		Node: nil,
	}
}

// CommentNode represent the comment nodes in the file.
type CommentNode struct {
	StartPos NodePosition
	EndPos   NodePosition
	Text     *string
}

func (d *CommentNode) String() string {
	strBuilder := strings.Builder{}

	if d.Text != nil {
		if !strings.HasPrefix(*d.Text, "#") {
			strBuilder.WriteString("# ")
		}
		strBuilder.WriteString(*d.Text)
		strBuilder.WriteRune('\n')
	}

	return strBuilder.String()
}

// SectionNode represent the sections, like '[Unit]', in the file.
type SectionNode struct {
	StartPos    NodePosition
	EndPos      NodePosition
	Text        *string
	Assignments []*AssignNode
	Documents   []*CommentNode
}

func (s *SectionNode) String() string {
	strBuilder := strings.Builder{}

	if s.Documents != nil {
		for _, doc := range s.Documents {
			strBuilder.WriteString(doc.String())
		}
	}

	if s.Text != nil {
		if !strings.HasPrefix(*s.Text, "[") {
			strBuilder.WriteString("[")
		}
		strBuilder.WriteString(*s.Text)
		if !strings.HasSuffix(*s.Text, "]") {
			strBuilder.WriteString("]")
		}
		strBuilder.WriteRune('\n')
	}

	for _, assign := range s.Assignments {
		if assign != nil {
			strBuilder.WriteString(assign.String())
		}
	}

	return strBuilder.String()
}

// AssignNode represents a key-value pair like 'Image=foo.image', altogether
// with its document comment.
type AssignNode struct {
	StartPos  NodePosition
	EndPos    NodePosition
	Name      *string
	Value     *ValueNode
	Documents []*CommentNode
}

func (a *AssignNode) String() string {
	strBuilder := strings.Builder{}

	if a.Documents != nil {
		strBuilder.WriteRune('\n')
		for _, doc := range a.Documents {
			if doc != nil {
				strBuilder.WriteString(doc.String())
			}
		}
	}

	if a.Name != nil {
		strBuilder.WriteString(*a.Name)
		strBuilder.WriteString("=")
	}
	if a.Value != nil {
		valueText := a.Value.String()
		if len(valueText)+len(*a.Name)+1 <= 80 {
			strBuilder.WriteString(a.Value.String())
		} else {
			// This is a long line, split onto multiple ones by word boundary
			newValueBuilder := strings.Builder{}

			linePos := len(*a.Name) // Keyword itself plus the '=' sign
			for i := range valueText {
				newValueBuilder.WriteByte(valueText[i])
				if i < len(valueText)-2 {
					if linePos > 80 && valueText[i] == ' ' && valueText[i+1] != ' ' {
						// Put a continuation sign, new line and start the next line with 2 character tab
						newValueBuilder.WriteString("\\\n  ")
						linePos = 2
						continue
					}
				}

				linePos++
			}

			strBuilder.WriteString(newValueBuilder.String())
		}
	}
	strBuilder.WriteRune('\n')

	return strBuilder.String()
}

type ValueNode struct {
	StartPos NodePosition
	EndPos   NodePosition
	Value    *string
}

func (v *ValueNode) String() string {
	strBuilder := strings.Builder{}

	if v.Value != nil {
		strBuilder.WriteString(*v.Value)
	}

	return strBuilder.String()
}
