package semantic

import (
	"regexp"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// readImageValue parses image value like: 'docker.io/gitea/gitea:rootless@sha256...'
func (l *lexer) readImageValue() {
	hostCheck := regexp.MustCompile(`(?:[a-z0-9]+(?:[a-z0-9._-]+)*\.(?:[a-z0-9]+)|localhost)`)
	hostFound := false

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
		case '/', ':':
			l.queue = append(l.queue, l.readOperator())
		case '@':
			l.queue = append(l.queue, l.readOperator())
			l.queue = append(l.queue, l.readUntil(map[rune]struct{}{}, string(protocol.SemanticTokenTypeString)))
		default:
			if utils.IsLetter(l.ch) {
				delimiters := map[rune]struct{}{
					':': {},
					'@': {},
					'/': {},
				}
				token := l.readUntil(delimiters, string(protocol.SemanticTokenTypeParameter))

				if !hostFound && hostCheck.MatchString(token.text) {
					token.tokenType = string(protocol.SemanticTokenTypeString)
					hostFound = true
				}

				l.queue = append(l.queue, token)
			} else {
				l.readRune() // Avoid inifinite loop on unkonw character
			}
		}
	}
}
