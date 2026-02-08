package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR012_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]Secret=my-secret",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]Secret=my-secret,target=/app/pw",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]Secret=my-secret,type=mount,target=/app/pw",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]Secret=my-secret,uid=69,gid=420,mode=0400",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]Secret=my-secret,type=mount,uid=69,gid=420,mode=0400",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr012(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR012_UnfinishedOption(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nSecret=my-secret,type\n",
		"foo.container",
	)

	diags := qsr012(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr012", *diags[0].Source)
	assert.Equal(t, "Invalid format of secret specification: 'type' has no value", diags[0].Message)
}

func TestQSR012_InvalidType(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nSecret=my-secret,type=foo\n",
		"foo.container",
	)

	diags := qsr012(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr012", *diags[0].Source)
	assert.Equal(t, "Invalid format of secret specification: 'type' can be either 'mount' or 'env'", diags[0].Message)
}

func TestQSR012_InvalidOption(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nSecret=my-secret,foo=bar\n",
		"foo.container",
	)

	diags := qsr012(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr012", *diags[0].Source)
	assert.Equal(t, "Invalid format of secret specification: 'foo' is invalid option", diags[0].Message)
}

func TestQSR012_InvalidWithEnv(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nSecret=my-secret,type=env,uid=69,gid=420,mode=0400\n",
			"test1.container",
		),
	}

	for _, s := range cases {
		diags := qsr012(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr012", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Invalid format of secret specification: ", diags[0].Message)
		assert.Contains(t, diags[0].Message, "' only allowed if type=mount", diags[0].Message)
	}
}
