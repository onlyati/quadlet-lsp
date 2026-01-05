package utils_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestFindItems(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	text := `[Unit]
Description=description

[Container]
Image=docker.io/library/debian:bookworm-slim
AutoUpdate=registry
Environment=ENV1=VALUE1
Volume=dev-db-volume:/app:rw
Environment=ENV2=VALUE2
# Environment=ENV3=VALUE3

[Service]
Restart=on-failure
RestartSec=5
StartLimitBurst=5
`
	findings := utils.FindItems(
		utils.FindItemProperty{
			RootDirectory: tmpDir,
			Text:          text,
			Section:       "[Container]",
			Property:      "Environment",
			URI:           "file://" + tmpDir + "foo.container",
		},
	)

	if len(findings) != 2 {
		t.Fatalf("Expected 2 founds, got %d", len(findings))
	}

	if findings[0].Property != "Environment" {
		t.Fatalf("Expected to get 'Environment' property, but got %s", findings[0].Property)
	}

	if findings[0].Value != "ENV1=VALUE1" {
		t.Fatalf("Expected to get 'ENV1=VALUE1' value, but got %s", findings[0].Value)
	}

	if findings[0].LineNumber != 6 {
		t.Fatalf("Expected to find in 6th line, but got %d", findings[0].LineNumber)
	}
}

func TestFindItemsWithExec(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	text := `[Unit]
Description=description

[Container]
Image=docker.io/library/debian:bookworm-slim
AutoUpdate=registry
Environment=ENV1=VALUE1
Volume=dev-db-volume:/app:rw
Exec=tail \
    -f \
    /dev/null
Environment=ENV2=VALUE2
# Environment=ENV3=VALUE3

[Service]
Restart=on-failure
RestartSec=5
StartLimitBurst=5
`
	findings := utils.FindItems(
		utils.FindItemProperty{
			RootDirectory: tmpDir,
			Text:          text,
			Section:       "[Container]",
			Property:      "Exec",
			URI:           "file://" + tmpDir + "foo.container",
		},
	)

	if len(findings) != 1 {
		t.Fatalf("Expected 1 founds, got %d", len(findings))
	}

	if findings[0].Property != "Exec" {
		t.Fatalf("Expected to get 'Environment' property, but got %s", findings[0].Property)
	}

	if findings[0].Value != "tail -f /dev/null" {
		t.Fatalf("Expected to get 'tail -f /dev/null' value, but got '%s'", findings[0].Value)
	}

	if findings[0].LineNumber != 8 {
		t.Fatalf("Expected to find in 8th line, but got %d", findings[0].LineNumber)
	}
}

func TestScanQuadlet(t *testing.T) {
	text := `[Unit]
Description=description

[Container]
Image=docker.io/library/debian:bookworm-slim
AutoUpdate=registry
Environment=ENV1=VALUE1
Volume=dev-db-volume:/app:rw
Exec=tail \
    -f \
    /dev/null
Environment=ENV2=VALUE2
# Environment=ENV3=VALUE3

[Service]
Restart=on-failure
RestartSec=5
StartLimitBurst=5
`

	findings := []struct {
		p string
		v string
		c string
	}{}
	mockFn := func(q utils.QuadletLine, _ utils.PodmanVersion, _ any) []protocol.Diagnostic {
		findings = append(findings, struct {
			p string
			v string
			c string
		}{p: q.Property, v: q.Value, c: q.Section})
		return nil
	}

	_ = utils.ScanQadlet(
		text,
		utils.PodmanVersion{},
		map[utils.ScanProperty]struct{}{
			{Section: "[Container]", Property: "Environment"}: {},
			{Section: "[Container]", Property: "Exec"}:        {},
		},
		mockFn,
		nil,
	)

	if len(findings) != 3 {
		t.Fatalf("execpted 3 finding but got %d", len(findings))
	}

	if findings[0].c != "[Container]" {
		t.Fatalf("expected '[Container]' section but got '%s'", findings[0].c)
	}

	if findings[0].p != "Environment" && findings[0].v != "ENV1=VALUE1" {
		t.Fatalf("unexpected finding for 0: '%s=%s'", findings[0].p, findings[0].v)
	}

	if findings[1].p != "Exec" && findings[1].v != "tail -f /dev/null" {
		t.Fatalf("unexpected finding for 1: '%s=%s'", findings[1].p, findings[1].v)
	}

	if findings[2].p != "Environment" && findings[2].v != "ENV2=VALUE2" {
		t.Fatalf("unexpected finding for 2: '%s=%s'", findings[2].p, findings[2].v)
	}
}

