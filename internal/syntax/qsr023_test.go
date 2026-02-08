package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR023_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Unit]\nWants=%N-db.container",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nLabel=Unit=%n",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Pod]\nVolume=%h/app1:/app:ro",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nLabel=Unit=%",
			"test4.container",
		),
	}

	for _, s := range cases {
		diags := qsr023(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR023_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Unit]\nWants=%r-db.container",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nLabel=Unit=%ö",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Pod]\nVolume=%5/app1:/app:ro",
			"test3.pod",
		),
	}

	for _, s := range cases {
		diags := qsr023(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr023", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Specifier ")
		assert.Contains(t, diags[0].Message, "is invalid")
		assert.Equal(t, uint32(1), diags[0].Range.Start.Line)
	}
}
