package syntax

import (
	"strings"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func TestQSR021_InvalidOld(t *testing.T) {
	types := []string{
		"Wants",
		"Requires",
		"Requisite",
		"BindsTo",
		"PartOf",
		"Upholds",
		"Conflicts",
		"Before",
		"After",
	}

	for _, ty := range types {
		s := SyntaxChecker{
			documentText: "[Unit]\n" + ty + "=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		}
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
	types := []string{
		"Wants",
		"Requires",
		"Requisite",
		"BindsTo",
		"PartOf",
		"Upholds",
		"Conflicts",
		"Before",
		"After",
	}

	for _, ty := range types {
		s := SyntaxChecker{
			documentText: "[Unit]\n" + ty + "=test1.volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		}
		diags := qsr021(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, but got %d", len(diags))
		}
	}
}

func TestQSR021_ValidOld(t *testing.T) {
	types := []string{
		"Wants",
		"Requires",
		"Requisite",
		"BindsTo",
		"PartOf",
		"Upholds",
		"Conflicts",
		"Before",
		"After",
	}

	for _, ty := range types {
		s := SyntaxChecker{
			documentText: "[Unit]\n" + ty + "=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 4, 0),
			},
		}
		diags := qsr021(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR021_ValidOldWithNew(t *testing.T) {
	types := []string{
		"Wants",
		"Requires",
		"Requisite",
		"BindsTo",
		"PartOf",
		"Upholds",
		"Conflicts",
		"Before",
		"After",
	}

	for _, ty := range types {
		s := SyntaxChecker{
			documentText: "[Unit]\n" + ty + "=test1-volume.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		}
		diags := qsr021(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, but got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR021_ValidWantsTemplate(t *testing.T) {
	cases := []SyntaxChecker{
		{
			documentText: "[Unit]\nWants=webapp@8081.service\nWants=webapp@8082.service\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		},
		{
			documentText: "[Unit]\nWants=webapp@8081.container\nWants=webapp@8082.service\n[Container]\nImage=my-image.image",
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
	types := []string{
		"Wants",
		"Requires",
		"Requisite",
		"BindsTo",
		"PartOf",
		"Upholds",
		"Conflicts",
		"Before",
		"After",
	}
	for _, ty := range types {
		s := SyntaxChecker{
			documentText: "[Unit]\n" + ty + "=test1-volume\n[Container]\nImage=my-image.image",
			uri:          "test1.container",
			config: &utils.QuadletConfig{
				Podman: utils.BuildPodmanVersion(5, 5, 2),
			},
		}
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

func TestQSR021_TestServiceRegexp(t *testing.T) {
	inputs := []string{
		"foobar@:\\_.-.service",
		"foobar@:\\_.-.socket",
		"foobar@:\\_.-.device",
		"foobar@:\\_.-.mount",
		"foobar@:\\_.-.automount",
		"foobar@:\\_.-.swap",
		"foobar@:\\_.-.target",
		"foobar@:\\_.-.path",
		"foobar@:\\_.-.timer",
		"foobar@:\\_.-.slice",
		"foobar@:\\_.-.scope",
	}

	for _, s := range inputs {
		if !qsr021ServiceNamingConvention.MatchString(s) {
			t.Fatalf("expected to match but does not: %s", s)
		}
	}
}

func TestQSR021_TestQuadletRegexp(t *testing.T) {
	inputs := []string{
		"foobar@_.-.image",
		"foobar@_.-.container",
		"foobar@_.-.volume",
		"foobar@_.-.network",
		"foobar@_.-.kube",
		"foobar@_.-.pod",
		"foobar@_.-.build",
	}

	for _, s := range inputs {
		if !qsr021QuadletNamingConvention.MatchString(s) {
			t.Fatalf("expected to match but does not: %s", s)
		}
	}
}
