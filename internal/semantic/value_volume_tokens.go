package semantic

import (
	"slices"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// readVolumeValue parses volume value like: 'foo.volume:/etc/asd:ro,z'
func (l *lexer) readVolumeValue() {
	hostFound := true
	continueDetected := false
	for {
		l.skipInlineWhitespace()

		switch l.ch {
		case '\\':
			continueDetected = true
			l.readRune()
		case '\n':
			l.handleNewLine()
			if !continueDetected {
				return
			}
			continueDetected = false
		case 0:
			return
		case ':', ',':
			l.queue = append(l.queue, l.readOperator())
		default:
			if utils.IsLetter(l.ch) || l.ch == '/' {
				startByte := l.position
				charPos := utils.Utf16Len(l.input[l.lineStart:l.position])

				stoppers := []rune{'\n', '\\', 0, ':', ','}
				for !slices.Contains(stoppers, l.ch) {
					l.readRune()
				}

				text := l.input[startByte:l.position]

				tokenType := string(protocol.SemanticTokenTypeParameter)
				if !hostFound {
					tokenType = string(protocol.SemanticTokenTypeString)
				}
				hostFound = false

				l.queue = append(l.queue, token{
					line:      l.lineNumber,
					charPos:   charPos,
					length:    utils.Utf16Len(text),
					tokenType: tokenType,
				})
			} else {
				l.readRune() // Avoid infinite loop on unkown field
			}

		}
	}
}
