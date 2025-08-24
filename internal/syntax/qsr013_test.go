package syntax

import (
	"os"
	"testing"
)

func TestQSR013_Valid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	createTempFile(
		t,
		tmpDir,
		"data1.volume",
		"[Volume]",
	)
	createTempFile(
		t,
		tmpDir,
		"data2.volume",
		"[Volume]",
	)
	createTempFile(
		t,
		tmpDir,
		"data@.volume",
		"[Volume]",
	)

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test1.container",
		),
		NewSyntaxChecker(
			"[Pod]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test2.pod",
		),
		NewSyntaxChecker(
			"[Build]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test2.build",
		),
		NewSyntaxChecker(
			"[Build]\nVolume=data@%i.volume:/app:r\nVolume=data@test.volume:/data/:rw",
			"file://"+tmpDir+"/test2.build",
		),
	}

	for _, s := range cases {
		diags := qsr013(s)

		if len(diags) != 0 {
			t.Fatalf("Expected 0 diagnostics, got %d at %s", len(diags), s.uri)
		}
	}
}

func TestQSR013_Invalid(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)

	cases := []SyntaxChecker{
		NewSyntaxChecker(
			"[Container]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test1.container",
		),
		NewSyntaxChecker(
			"[Pod]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test2.pod",
		),
		NewSyntaxChecker(
			"[Build]\nVolume=data1.volume:/app:r\nVolume=data2.volume:/data/:rw",
			"file://"+tmpDir+"/test2.build",
		),
	}

	for _, s := range cases {
		diags := qsr013(s)

		if len(diags) != 2 {
			t.Fatalf("Expected 2 diagnostics, got %d at %s", len(diags), s.uri)
		}

		if *diags[0].Source != "quadlet-lsp.qsr013" {
			t.Fatalf("Wrong source found: %s at %s", *diags[0].Source, s.uri)
		}

		if diags[0].Message != "Volume file does not exists: data1.volume" {
			t.Fatalf("Unexpected message: '%s' at %s", diags[0].Message, s.uri)
		}
	}
}