func TestScanQuadletAll(t *testing.T) {
	text := `[Unit]
Description=description

[Container]
Image=docker.io/library/debian:bookworm-slim
Exec=tail \
    -f \
    /dev/null
AutoUpdate=registry

[Service]
Restart=on-failure

[Test]
NoExist=true
`

	findings := []struct {
		p string
		v string
		c string
	}{}
	mockFn := func(q utils.QuadletLine, _ utils.PodmanVersion, _ any) []protocol.Diagnostic {
		findings = append(findings, struct {
			p string
			v string
			c string
		}{p: q.Property, v: q.Value, c: q.Section})
		return nil
	}

	_ = utils.ScanQadlet(
		text,
		utils.PodmanVersion{},
		map[utils.ScanProperty]struct{}{
			{Section: "*", Property: "*"}: {},
		},
		mockFn,
		nil,
	)

	if len(findings) != 10 {
		t.Fatalf("execpted 6 finding but got %d", len(findings))
	}

	expectedFindings := []struct {
		p string
		v string
		c string
	}{
		{c: "[Unit]", p: "", v: ""},
		{c: "[Unit]", p: "Description", v: "description"},
		{c: "[Container]", p: "", v: ""},
		{c: "[Container]", p: "Image", v: "docker.io/library/debian:bookworm-slim"},
		{c: "[Container]", p: "Exec", v: "tail -f /dev/null"},
		{c: "[Container]", p: "AutoUpdate", v: "registry"},
		{c: "[Service]", p: "", v: ""},
		{c: "[Service]", p: "Restart", v: "on-failure"},
		{c: "[Test]", p: "", v: ""},
		{c: "[Test]", p: "NoExist", v: "true"},
	}

	for i, e := range expectedFindings {
		if findings[i] != e {
			t.Fatalf("unexpected finding: %+v", findings[0])
		}
	}
}

func TestFindReferences(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "example.container", "[Container]\nNetwork=example.network\nAnotherLine")

	locations, err := findReferences(
		utils.GoReferenceProperty{
			Property: "Network",
			SearchIn: []string{"container", "pod", "kube"},
			DirLevel: 2,
		}, "example.network", tmpDir)
	assert.NoError(t, err)
	assert.Len(t, locations, 1)
	assert.Contains(t, string(locations[0].URI), "example.container")
	assert.Equal(t, uint32(1), locations[0].Range.Start.Line)
}

func TestFindReferencesTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "web@.container", "[Container]\nVolume=web@%i.volume:/app")
	createTempFile(t, tmpDir, "builder@.container", "[Container]\nVolume=web@%i.volume:/app")

	locations, err := findReferences(
		utils.GoReferenceProperty{
			Property: "Volume",
			SearchIn: []string{"container", "pod"},
			DirLevel: 2,
		}, "web@.volume", tmpDir)
	assert.NoError(t, err)
	assert.Len(t, locations, 2)

	for _, loc := range locations {
		if !strings.Contains(loc.URI, "web@.container") && !strings.Contains(loc.URI, "builder@.container") {
			t.Fatalf("Unexpected finding: %+v", loc)
		}
	}
}

func TestFindReferencesDropIns(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempFile(t, tmpDir, "foo.container", "[Container]\nImage=foo.image\n")
	createTempFile(t, tmpDir, "foo.pod", "[Pod]\n")
	createTempDir(t, tmpDir, "foo.container.d")
	createTempFile(t, path.Join(tmpDir, "foo.container.d"), "pod.conf", "[Container]\nPod=foo.pod\n")

	locations, err := findReferences(
		utils.GoReferenceProperty{
			Property: "Pod",
			SearchIn: []string{"container", "kube", "volume", "network", "image", "build"},
			DirLevel: 2,
		}, "foo.pod", tmpDir)
	assert.NoError(t, err)
	assert.Len(t, locations, 1)
	assert.Contains(t, string(locations[0].URI), "foo.container.d/pod.conf")
	assert.Equal(t, uint32(1), locations[0].Range.Start.Line)
}

func TestFindReferencesNested(t *testing.T) {
	tmpDir := t.TempDir()
	_ = os.Chdir(tmpDir)

	createTempDir(t, tmpDir, "foo")
	createTempFile(t, path.Join(tmpDir, "foo"), "foo.container", "[Container]\nImage=foo.image\n")
	createTempFile(t, path.Join(tmpDir, "foo"), "foo.pod", "[Pod]\n")
	createTempDir(t, path.Join(tmpDir, "foo"), "foo.container.d")
	createTempFile(t, path.Join(tmpDir, "foo", "foo.container.d"), "pod.conf", "[Container]\nPod=foo.pod\n")

	locations, err := findReferences(
		utils.GoReferenceProperty{
			Property: "Pod",
			SearchIn: []string{"container", "kube", "volume", "network", "image", "build"},
			DirLevel: 2,
		}, "foo.pod", tmpDir)
	assert.NoError(t, err)
	assert.Len(t, locations, 1)
	assert.Contains(t, string(locations[0].URI), "foo/foo.container.d/pod.conf")
	assert.Equal(t, uint32(1), locations[0].Range.Start.Line)
}
