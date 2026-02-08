package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR004_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nImage=docker.io/library/debian:bookworm-slim",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=ghcr.io/henrygd/beszel/beszel-agent:latest",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=localhost/test",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=localhost:5000/test",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=example.com:5000/test",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Artifact]\nArtifact=example.com:5000/test",
			"test3.artifact",
		),
	}

	for _, s := range cases {
		diags := qsr004(s)
		require.Len(t, diags, 0)
	}
}

func TestQSR004_ValidImageFile(t *testing.T) {
	s := NewSyntaxChecker("[Image]\nImage=docker.io/library/debian:bookworm-slim", "test.image")
	diags := qsr004(s)
	require.Len(t, diags, 0)
}

func TestQSR004_ValidWithImage(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nImage=db.image", "test.container")
	diags := qsr004(s)
	require.Len(t, diags, 0)
}

func TestQSR004_ValidWithBuild(t *testing.T) {
	s := NewSyntaxChecker("[Container]\nImage=db.build", "test.container")
	diags := qsr004(s)
	require.Len(t, diags, 0)
}

func TestQSR004_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nImage=library/debian:bookworm-slim",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=localhost",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nImage=localhost/",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Artifact]\nArtifact=localhost/",
			"test1.artifact",
		),
	}

	for _, s := range cases {
		diags := qsr004(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr004", *diags[0].Source)
		assert.Equal(t, "Image name is not fully qualified", diags[0].Message)
	}
}

func TestQSR004_InvalidImageFile(t *testing.T) {
	s := NewSyntaxChecker("[Image]\nImage=library/debian:bookworm-slim", "test.image")
	diags := qsr004(s)
	require.Len(t, diags, 1)
	require.NotNil(t, diags[0].Source)
	assert.Equal(t, "quadlet-lsp.qsr004", *diags[0].Source)
	assert.Equal(t, "Image name is not fully qualified", diags[0].Message)
}

func TestQSR004_NonContainer(t *testing.T) {
	s := NewSyntaxChecker("[Pod]\nImage=library/debian:bookworm-slim", "test.pod")
	diags := qsr004(s)
	require.Len(t, diags, 0)
}
