package semantic

import (
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// readEnvValue parses environment values. It can be simple such as:
// `Environment="GITEA__database__USER=gitea"`.
// But it can be complex like:
// `Environment=FOO=BAR FOO2=BAR2 "MyVar=MyValue is=>here" 'foo=bar' FOO3=BAR3`.
func (l *lexer) readEnvValue() {
	extraDelimiter := ' '
	foundEqual := false

	l.customReader(func(l *lexer) {
		switch l.ch {
		case '\'', '"':
			extraDelimiter = l.ch
			l.queue = append(l.queue, l.readOperator())
			foundEqual = false
		default:
			if utils.IsLetter(l.ch) || foundEqual {
				delimiters := map[rune]struct{}{}
				delimiters[extraDelimiter] = struct{}{}
				if !foundEqual {
					delimiters['='] = struct{}{}
				}
				token := l.readUntil(delimiters, string(protocol.SemanticTokenTypeString))

				if l.ch == '=' && !foundEqual {
					foundEqual = true
					token.tokenType = string(protocol.SemanticTokenTypeParameter)
				}

				l.queue = append(l.queue, token)

				if l.ch == extraDelimiter {
					extraDelimiter = ' '
				}
				if l.ch == '=' || l.ch == extraDelimiter {
					if l.ch == ' ' {
						foundEqual = false
						l.readRune()
					} else {
						l.queue = append(l.queue, l.readOperator())
					}
				}
			} else {
				l.readRune() // Avoid infinite loop on unkown field
			}
		}
	})
}
