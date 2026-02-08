package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR005_Valid(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nImage=cr.io/org/cont\nAutoUpdate=registry",
		"test.container",
	)

	diags := qsr005(s)
	require.Len(t, diags, 0)
}

func TestQSR005_ValidKube(t *testing.T) {
	s := NewSyntaxChecker(
		"[Kube]\nYAML=test.yaml\nAutoUpdate=local",
		"test.kube",
	)

	diags := qsr005(s)
	require.Len(t, diags, 0)
}

func TestQSR005_Invalid(t *testing.T) {
	s := NewSyntaxChecker(
		"[Container]\nImage=cr.io/org/cont\nAutoUpdate=foo",
		"test.container",
	)

	diags := qsr005(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr005", *diags[0].Source)
	assert.Equal(t, "Invalid value of AutoUpdate: foo", diags[0].Message)
}
