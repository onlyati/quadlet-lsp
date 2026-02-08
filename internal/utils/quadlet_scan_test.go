package utils_test

import (
	"path"
	"strings"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestFindItems(t *testing.T) {
	tmpDir := t.TempDir()

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
	require.Len(t, findings, 2)
	assert.Equal(t, "Environment", findings[0].Property)
	assert.Equal(t, "ENV1=VALUE1", findings[0].Value)
	assert.Equal(t, uint32(6), findings[0].LineNumber)

	assert.Equal(t, "Environment", findings[1].Property)
	assert.Equal(t, "ENV2=VALUE2", findings[1].Value)
	assert.Equal(t, uint32(8), findings[1].LineNumber)
}

func TestFindItemsWithExec(t *testing.T) {
	tmpDir := t.TempDir()

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
	require.Len(t, findings, 1)
	assert.Equal(t, "Exec", findings[0].Property)
	assert.Equal(t, "tail -f /dev/null", findings[0].Value)
	assert.Equal(t, uint32(8), findings[0].LineNumber)
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

	require.Len(t, findings, 3)
	assert.Equal(t, "[Container]", findings[0].c)
	assert.Equal(t, "Environment", findings[0].p)
	assert.Equal(t, "ENV1=VALUE1", findings[0].v)

	assert.Equal(t, "[Container]", findings[1].c)
	assert.Equal(t, "Exec", findings[1].p)
	assert.Equal(t, "tail -f /dev/null", findings[1].v)

	assert.Equal(t, "[Container]", findings[2].c)
	assert.Equal(t, "Environment", findings[2].p)
	assert.Equal(t, "ENV2=VALUE2", findings[2].v)
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
		assert.Equal(t, e, findings[i])
	}
}

func TestFindReferences(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "example.container", "[Container]\nNetwork=example.network\nAnotherLine")

	locations, err := utils.FindReferences(
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

	testutils.CreateTempFile(t, tmpDir, "web@.container", "[Container]\nVolume=web@%i.volume:/app")
	testutils.CreateTempFile(t, tmpDir, "builder@.container", "[Container]\nVolume=web@%i.volume:/app")

	locations, err := utils.FindReferences(
		utils.GoReferenceProperty{
			Property: "Volume",
			SearchIn: []string{"container", "pod"},
			DirLevel: 2,
		}, "web@.volume", tmpDir)
	assert.NoError(t, err)
	assert.Len(t, locations, 2)

	for _, loc := range locations {
		if !strings.Contains(loc.URI, "web@.container") && !strings.Contains(loc.URI, "builder@.container") {
			assert.Fail(t, "Unexpected finding: %+v", loc)
		}
	}
}

func TestFindReferencesDropIns(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempDir(t, tmpDir, "foo.container.d")
	testutils.CreateTempFile(t, tmpDir, "foo.container", "[Container]\nImage=foo.image\n")
	testutils.CreateTempFile(t, tmpDir, "foo.pod", "[Pod]\n")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo.container.d"), "pod.conf", "[Container]\nPod=foo.pod\n")

	locations, err := utils.FindReferences(
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

	testutils.CreateTempDir(t, tmpDir, "foo")
	testutils.CreateTempDir(t, path.Join(tmpDir, "foo"), "foo.container.d")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo"), "foo.container", "[Container]\nImage=foo.image\n")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo"), "foo.pod", "[Pod]\n")
	testutils.CreateTempFile(t, path.Join(tmpDir, "foo", "foo.container.d"), "pod.conf", "[Container]\nPod=foo.pod\n")

	locations, err := utils.FindReferences(
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
