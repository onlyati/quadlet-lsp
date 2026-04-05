// Package semantic contains functions, structs and constants that is related
// for semantic tokens.
package semantic

import (
	"strings"

	quadlet_lexer "github.com/onlyati/quadlet-lsp/pkg/quadlet/lexer"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var TokenLegends = []string{
	string(protocol.SemanticTokenTypeComment),   // Comment lines
	string(protocol.SemanticTokenTypeKeyword),   // Like 'Image=', 'Exec=', 'Pod='
	string(protocol.SemanticTokenTypeNamespace), // Section like '[Container]', '[Unit]'
	string(protocol.SemanticTokenTypeString),    // Value belongs to keywords
	string(protocol.SemanticTokenTypeOperator),  // Operators like '=', ':', ','
	string(protocol.SemanticTokenTypeClass),     // Used within values to highlight things
	string(protocol.SemanticTokenTypeParameter), // Used within values to highlight things
}

var LegendMap = func() map[string]uint32 {
	m := make(map[string]uint32)
	for i, t := range TokenLegends {
		m[t] = uint32(i)
	}
	return m
}()

type semanticToken struct {
	line      protocol.UInteger
	charPos   protocol.UInteger
	length    protocol.UInteger
	tokenType string
	text      string
}

func CalculateSemanticTokens(lexerTokens []quadlet_lexer.Token) (*protocol.SemanticTokens, error) {
	converter := tokenConverter{
		lexerTokens:    lexerTokens,
		index:          -1,
		semanticTokens: make([]semanticToken, 0, len(lexerTokens)),
	}
	converter.parseQuadlet()

	data := make([]uint32, 0, len(converter.semanticTokens)*5)

	var lastLine, lastChar uint32

	for _, token := range converter.semanticTokens {
		deltaLine := token.line - lastLine
		deltaStart := token.charPos
		if deltaLine == 0 {
			deltaStart = token.charPos - lastChar
		}

		typeIndex, ok := LegendMap[token.tokenType]
		if !ok {
			typeIndex = 0
		}

		data = append(data,
			deltaLine,
			deltaStart,
			token.length,
			typeIndex,
			0, // No modifiers
		)

		lastLine = token.line
		lastChar = token.charPos
	}

	return &protocol.SemanticTokens{Data: data}, nil
}

type tokenConverter struct {
	lexerTokens    []quadlet_lexer.Token
	index          int
	semanticTokens []semanticToken
}

var specialParsers = map[string]func(*tokenConverter, *quadlet_lexer.Token) []semanticToken{
	"Network":     (*tokenConverter).readNetworkValue,
	"Image":       (*tokenConverter).readImageValue,
	"Pod":         (*tokenConverter).readPodValue,
	"Label":       (*tokenConverter).readLabelValue,
	"Annotation":  (*tokenConverter).readLabelValue,
	"Environment": (*tokenConverter).readLabelValue,
	"Volume":      (*tokenConverter).readVolumeValue,
	"Secret":      (*tokenConverter).readSecretValue,
}

func (t *tokenConverter) readToken() *quadlet_lexer.Token {
	t.index++
	if t.index == len(t.lexerTokens)-1 {
		return nil
	}
	return &t.lexerTokens[t.index]
}

// parseQuadlet translate the regular lexer tokens to semantic tokens
func (t *tokenConverter) parseQuadlet() {
	token := t.readToken()
	lastKeyword := ""

	for token != nil {
		switch token.Type {
		case quadlet_lexer.TokenTypeComment:
			comment := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position,
				length:    token.EndPos.Position - token.StartPos.Position,
				tokenType: string(protocol.SemanticTokenTypeComment),
				text:      token.Text,
			}
			t.semanticTokens = append(t.semanticTokens, comment)
		case quadlet_lexer.TokenTypeSection:
			section := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position,
				length:    token.EndPos.Position - token.StartPos.Position,
				tokenType: string(protocol.SemanticTokenTypeNamespace),
				text:      token.Text,
			}
			t.semanticTokens = append(t.semanticTokens, section)
		case quadlet_lexer.TokenTypeKeyword:
			keyword := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position,
				length:    token.EndPos.Position - token.StartPos.Position,
				tokenType: string(protocol.SemanticTokenTypeKeyword),
				text:      token.Text,
			}
			t.semanticTokens = append(t.semanticTokens, keyword)
			lastKeyword = token.Text
		case quadlet_lexer.TokenTypeValue:
			if fn, ok := specialParsers[lastKeyword]; ok {
				valueTokens := fn(t, token)
				t.semanticTokens = append(t.semanticTokens, valueTokens...)
			} else {
				value := semanticToken{
					line:      token.StartPos.LineNumber,
					charPos:   token.StartPos.Position,
					length:    token.EndPos.Position - token.StartPos.Position,
					tokenType: string(protocol.SemanticTokenTypeString),
					text:      token.Text,
				}
				t.semanticTokens = append(t.semanticTokens, value)
			}
		case quadlet_lexer.TokenTypeAssign, quadlet_lexer.TokenTypeContSign:
			op := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position,
				length:    token.EndPos.Position - token.StartPos.Position,
				tokenType: string(protocol.SemanticTokenTypeOperator),
				text:      token.Text,
			}
			t.semanticTokens = append(t.semanticTokens, op)
		}

		token = t.readToken()
	}
}

