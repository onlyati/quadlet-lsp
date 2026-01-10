package syntax

import (
	"os"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestQSR014_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"net1.network",
		"[Network]",
	)
	createTempFile(
		t,
		tmpDir,
		"net2.network",
		"[Network]",
	)

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test1.container",
		),
		NewSyntaxChecker(
			"[Pod]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.pod",
		),
		NewSyntaxChecker(
			"[Build]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.build",
		),
		NewSyntaxChecker(
			"[Kube]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.kube",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr014(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR014_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test1.container",
		),
		NewSyntaxChecker(
			"[Pod]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.pod",
		),
		NewSyntaxChecker(
			"[Build]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.build",
		),
		NewSyntaxChecker(
			"[Kube]\nNetwork=net1.network\nNetwork=net2.network",
			"file://"+tmpDir+"/test2.kube",
		),
	}

	for _, s := range cases {
		s.config = &utils.QuadletConfig{
			WorkspaceRoot: tmpDir,
			Project: utils.ProjectProperty{
				DirLevel: utils.ReturnAsPtr(2),
			},
		}
		diags := qsr014(s)

		if len(diags) != 2 {
			t.Fatalf("Expected 2 diagnostics, got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr014" {
			t.Fatalf("Wrong source found: %s at %s", *diags[0].Source, s.uri)
		}

		if diags[0].Message != "Network file does not exists: net1.network" {
			t.Fatalf("Unexpected message: '%s' at %s", diags[0].Message, s.uri)
		}
	}
}
