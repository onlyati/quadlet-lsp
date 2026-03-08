package semantic

import (
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func ImageValueTokens(value string) []Token {
	tokens := []Token{}

	tmp := strings.Split(value, "/")

	if len(tmp) == 1 {
		tokens = append(tokens, Token{
			CharPos:   0,
			Length:    uint32(len(value)),
			TokenType: string(protocol.SemanticTokenTypeString),
		})
	} else {
		offset := uint32(0)
		for i, part := range tmp {
			tokenType := ""
			if i == 0 {
				tokenType = string(protocol.SemanticTokenTypeString)
			} else {
				tokenType = string(protocol.SemanticTokenTypeClass)
			}
			tokens = append(tokens, Token{
				CharPos:   offset,
				Length:    uint32(len(part)),
				TokenType: tokenType,
			})
			offset += uint32(len(part))

			if i != len(tmp)-1 {
				tokens = append(tokens, Token{
					CharPos:   offset,
					Length:    uint32(1),
					TokenType: string(protocol.SemanticTokenTypeOperator),
				})
			}

			offset += uint32(1)
		}
	}

	return tokens
}