// readNetworkValue is part of semantic token parsing for Network keyword.
func (t *tokenConverter) readNetworkValue(token *quadlet_lexer.Token) []semanticToken {
	tokenType := protocol.SemanticTokenTypeString
	if strings.HasSuffix(token.Text, ".network") {
		tokenType = protocol.SemanticTokenTypeParameter
	}

	return []semanticToken{{
		line:      token.StartPos.LineNumber,
		charPos:   token.StartPos.Position,
		length:    token.EndPos.Position - token.StartPos.Position,
		tokenType: string(tokenType),
		text:      token.Text,
	}}
}

// readPodValue is part of semantic token parsing for Pod keyword.
func (t *tokenConverter) readPodValue(token *quadlet_lexer.Token) []semanticToken {
	tokenType := protocol.SemanticTokenTypeString
	if strings.HasSuffix(token.Text, ".pod") {
		tokenType = protocol.SemanticTokenTypeParameter
	}

	return []semanticToken{{
		line:      token.StartPos.LineNumber,
		charPos:   token.StartPos.Position,
		length:    token.EndPos.Position - token.StartPos.Position,
		tokenType: string(tokenType),
		text:      token.Text,
	}}
}

// readImageValue is part of semantic token parsing for Image keyword.
func (t *tokenConverter) readImageValue(token *quadlet_lexer.Token) []semanticToken {
	if strings.HasSuffix(token.Text, ".image") || strings.HasSuffix(token.Text, ".build") {
		return []semanticToken{{
			line:      token.StartPos.LineNumber,
			charPos:   token.StartPos.Position,
			length:    token.EndPos.Position - token.StartPos.Position,
			tokenType: string(protocol.SemanticTokenTypeParameter),
			text:      token.Text,
		}}
	}

	tokens := []semanticToken{}

	parts := strings.Split(token.Text, "/")
	lastPos := token.StartPos.Position
	for i, part := range parts {
		switch i {
		case 0:
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    uint32(len(part)),
				tokenType: string(protocol.SemanticTokenTypeString),
				text:      part,
			})
			lastPos += uint32(len(part))
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    1,
				tokenType: string(protocol.SemanticTokenTypeOperator),
				text:      "/",
			})
			lastPos += uint32(1)
		case len(parts) - 1:
			// Last token handled due to ':' and '@' characters
			imageParts := strings.SplitN(part, ":", 2)
			lastIndex := 0
			if len(imageParts) == 2 {
				tokens = append(tokens, semanticToken{
					line:      token.StartPos.LineNumber,
					charPos:   lastPos,
					length:    uint32(len(imageParts[0])),
					tokenType: string(protocol.SemanticTokenTypeParameter),
					text:      imageParts[0],
				})
				lastPos += uint32(len(imageParts[0]))
				tokens = append(tokens, semanticToken{
					line:      token.StartPos.LineNumber,
					charPos:   lastPos,
					length:    1,
					tokenType: string(protocol.SemanticTokenTypeOperator),
					text:      ":",
				})
				lastPos += uint32(1)
				lastIndex = 1
			}

			tagParts := strings.SplitN(imageParts[lastIndex], "@", 2)
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    uint32(len(tagParts[0])),
				tokenType: string(protocol.SemanticTokenTypeParameter),
				text:      tagParts[0],
			})
			lastPos += uint32(len(tagParts[0]))
			if len(tagParts) == 2 {
				tokens = append(tokens, semanticToken{
					line:      token.StartPos.LineNumber,
					charPos:   lastPos,
					length:    1,
					tokenType: string(protocol.SemanticTokenTypeOperator),
					text:      "@",
				})
				lastPos += uint32(1)
				tokens = append(tokens, semanticToken{
					line:      token.StartPos.LineNumber,
					charPos:   lastPos,
					length:    uint32(len(tagParts[1])),
					tokenType: string(protocol.SemanticTokenTypeString),
					text:      tagParts[1],
				})
			}
		default:
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    uint32(len(part)),
				tokenType: string(protocol.SemanticTokenTypeParameter),
				text:      part,
			})
			lastPos += uint32(len(part))
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    1,
				tokenType: string(protocol.SemanticTokenTypeOperator),
				text:      "/",
			})
			lastPos += uint32(1)
		}
	}

	return tokens
}

