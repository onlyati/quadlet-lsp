package syntax

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQSR010_Valid(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:420:69",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Pod]\nPublishPort=10.0.0.1:420:69",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Kube]\nPublishPort=10.0.0.1:420:69",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=420:69",
			"test4.container",
		),
		NewSyntaxChecker(
			"[Pod]\nPublishPort=420:69",
			"test5.container",
		),
		NewSyntaxChecker(
			"[Kube]\nPublishPort=420:69",
			"test6.container",
		),
		NewSyntaxChecker(
			"[Kube]\nPublishPort=:69",
			"test7.container",
		),
		NewSyntaxChecker(
			"[Kube]\nPublishPort=10.0.0.1::69",
			"test8.container",
		),
		// Test cases for protocol suffix support
		// See: https://docs.podman.io/en/v5.0.1/markdown/podman-run.1.html
		// Format: [[ip:]hostPort:]containerPort[/protocol]
		NewSyntaxChecker(
			"[Container]\nPublishPort=22000:22000/tcp",
			"test9.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=22000:22000/udp",
			"test10.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=22000:22000/sctp",
			"test11.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:22000:22000/tcp",
			"test12.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:22000:22000/udp",
			"test13.container",
		),
	}

	for _, s := range cases {
		d := qsr010(s)
		require.Len(t, d, 0)
	}
}

func TestQSR010_InvalidFormat(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=420",
			"test2.container",
		),
	}

	for _, s := range cases {
		diags := qsr010(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr010", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Incorrect format of PublishPort")
	}
}

func TestQSR010_InvalidPortIsText(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:nice:420",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:69:ez",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=nice:420",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=69:ez",
			"test4.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=69:",
			"test4.container",
		),
	}

	for _, s := range cases {
		diags := qsr010(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr010", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Incorrect format of PublishPort")
	}
}

func TestQSR010_InvalidInvalidPortNumber(t *testing.T) {
	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:-69:420",
			"test1.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=10.0.0.1:69:80000",
			"test2.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=-69:420",
			"test3.container",
		),
		NewSyntaxChecker(
			"[Container]\nPublishPort=69:80000",
			"test4.container",
		),
	}

	for _, s := range cases {
		diags := qsr010(s)
		require.Len(t, diags, 1)
		require.NotNil(t, diags[0].Source)
		assert.Equal(t, "quadlet-lsp.qsr010", *diags[0].Source)
		assert.Contains(t, diags[0].Message, "Incorrect format of PublishPort")
	}
}
