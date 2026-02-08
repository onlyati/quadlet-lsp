package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR016_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nUserNS=keep-id",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nUserNS=keep-id:uid=101,gid=101",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nUserNS=auto",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nUserNS=host",
			"test4.container",
		),
		NewSyntaxChecker(
			"[Container]\nUserNS=nomap",
			"test5.container",
		),
	}

	for _, s := range cases {
		diags := qsr015(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR016_NoParameters(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nUserNS=auto:gid=101",
		"test.container",
	)

	diags := qsr016(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr016", *diags[0].Source)
	assert.Equal(t, "Invalid value of UserNS: 'auto' has no parameters", diags[0].Message)
}

func TestQSR016_KeepIdWrongParameter(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nUserNS=keep-id:gid=101,foo=101",
		"test.container",
	)

	diags := qsr016(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr016", *diags[0].Source)
	assert.Equal(t, "Invalid value of UserNS: [uid gid] allowed but found foo=101", diags[0].Message)
}

func TestQSR016_InvalidValue(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nUserNS=foo",
		"test.container",
	)

	diags := qsr016(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr016", *diags[0].Source)
	assert.Equal(t, "Invalid value of UserNS: allowed values: '[auto host keep-id nomap]' and found foo", diags[0].Message)
}
