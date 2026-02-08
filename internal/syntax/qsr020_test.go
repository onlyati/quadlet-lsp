package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR020_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nContainerName=foo\n",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Volume]\nVolumeName=foo\n",
			"test1.volume",
		),
		NewSyntaxChecker(
			"[Network]\nNetworkName=foo\n",
			"test1.network",
		),
		NewSyntaxChecker(
			"[Pod]\nPodName=foo\n",
			"test1.pod",
		),
	}

	for _, s := range cases {
		diags := qsr020(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR020_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nContainerName=.foo\n",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Volume]\nVolumeName=_foo\n",
			"test1.volume",
		),
		NewSyntaxChecker(
			"[Network]\nNetworkName=-foo\n",
			"test1.network",
		),
		NewSyntaxChecker(
			"[Pod]\nPodName=*foo\n",
			"test1.pod",
		),
	}

	for _, s := range cases {
		diags := qsr020(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr020", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Invalid name of unit: ")
	}
}
