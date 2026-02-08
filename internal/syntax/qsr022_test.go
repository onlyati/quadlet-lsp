package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR022_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/container/%n",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%t:%t",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/%N",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/%",
			"test3.container",
		),
	}

	for _, s := range cases {
		diags := qsr022(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR022_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/container/%t",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/%t",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=/%t:%t",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr022(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr022", *diags[0].Source)
		assert.Equal(t, "Specifier, %t, already starts with '/' sign", diags[0].Message)
		assert.Equal(t, uint32(1), diags[0].Range.Start.Line)
	}
}

func TestQSR022_InvalidMultipleInOneLine(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=/%t:/container/%h/test",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr022(s)
		require.Len(t, diags, 2)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr022", *diags[0].Source)
		assert.Equal(t, "Specifier, %t, already starts with '/' sign", diags[0].Message)
		assert.Equal(t, uint32(1), diags[0].Range.Start.Line)

		require.NotNil(t, diags[1].Source)
		assert.Equal(t, "quadlet-lsp.qsr022", *diags[1].Source)
		assert.Equal(t, "Specifier, %h, already starts with '/' sign", diags[1].Message)
		assert.Equal(t, uint32(1), diags[1].Range.Start.Line)

		line := "Volume=/%t:/container/%h/test"
		problemPart := line[diags[0].Range.Start.Character : diags[0].Range.End.Character+1]
		assert.Equal(t, "/%t", problemPart)

		problemPart = line[diags[1].Range.Start.Character : diags[1].Range.End.Character+1]
		assert.Equal(t, "/%h", problemPart)
	}
}
