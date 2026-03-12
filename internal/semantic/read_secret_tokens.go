package semantic

import (
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// readSecretValue parses values like 'gitea-db-password,type=env,target=GITEA__database__PASSWD'
func (l *lexer) readSecretValue() {
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
		case ',', '=':
			l.queue = append(l.queue, l.readOperator())
		case 0:
			return
		default:
			if utils.IsLetter(l.ch) || l.ch == '/' {
				delimiters := map[rune]struct{}{
					',': {},
					'=': {},
				}
				token := l.readUntil(delimiters, string(protocol.SemanticTokenTypeString))
				if l.ch == ',' || l.ch == '\n' || l.ch == '\\' || l.ch == 0 {
					token.tokenType = string(protocol.SemanticTokenTypeParameter)
				}
				l.queue = append(l.queue, token)
			} else {
				l.readRune() // Avoid infinite loop on unkown field
			}

		}
	}
}
