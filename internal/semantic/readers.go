package semantic

import (
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func (l *lexer) readOperator() token {
	startByte := l.position
	charPos := utils.Utf16Len(l.input[l.lineStart:l.position]) // Calc column in UTF-16

	l.readRune()

	// Measure the content we just read in UTF-16 units
	sectionText := l.input[startByte:l.position]

	return token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    utils.Utf16Len(sectionText),
		tokenType: string(protocol.SemanticTokenTypeOperator),
		text:      sectionText,
	}
}

func (l *lexer) readUntil(delimiters map[rune]struct{}, tokenType string) token {
	startByte := l.position
	charPos := utils.Utf16Len(l.input[l.lineStart:l.position]) // Calc column in UTF-16

	// Append default delimiters
	delimiters['\n'] = struct{}{}
	delimiters['\\'] = struct{}{}
	delimiters[0] = struct{}{}

	inList := func(ch rune) bool {
		_, ok := delimiters[ch]
		return ok
	}

	for !inList(l.ch) {
		l.readRune()
	}

	// Measure the content we just read in UTF-16 units
	text := l.input[startByte:l.position]

	return token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    utils.Utf16Len(text),
		tokenType: tokenType,
		text:      text,
	}
}

func (l *lexer) readAssignment() {
	startByte := l.position
	charPos := utils.Utf16Len(l.input[l.lineStart:l.position]) // Calc column in UTF-16

	for l.ch != '=' && l.ch != '\n' && l.ch != '\\' && l.ch != 0 {
		l.readRune()
	}

	// Measure the content we just read in UTF-16 units
	propText := l.input[startByte:l.position]

	// If we hit a \n, it means it was a value line, if we hit '=' it is property
	if l.ch == '\n' || l.ch == '\\' || l.ch == 0 {
		l.queue = append(l.queue, token{
			line:      l.lineNumber,
			charPos:   charPos,
			length:    utils.Utf16Len(propText),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      propText,
		})
		return
	}

	// We hit '=' a sign, put property to queue, then analyze line further
	l.queue = append(l.queue, token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    utils.Utf16Len(propText),
		tokenType: string(protocol.SemanticTokenTypeKeyword),
		text:      propText,
	})

	if l.ch == '=' {
		l.queue = append(l.queue, l.readOperator())

		// Check if we had any special parse for property, if not just read value
		if fn, ok := specialFunctionMap[propText]; ok {
			fn(l)
		} else {
			l.queue = append(l.queue, l.readUntil(map[rune]struct{}{}, string(protocol.SemanticTokenTypeString)))
		}
	}
}

func (l *lexer) readSection() token {
	startByte := l.position
	charPos := utils.Utf16Len(l.input[l.lineStart:l.position]) // Calc column in UTF-16

	for l.ch != ']' && l.ch != '\n' && l.ch != 0 {
		l.readRune()
	}
	if l.ch == ']' {
		l.readRune()
	}

	// Measure the content we just read in UTF-16 units
	sectionText := l.input[startByte:l.position]

	return token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    utils.Utf16Len(sectionText),
		tokenType: string(protocol.SemanticTokenTypeNamespace),
		text:      sectionText,
	}
}

func (l *lexer) readComment() token {
	startByte := l.position
	charPos := utils.Utf16Len(l.input[l.lineStart:l.position]) // Calc column in UTF-16

	for l.ch != '\n' && l.ch != 0 {
		l.readRune()
	}

	// Measure the content we just read in UTF-16 units
	commentText := l.input[startByte:l.position]

	return token{
		line:      l.lineNumber,
		charPos:   charPos,
		length:    utils.Utf16Len(commentText),
		tokenType: string(protocol.SemanticTokenTypeComment),
		text:      commentText,
	}
}
