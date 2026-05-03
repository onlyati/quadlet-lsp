package parser

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/pkg/quadlet/lexer"
)

type NodePosition = lexer.TokenPosition

type Node interface {
	String() string
}

type FindTokenOutput struct {
	CurrentNode Node
	ParentNodes []Node
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

	isOverFunc := func(position, prevEndPos, currStartPos NodePosition) bool {
		// This scerio that previous end's line is over but not reach the
		// current start's line, so between them.
		if currStartPos.LineNumber > position.LineNumber {
			return true
		}

		// This scenario we are the line of previous end's line, position must
		// be checked.
		if position.LineNumber == prevEndPos.LineNumber && position.Position > prevEndPos.Position {
			return true
		}

		// Similar scenario like previous but now we are line of current's line.
		if position.LineNumber == currStartPos.LineNumber && position.Position < currStartPos.Position {
			return true
		}

		return false
	}

	if q == nil {
		return FindTokenOutput{nil, nil}
	}

	// Search in document's comment
	var prevDocument *CommentNode = nil
	if q.Documents != nil {
		for _, node := range q.Documents {
			if prevDocument != nil {
				if isOverFunc(position, prevDocument.EndPos, node.StartPos) {
					return FindTokenOutput{
						CurrentNode: nil,
						ParentNodes: nil,
					}
				}
			}

			if inLineFunc(position, node.StartPos, node.EndPos) {
				return FindTokenOutput{
					CurrentNode: node,
					ParentNodes: nil,
				}
			}

			prevDocument = node
		}
	}

	var prevSection *SectionNode = nil
	var prevAssign *AssignNode = nil

	if q.Sections != nil {
		for _, section := range q.Sections {

			if prevAssign != nil && *prevAssign.Value.Value == "" {
				if isOverFunc(position, prevAssign.EndPos, section.StartPos) {
					return FindTokenOutput{
						CurrentNode: nil,
						ParentNodes: []Node{prevAssign, prevSection},
					}
				}
			}
			if prevSection != nil {
				if isOverFunc(position, prevSection.EndPos, section.StartPos) {
					return FindTokenOutput{
						CurrentNode: nil,
						ParentNodes: []Node{prevSection},
					}
				}
			}
			prevAssign = nil

			// Search in section's comment
			for _, sectionDoc := range section.Documents {
				if inLineFunc(position, sectionDoc.StartPos, sectionDoc.EndPos) {
					return FindTokenOutput{
						CurrentNode: sectionDoc,
						ParentNodes: nil,
					}
				}
			}
			if inLineFunc(position, section.StartPos, section.EndPos) {
				return FindTokenOutput{
					CurrentNode: section,
					ParentNodes: nil,
				}
			}

			// Search in assigment
			for _, assingment := range section.Assignments {
				if prevAssign != nil && *prevAssign.Value.Value == "" {
					if isOverFunc(position, prevAssign.EndPos, assingment.StartPos) {
						return FindTokenOutput{
							CurrentNode: nil,
							ParentNodes: []Node{prevAssign, section},
						}
					}
				}
				// Search in assingment's comment
				for _, assigmentDoc := range assingment.Documents {
					if inLineFunc(position, assigmentDoc.StartPos, assigmentDoc.EndPos) {
						return FindTokenOutput{
							CurrentNode: assigmentDoc,
							ParentNodes: []Node{section},
						}
					}
				}
				// Search in assigment's value
				if assingment.Value != nil {
					if inLineFunc(position, assingment.Value.StartPos, assingment.Value.EndPos) {
						return FindTokenOutput{
							CurrentNode: assingment.Value,
							ParentNodes: []Node{assingment, section},
						}
					}
					// If cursor is past the '=' but before the value ends (or at the end of an empty value)
					if position.LineNumber == assingment.Value.StartPos.LineNumber &&
						position.Position >= assingment.Value.StartPos.Position {
						return FindTokenOutput{
							CurrentNode: assingment.Value,
							ParentNodes: []Node{assingment, section},
						}
					}
				}

				// Search in assigment
				if inLineFunc(position, assingment.StartPos, assingment.EndPos) {
					return FindTokenOutput{
						CurrentNode: assingment,
						ParentNodes: []Node{section},
					}
				}

				if assingment.Value != nil {
					if isOverFunc(position, assingment.EndPos, assingment.Value.StartPos) {
						prevAssign = assingment
						break
					}
				}

				prevAssign = assingment
			}

			prevSection = section
		}
	}

	if prevAssign != nil {
		if prevAssign.Value == nil {
			return FindTokenOutput{
				CurrentNode: nil,
				ParentNodes: []Node{prevAssign, prevSection},
			}
		}
		if *prevAssign.Value.Value == "" {
			return FindTokenOutput{
				CurrentNode: nil,
				ParentNodes: []Node{prevAssign, prevSection},
			}
		}
	}

	return FindTokenOutput{
		CurrentNode: nil,
		ParentNodes: []Node{prevSection},
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

	if len(a.Documents) > 0 {
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
