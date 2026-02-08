package syntax

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR007_Valid(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO1=BAR\nEnvironment=FOO2=BAR",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Build]\nEnvironment=FOO1=BAR\nEnvironment=FOO2=BAR",
			"test.build",
		),
		NewSyntaxChecker(
			`[Container]\nEnvironment=FOO=BAR "MyVar=MyValue" 'foo=bar'\n`,
			"test2.build",
		),
		NewSyntaxChecker(
			`[Container]\nEnvironment=FOO=\n`,
			"test3.build",
		),
		NewSyntaxChecker(
			`[Container]\nEnvironment='fooVariable=barValue'\n`,
			"test4.build",
		),
	}

	for _, s := range variants {
		s.config = &utils.QuadletConfig{}
		diags := qsr007(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR007_InvalidUnfinished(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO",
			"test1.container",
		),
	}

	for _, s := range variants {
		s.config = &utils.QuadletConfig{}
		diags := qsr007(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr007", *diags[0].Source)
		assert.Equal(t, "Invalid format: bad delimiter usage at FOO", diags[0].Message)
	}
}

func TestQSR007_InvalidSpaceFound(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO = BAR",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO =",
			"test2.container",
		),
	}

	for _, s := range variants {
		s.config = &utils.QuadletConfig{}
		diags := qsr007(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr007", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Invalid format: bad delimiter usage at")
	}
}

func TestQSR007_ValidAfter560(t *testing.T) {
	variants := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nEnvironment=FOO",
			"test1.container",
		),
	}

	for _, s := range variants {
		s.config = &utils.QuadletConfig{}
		s.config.Podman = utils.BuildPodmanVersion(5, 6, 0)
		diags := qsr007(s)
		require.Len(t, diags, 0)
	}
}
