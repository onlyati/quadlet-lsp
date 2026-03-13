package semantic

import (
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// readSecretValue parses values like 'gitea-db-password,type=env,target=GITEA__database__PASSWD'
func (l *lexer) readSecretValue() {
	l.customReader(func(l *lexer) {
		switch l.ch {
		case ',', '=':
			l.queue = append(l.queue, l.readOperator())
		default:
			if utils.IsLetter(l.ch) {
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
	})
}
