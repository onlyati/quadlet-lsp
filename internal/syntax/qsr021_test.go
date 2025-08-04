package syntax

import (
	"strings"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestQSR021_InvalidOld(t *testing.T) {
	cases := []SyntaxChecker{
		{
			documentText: "[Unit]\nWants=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nRequires=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nRequisite=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nBindsTo=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nPartOf=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nUpholds=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nConflicts=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nBefore=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nAfter=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
	}

	for _, s := range cases {
		diags := qsr021(s)

		if len(diags) != 1 {
			t.Fatalf("Expected 1 diagnostics, but got %d", len(diags))
		}

		if *diags[0].Source != "quadlet-lsp.qsr021" {
			t.Fatalf("Exptexted quadlet-lsp.qsr012 source, but got %s", *diags[0].Source)
		}

		checkMessage := strings.HasPrefix(diags[0].Message, "Invalid depdency is specified: ")
		if !checkMessage {
			t.Fatalf("Unexpected message returned: %s", diags[0].Message)
		}
	}
}

func TestQSR021_ValidNew(t *testing.T) {
	cases := []SyntaxChecker{
		{
			documentText: "[Unit]\nWants=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nRequires=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nRequisite=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nBindsTo=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nPartOf=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nUpholds=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nConflicts=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nBefore=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nAfter=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
	}

	for _, s := range cases {
		diags := qsr021(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, but got %d", len(diags))
		}
	}
}

func TestQSR021_ValidOld(t *testing.T) {
	cases := []SyntaxChecker{
		{
			documentText: "[Unit]\nWants=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nRequires=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nRequisite=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nBindsTo=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nPartOf=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nUpholds=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nConflicts=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nBefore=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
		{
			documentText: "[Unit]\nAfter=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		},
	}

	for _, s := range cases {
		diags := qsr021(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, but got %d", len(diags))
		}
	}
}

func TestQSR021_ValidOldWithNew(t *testing.T) {
	cases := []SyntaxChecker{
		{
			documentText: "[Unit]\nWants=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nRequires=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nRequisite=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nBindsTo=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nPartOf=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nUpholds=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nConflicts=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nBefore=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nAfter=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
	}

	for _, s := range cases {
		diags := qsr021(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, but got %d", len(diags))
		}
	}
}

func TestQSR021_Invalid(t *testing.T) {
	cases := []SyntaxChecker{
		{
			documentText: "[Unit]\nWants=test1-volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nRequires=test1-volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nRequisite=test1-volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nBindsTo=test1-volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nPartOf=test1-volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nUpholds=test1-volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nConflicts=test1-volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nBefore=test1-volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nAfter=test1-volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
	}

	for _, s := range cases {
		diags := qsr021(s)

		if len(diags) != 1 {
			t.Fatalf("Expected 1 diagnostics, but got %d", len(diags))
		}

		if *diags[0].Source != "quadlet-lsp.qsr021" {
			t.Fatalf("Exptexted quadlet-lsp.qsr012 source, but got %s", *diags[0].Source)
		}

		checkMessage := strings.HasPrefix(diags[0].Message, "Invalid depdency is specified: ")
		if !checkMessage {
			t.Fatalf("Unexpected message returned: %s", diags[0].Message)
		}
	}
}
