package semantic

import (
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// readEnvValue parses label and annotation values.
func (l *lexer) readLabelValue() {
	extraDelimiter := ' '
	foundEqual := false

	l.customReader(func(l *lexer) {
		switch l.ch {
		case '\'', '"':
			extraDelimiter = l.ch
			l.queue = append(l.queue, l.readOperator())
		default:
			if utils.IsLetter(l.ch) {
				delimiters := map[rune]struct{}{}
				if extraDelimiter != ' ' {
					delimiters[extraDelimiter] = struct{}{}
				}
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
