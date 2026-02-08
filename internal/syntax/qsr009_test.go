package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR009_Valid(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nLabel=FOO1=BAR\nLabel=FOO2=BAR",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Build]\nLabel=FOO1=BAR\nLabel=FOO2=BAR",
			"test.build",
		),
		NewSyntaxChecker(
			`[Container]\nLabel=FOO=BAR "MyVar=MyValue" 'foo=bar'\n`,
			"test2.build",
		),
		NewSyntaxChecker(
			`[Container]\nLabel=FOO=\n`,
			"test3.build",
		),
		NewSyntaxChecker(
			`[Container]\nLabel='fooVariable=barValue'\n`,
			"test4.build",
		),
		NewSyntaxChecker(
			`[Container]\nLabel=sample.valid-key=true\n`,
			"test5.container",
		),
		NewSyntaxChecker(
			`[Container]\nLabel=com.example.my-app.version=1.0.0\n`,
			"test6.container",
		),
		NewSyntaxChecker(
			`[Container]\nLabel="com.example.my-app/subkey=value"\n`,
			"test7.container",
		),
	}

	for _, s := range variants {
		diags := qsr009(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR009_InvalidUnfinished(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nLabel=FOO",
			"test1.container",
		),
	}

	for _, s := range variants {
		diags := qsr009(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr009", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Invalid format: bad delimiter usage at")
	}
}

func TestQSR009_InvalidSpaceFound(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nLabel=FOO = BAR",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nLabel=FOO =",
			"test2.container",
		),
	}

	for _, s := range variants {
		diags := qsr009(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr009", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Invalid format: bad delimiter usage at")
	}
}
