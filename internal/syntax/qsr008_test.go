package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR008_Valid(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nAnnotation=FOO1=BAR\nAnnotation=FOO2=BAR",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Build]\nAnnotation=FOO1=BAR\nAnnotation=FOO2=BAR",
			"test.build",
		),
		NewSyntaxChecker(
			`[Container]\nAnnotation=FOO=BAR "MyVar=MyValue" 'foo=bar'\n`,
			"test2.build",
		),
		NewSyntaxChecker(
			`[Container]\nAnnotation=FOO=\n`,
			"test3.build",
		),
		NewSyntaxChecker(
			`[Container]\nAnnotation='fooVariable=barValue'\n`,
			"test4.build",
		),
		NewSyntaxChecker(
			`[Container]\nAnnotation=sample.valid-key=true\n`,
			"test5.container",
		),
		NewSyntaxChecker(
			`[Container]\nAnnotation=com.example.my-app.version=1.0.0\n`,
			"test6.container",
		),
		NewSyntaxChecker(
			`[Container]\nAnnotation="com.example.my-app/subkey=value"\n`,
			"test7.container",
		),
	}

	for _, s := range variants {
		diags := qsr008(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR008_InvalidUnfinished(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nAnnotation=FOO",
			"test1.container",
		),
	}

	for _, s := range variants {
		diags := qsr008(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr008", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Invalid format: bad delimiter usage at")
	}
}

func TestQSR008_InvalidSpaceFound(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nAnnotation=FOO = BAR",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nAnnotation=FOO =",
			"test2.container",
		),
	}

	for _, s := range variants {
		diags := qsr008(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr008", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Invalid format: bad delimiter usage at")
	}
}
