package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR015_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=foo.volume:/app",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=foo.volume:/app:ro",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=foo.volume:/app:ro,Z",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%h:%h:ro,Z",
			"test4.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%t:%t",
			"test5.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%t:/container/%t",
			"test6.container",
		),
		NewSyntaxChecker(
			"[Container]\nVolume=%h:%h/container",
			"test7.container",
		),
	}

	for _, s := range cases {
		diags := qsr015(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR015_InvalidContainerDirectory(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nVolume=foo.volume:data/config",
		"test1.container",
	)

	diags := qsr015(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr015", *diags[0].Source)
	assert.Equal(t, "Invalid format of Volume specification: container directory is not absolute", diags[0].Message)
}

func TestQSR016_UnkownFlag(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nVolume=foo.volume:/app/data/config:rw,Z,foo,nocopy",
		"test1.container",
	)

	diags := qsr015(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr015", *diags[0].Source)
	assert.Equal(t, "Invalid format of Volume specification: 'foo' flag is unknown", diags[0].Message)
}
