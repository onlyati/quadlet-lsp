package semantic

import (
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
				delimiters := map[rune]struct{}{
					':': {},
					',': {},
				}
				token := l.readUntil(delimiters, string(protocol.SemanticTokenTypeParameter))

				if !hostFound {
					token.tokenType = string(protocol.SemanticTokenTypeString)
				}
				hostFound = false

				l.queue = append(l.queue, token)
			} else {
				l.readRune() // Avoid infinite loop on unkown field
			}

		}
	}
}