// readSecretValue is part of semantic token parsing for Secret keyword.
func (t *tokenConverter) readSecretValue(token *quadlet_lexer.Token) []semanticToken {
	tokens := []semanticToken{}

	lastPos := token.StartPos.Position
	parts := strings.Split(token.Text, ",")

	for i, part := range parts {
		sparts := strings.SplitN(part, "=", 2)
		if len(sparts) == 1 {
			// This is most probably name of the secret
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    uint32(len(sparts[0])),
				tokenType: string(protocol.SemanticTokenTypeParameter),
				text:      sparts[0],
			})
			lastPos += uint32(len(sparts[0]))
		} else {
			// This is paramter pairs for the secret
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    uint32(len(sparts[0])),
				tokenType: string(protocol.SemanticTokenTypeString),
				text:      sparts[0],
			})
			lastPos += uint32(len(sparts[0]))
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    1,
				tokenType: string(protocol.SemanticTokenTypeOperator),
				text:      "=",
			})
			lastPos += uint32(1)
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    uint32(len(sparts[1])),
				tokenType: string(protocol.SemanticTokenTypeParameter),
				text:      sparts[1],
			})
			lastPos += uint32(len(sparts[1]))
		}

		if i != len(parts)-1 {
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    1,
				tokenType: string(protocol.SemanticTokenTypeOperator),
				text:      ",",
			})
			lastPos += uint32(1)
		}
	}

	return tokens
}

