package semantic

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// readNetworkValue parses network value like: 'foo.network'
func (l *lexer) readNetworkValue() {
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
		default:
			if utils.IsLetter(l.ch) || l.ch == '/' {
				token := l.readUntil(map[rune]struct{}{}, string(protocol.SemanticTokenTypeString))
				if strings.HasSuffix(strings.TrimSpace(token.text), ".network") {
					token.tokenType = string(protocol.SemanticTokenTypeParameter)
				}

				l.queue = append(l.queue, token)
			} else {
				l.readRune() // Avoid infinite loop on unkown field
			}

		}
	}
}
