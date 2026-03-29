package lexer

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/require"
)

func Test_LexerComment(t *testing.T) {
	input := `# Line comment with 🫠 emoji
    # 日本語 comment
`
	expected := []Token{
		{
			StartPos: TokenPosition{0, 0},
			EndPos:   TokenPosition{0, 28},
			Position: 0,
			Length:   utils.Utf16Len("# Line comment with 🫠 emoji"),
			Type:     TokenTypeComment,
			Text:     "# Line comment with 🫠 emoji",
		},
		{
			StartPos: TokenPosition{1, 4},
			EndPos:   TokenPosition{1, 17},
			Position: 35,
			Length:   utils.Utf16Len("# 日本語 comment"),
			Type:     TokenTypeComment,
			Text:     "# 日本語 comment",
		},
		{
			StartPos: TokenPosition{2, 0},
			EndPos:   TokenPosition{2, 0},
			Position: 55,
			Length:   0,
			Type:     TokenTypeEOF,
			Text:     "",
		},
	}

	l := NewLexer(input)
	l.Run()

	require.Len(t, l.Tokens, len(expected), "invalid length of result instead of %d", len(expected))
	for i, tok := range l.Tokens {
		require.Equal(t, expected[i], tok, "unexpected token at %d", i)
	}
}

func Test_LexerSection(t *testing.T) {
	input := `[Unit]
    [Container]
[Foo
]
`

	expected := []Token{
		{
			StartPos: TokenPosition{0, 0},
			EndPos:   TokenPosition{0, 6},
			Position: 0,
			Length:   utils.Utf16Len("[Unit]"),
			Type:     TokenTypeSection,
			Text:     "[Unit]",
		},
		{
			StartPos: TokenPosition{1, 4},
			EndPos:   TokenPosition{1, 15},
			Position: 11,
			Length:   utils.Utf16Len("[Container]"),
			Type:     TokenTypeSection,
			Text:     "[Container]",
		},
		{
			StartPos: TokenPosition{2, 0},
			EndPos:   TokenPosition{2, 4},
			Position: 23,
			Length:   utils.Utf16Len("[Foo"),
			Type:     TokenTypeSection,
			Text:     "[Foo",
		},
		{
			StartPos: TokenPosition{3, 0},
			EndPos:   TokenPosition{3, 1},
			Position: 28,
			Length:   utils.Utf16Len("]"),
			Type:     TokenTypeSection,
			Text:     "]",
		},
		{
			StartPos: TokenPosition{4, 0},
			EndPos:   TokenPosition{4, 0},
			Position: 30,
			Length:   0,
			Type:     TokenTypeEOF,
			Text:     "",
		},
	}

	l := NewLexer(input)
	l.Run()

	require.Len(t, l.Tokens, len(expected), "invalid length of result instead of %d", len(expected))
	for i, tok := range l.Tokens {
		require.Equal(t, expected[i], tok, "unexpected token at %d", i)
	}
}