// readVolumeValue is part of semantic token parsing for Volume keyword.
func (t *tokenConverter) readVolumeValue(token *quadlet_lexer.Token) []semanticToken {
	tokens := []semanticToken{}

	lastPos := token.StartPos.Position
	parts := strings.Split(token.Text, ":")

	// The source volume
	tokens = append(tokens, semanticToken{
		line:      token.StartPos.LineNumber,
		charPos:   lastPos,
		length:    uint32(len(parts[0])),
		tokenType: string(protocol.SemanticTokenTypeParameter),
		text:      parts[0],
	})
	lastPos += uint32(len(parts[0]))
	if len(parts) > 1 {
		tokens = append(tokens, semanticToken{
			line:      token.StartPos.LineNumber,
			charPos:   lastPos,
			length:    1,
			tokenType: string(protocol.SemanticTokenTypeOperator),
			text:      ":",
		})
		lastPos += uint32(1)
	}

	if len(parts) == 2 || len(parts) == 3 {
		// This is the destination volume
		tokens = append(tokens, semanticToken{
			line:      token.StartPos.LineNumber,
			charPos:   lastPos,
			length:    uint32(len(parts[1])),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      parts[1],
		})
		lastPos += uint32(len(parts[1]))
		if len(parts) > 2 {
			tokens = append(tokens, semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    1,
				tokenType: string(protocol.SemanticTokenTypeOperator),
				text:      ":",
			})
			lastPos += uint32(1)
		}
	}

	if len(parts) == 3 {
		// Part where flags are specified
		sparts := strings.Split(parts[2], ",")
		for i, spart := range sparts {
			if i != len(sparts)-1 {
				tokens = append(tokens, semanticToken{
					line:      token.StartPos.LineNumber,
					charPos:   lastPos,
					length:    uint32(len(spart)),
					tokenType: string(protocol.SemanticTokenTypeString),
					text:      sparts[i],
				})
				lastPos += uint32(len(spart))
				tokens = append(tokens, semanticToken{
					line:      token.StartPos.LineNumber,
					charPos:   lastPos,
					length:    1,
					tokenType: string(protocol.SemanticTokenTypeOperator),
					text:      ",",
				})
				lastPos += uint32(1)
			} else {
				tokens = append(tokens, semanticToken{
					line:      token.StartPos.LineNumber,
					charPos:   lastPos,
					length:    uint32(len(spart)),
					tokenType: string(protocol.SemanticTokenTypeString),
					text:      sparts[i],
				})
				lastPos += uint32(len(spart))
			}
		}
	}

	return tokens
}

// readLabelValue is part of semantic token parsing for Label keyword.
func (t *tokenConverter) readLabelValue(token *quadlet_lexer.Token) []semanticToken {
	tokens := []semanticToken{}

	lastPos := token.StartPos.Position
	tokenBuf := strings.Builder{}
	propFound := false
	valueWritten := false

	for i, c := range token.Text {
		switch c {
		case '"', '\'':
			if propFound {
				tokenStr := tokenBuf.String()
				tokenBuf = strings.Builder{}
				op := semanticToken{
					line:      token.StartPos.LineNumber,
					charPos:   lastPos,
					length:    uint32(len(tokenStr)),
					tokenType: string(protocol.SemanticTokenTypeString),
					text:      tokenStr,
				}
				tokens = append(tokens, op)
				valueWritten = true
			}
			op := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position + uint32(i),
				length:    uint32(1),
				tokenType: string(protocol.SemanticTokenTypeString),
				text:      string(c),
			}
			tokens = append(tokens, op)
			lastPos = token.StartPos.Position + uint32(i) + 1
		case '=':
			if propFound {
				tokenBuf.WriteRune(c)
				continue
			}
			propFound = true
			tokenStr := tokenBuf.String()
			tokenBuf = strings.Builder{}

			op := semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   lastPos,
				length:    uint32(len(tokenStr)),
				tokenType: string(protocol.SemanticTokenTypeParameter),
				text:      tokenStr,
			}
			tokens = append(tokens, op)

			op = semanticToken{
				line:      token.StartPos.LineNumber,
				charPos:   token.StartPos.Position + uint32(i),
				length:    uint32(1),
				tokenType: string(protocol.SemanticTokenTypeOperator),
				text:      string(c),
			}
			tokens[len(tokens)-1].tokenType = string(protocol.SemanticTokenTypeParameter)
			tokens = append(tokens, op)
			lastPos = token.StartPos.Position + uint32(i) + 1
		default:
			tokenBuf.WriteRune(c)
		}
	}

	if !valueWritten {
		tokenStr := tokenBuf.String()
		tokenBuf = strings.Builder{}
		op := semanticToken{
			line:      token.StartPos.LineNumber,
			charPos:   lastPos + 1,
			length:    uint32(len(tokenStr)),
			tokenType: string(protocol.SemanticTokenTypeString),
			text:      tokenStr,
		}
		tokens = append(tokens, op)
	}

	return tokens
}