func Test_LexerKeyValuePair(t *testing.T) {
	input := `Label="env=test"
Foo1=
Foo2
Foo3 \
  = \
  app=foo \`

	expected := []Token{
		{
			StartPos: TokenPosition{0, 0},
			EndPos:   TokenPosition{0, 5},
			Position: 0,
			Length:   utils.Utf16Len("Label"),
			Type:     TokenTypeKeyword,
			Text:     "Label",
		},
		{
			StartPos: TokenPosition{0, 5},
			EndPos:   TokenPosition{0, 6},
			Position: 5,
			Length:   utils.Utf16Len("="),
			Type:     TokenTypeAssign,
			Text:     "=",
		},
		{
			StartPos: TokenPosition{0, 6},
			EndPos:   TokenPosition{0, 16},
			Position: 6,
			Length:   utils.Utf16Len("\"env=test\""),
			Type:     TokenTypeValue,
			Text:     "\"env=test\"",
		},
		{
			StartPos: TokenPosition{1, 0},
			EndPos:   TokenPosition{1, 4},
			Position: 17,
			Length:   utils.Utf16Len("Foo1"),
			Type:     TokenTypeKeyword,
			Text:     "Foo1",
		},
		{
			StartPos: TokenPosition{1, 4},
			EndPos:   TokenPosition{1, 5},
			Position: 21,
			Length:   utils.Utf16Len("="),
			Type:     TokenTypeAssign,
			Text:     "=",
		},
		{
			StartPos: TokenPosition{2, 0},
			EndPos:   TokenPosition{2, 4},
			Position: 23,
			Length:   utils.Utf16Len("Foo2"),
			Type:     TokenTypeKeyword,
			Text:     "Foo2",
		},
		{
			StartPos: TokenPosition{3, 0},
			EndPos:   TokenPosition{3, 5},
			Position: 28,
			Length:   utils.Utf16Len("Foo3 "),
			Type:     TokenTypeKeyword,
			Text:     "Foo3 ",
		},
		{
			StartPos: TokenPosition{3, 5},
			EndPos:   TokenPosition{3, 6},
			Position: 33,
			Length:   utils.Utf16Len("\\"),
			Type:     TokenTypeContSign,
			Text:     "\\",
		},
		{
			StartPos: TokenPosition{4, 2},
			EndPos:   TokenPosition{4, 3},
			Position: 37,
			Length:   utils.Utf16Len("="),
			Type:     TokenTypeAssign,
			Text:     "=",
		},
		{
			StartPos: TokenPosition{4, 4},
			EndPos:   TokenPosition{4, 5},
			Position: 39,
			Length:   utils.Utf16Len("\\"),
			Type:     TokenTypeContSign,
			Text:     "\\",
		},
		{
			StartPos: TokenPosition{5, 2},
			EndPos:   TokenPosition{5, 10},
			Position: 43,
			Length:   utils.Utf16Len("app=foo "),
			Type:     TokenTypeValue,
			Text:     "app=foo ",
		},
		{
			StartPos: TokenPosition{5, 10},
			EndPos:   TokenPosition{5, 11},
			Position: 51,
			Length:   utils.Utf16Len("\\"),
			Type:     TokenTypeContSign,
			Text:     "\\",
		},
		{
			StartPos: TokenPosition{5, 11},
			EndPos:   TokenPosition{5, 11},
			Position: 52,
			Length:   0,
			Type:     TokenTypeEOF,
			Text:     "",
		},
	}

	l := NewLexer(input)
	l.Run()

	require.Len(t, l.Tokens, len(expected), "invalid length of result instead of %d", len(expected))
	for i, tok := range l.Tokens {
		require.Equal(t, expected[i], tok, "unexpected token at %d", i)
	}
}

func Test_LexerComplex(t *testing.T) {
	input := `# Test container

[Container]
Image=foo.image
Label= \
  env=test

[Unit]
Description=Foo container`

	expected := []Token{
		{
			StartPos: TokenPosition{0, 0},
			EndPos:   TokenPosition{0, 16},
			Position: 0,
			Length:   utils.Utf16Len("# Test container"),
			Type:     TokenTypeComment,
			Text:     "# Test container",
		},
		{
			StartPos: TokenPosition{2, 0},
			EndPos:   TokenPosition{2, 11},
			Position: 18,
			Length:   utils.Utf16Len("[Container]"),
			Type:     TokenTypeSection,
			Text:     "[Container]",
		},
		{
			StartPos: TokenPosition{3, 0},
			EndPos:   TokenPosition{3, 5},
			Position: 30,
			Length:   utils.Utf16Len("Image"),
			Type:     TokenTypeKeyword,
			Text:     "Image",
		},
		{
			StartPos: TokenPosition{3, 5},
			EndPos:   TokenPosition{3, 6},
			Position: 35,
			Length:   utils.Utf16Len("="),
			Type:     TokenTypeAssign,
			Text:     "=",
		},
		{
			StartPos: TokenPosition{3, 6},
			EndPos:   TokenPosition{3, 15},
			Position: 36,
			Length:   utils.Utf16Len("foo.image"),
			Type:     TokenTypeValue,
			Text:     "foo.image",
		},
		{
			StartPos: TokenPosition{4, 0},
			EndPos:   TokenPosition{4, 5},
			Position: 46,
			Length:   utils.Utf16Len("Label"),
			Type:     TokenTypeKeyword,
			Text:     "Label",
		},
		{
			StartPos: TokenPosition{4, 5},
			EndPos:   TokenPosition{4, 6},
			Position: 51,
			Length:   utils.Utf16Len("="),
			Type:     TokenTypeAssign,
			Text:     "=",
		},
		{
			StartPos: TokenPosition{4, 7},
			EndPos:   TokenPosition{4, 8},
			Position: 53,
			Length:   utils.Utf16Len("\\"),
			Type:     TokenTypeContSign,
			Text:     "\\",
		},
		{
			StartPos: TokenPosition{5, 2},
			EndPos:   TokenPosition{5, 10},
			Position: 57,
			Length:   utils.Utf16Len("env=test"),
			Type:     TokenTypeValue,
			Text:     "env=test",
		},
		{
			StartPos: TokenPosition{7, 0},
			EndPos:   TokenPosition{7, 6},
			Position: 67,
			Length:   utils.Utf16Len("[Unit]"),
			Type:     TokenTypeSection,
			Text:     "[Unit]",
		},
		{
			StartPos: TokenPosition{8, 0},
			EndPos:   TokenPosition{8, 11},
			Position: 74,
			Length:   utils.Utf16Len("Description"),
			Type:     TokenTypeKeyword,
			Text:     "Description",
		},
		{
			StartPos: TokenPosition{8, 11},
			EndPos:   TokenPosition{8, 12},
			Position: 85,
			Length:   utils.Utf16Len("="),
			Type:     TokenTypeAssign,
			Text:     "=",
		},
		{
			StartPos: TokenPosition{8, 12},
			EndPos:   TokenPosition{8, 25},
			Position: 86,
			Length:   utils.Utf16Len("Foo container"),
			Type:     TokenTypeValue,
			Text:     "Foo container",
		},
		{
			StartPos: TokenPosition{8, 25},
			EndPos:   TokenPosition{8, 25},
			Position: 99,
			Length:   0,
			Type:     TokenTypeEOF,
			Text:     "",
		},
	}

	l := NewLexer(input)
	l.Run()

	require.Len(t, l.Tokens, len(expected), "invalid length of result instead of %d", len(expected))
	for i, tok := range l.Tokens {
		require.Equal(t, expected[i], tok, "unexpected token at %d", i)
	}
}
