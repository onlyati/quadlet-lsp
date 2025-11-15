package parser_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/onlyati/quadlet-lsp/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0o644)
	assert.NoError(t, err)
	return path
}

func createTempDir(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	err := os.Mkdir(path, 0o755)
	assert.NoError(t, err)
	return path
}

func Test_ParseQuadlet(t *testing.T) {
	tmpDir := t.TempDir()

	// Create foo.container
	createTempFile(
		t,
		tmpDir,
		"foo-bar.container",
		`# disable-qsr: qsr014 qsr014
# disable-qsr: qsr007
;
; # Description
# Lorem ipsum dolor sit amet.

[Unit]
Description=Foo container

[Container]
Pod=foo.pod
Image=foo.image
Exec=tail \
  -f /dev/null
Volume=foo.volume:/app

# Network options
Network=foo.network

[Install]
WantedBy=default.target
`)

	// Create foo-.container.d/network.conf
	createTempDir(t, tmpDir, "foo-.container.d")
	createTempFile(
		t,
		path.Join(tmpDir, "foo-.container.d"),
		"network.conf",
		`
[Container]
PublishPort=8080:80
Network=bar.network
`)

	// Create foo-bar.container.d/network.conf
	createTempDir(t, tmpDir, "foo-bar.container.d")
	createTempFile(
		t,
		path.Join(tmpDir, "foo-bar.container.d"),
		"network.conf",
		`
[Container]
Network=foo-bar.network
`)

	// Create container.d/labels.conf
	createTempDir(t, tmpDir, "container.d")
	createTempFile(
		t,
		path.Join(tmpDir, "container.d"),
		"labels.conf",
		`
[Container]
Label="env.type=prod"
Label="env.server=app01"
`)

	expected := parser.Quadlet{
		DisabledQSR: []string{"qsr014", "qsr007"},
		Header: []string{
			"",
			"# Description",
			"Lorem ipsum dolor sit amet.",
		},
		Name: "foo-bar.container",
		Properties: map[string][]parser.QuadletProperty{
			"Unit": {
				{
					"Description", "Foo container",
				},
			},
			"Container": {
				{
					"Pod", "foo.pod",
				},
				{
					"Image", "foo.image",
				},
				{
					"Exec", "tail -f /dev/null",
				},
				{
					"Volume", "foo.volume:/app",
				},
				{
					"Network", "foo.network",
				},
			},
			"Install": {
				{
					"WantedBy", "default.target",
				},
			},
		},
		References: []string{
			"foo.pod",
			"foo.image",
			"foo.volume",
			"foo.network",
			"foo-bar.network",
			"bar.network",
		},
		SourceFile: `# disable-qsr: qsr014 qsr014
# disable-qsr: qsr007
;
; # Description
# Lorem ipsum dolor sit amet.

[Unit]
Description=Foo container

[Container]
Pod=foo.pod
Image=foo.image
Exec=tail \
  -f /dev/null
Volume=foo.volume:/app

# Network options
Network=foo.network

[Install]
WantedBy=default.target
`,
		Dropins: []parser.Dropin{
			{
				Directory: "foo-bar.container.d",
				FileName:  "network.conf",
				Properties: map[string][]parser.QuadletProperty{
					"Container": {
						{
							"Network", "foo-bar.network",
						},
					},
				},
				SourceFile: `
[Container]
Network=foo-bar.network
`,
			},
			{
				Directory: "foo-.container.d",
				FileName:  "network.conf",
				Properties: map[string][]parser.QuadletProperty{
					"Container": {
						{
							"PublishPort", "8080:80",
						},
						{
							"Network", "bar.network",
						},
					},
				},
				SourceFile: `
[Container]
PublishPort=8080:80
Network=bar.network
`,
			},
			{
				Directory: "container.d",
				FileName:  "labels.conf",
				Properties: map[string][]parser.QuadletProperty{
					"Container": {
						{
							"Label", "\"env.type=prod\"",
						},
						{
							"Label", "\"env.server=app01\"",
						},
					},
				},
				SourceFile: `
[Container]
Label="env.type=prod"
Label="env.server=app01"
`,
			},
		},
	}

	result, err := parser.ParseQuadlet(parser.ParseQuadletConfig{
		RootDirectory: tmpDir,
		FileName:      "foo-bar.container",
	})

	require.NoError(
		t,
		err,
		"error happened during quadlet scan",
	)
	assert.Equal(
		t,
		expected.Name,
		result.Name,
		"quadlet name does not match",
	)
	assert.Equal(
		t,
		expected.DisabledQSR,
		result.DisabledQSR,
		"wrongly parsed disabled QSRs",
	)
	assert.Equal(
		t,
		expected.Header,
		result.Header,
		"wrongly parsed header",
	)
	assert.Equal(
		t,
		expected.Properties,
		result.Properties,
		"wrongly parsed properties",
	)
	require.Equal(
		t,
		expected.Dropins,
		result.Dropins,
		"wrongly parsed dropins",
	)
	assert.ElementsMatch(
		t,
		expected.References,
		result.References,
		"wrongly calculated references",
	)
}

func Test_ParseQuadletImageOverride(t *testing.T) {
	tmpDir := t.TempDir()

	// Create foo.container
	createTempFile(
		t,
		tmpDir,
		"foo.container",
		` [Container]
Image=foo.image
`)

	// Create foo-.container.d/network.conf
	createTempDir(t, tmpDir, "foo.container.d")
	createTempFile(
		t,
		path.Join(tmpDir, "foo.container.d"),
		"image.conf",
		`
[Container]
Image=docker.io/library/debian
`)

	expected := parser.Quadlet{
		Name:        "foo.container",
		References:  nil,
		DisabledQSR: nil,
		Properties: map[string][]parser.QuadletProperty{
			"Container": {
				{
					"Image", "foo.image",
				},
			},
		},
		SourceFile: ` [Container]
Image=foo.image
`,
		Dropins: []parser.Dropin{
			{
				FileName:  "image.conf",
				Directory: "foo.container.d",
				Properties: map[string][]parser.QuadletProperty{
					"Container": {
						{
							"Image", "docker.io/library/debian",
						},
					},
				},
				SourceFile: `
[Container]
Image=docker.io/library/debian
`,
			},
		},
	}

	result, err := parser.ParseQuadlet(parser.ParseQuadletConfig{
		RootDirectory: tmpDir,
		FileName:      "foo.container",
	})

	require.NoError(
		t,
		err,
		"error happened during quadlet scan",
	)
	assert.Equal(
		t,
		expected.Name,
		result.Name,
		"quadlet name does not match",
	)
	assert.Equal(
		t,
		expected.DisabledQSR,
		result.DisabledQSR,
		"wrongly parsed disabled QSRs",
	)
	assert.Equal(
		t,
		expected.Header,
		result.Header,
		"wrongly parsed header",
	)
	assert.Equal(
		t,
		expected.Properties,
		result.Properties,
		"wrongly parsed properties",
	)
	require.Equal(
		t,
		expected.Dropins,
		result.Dropins,
		"wrongly parsed dropins",
	)
	assert.ElementsMatch(
		t,
		expected.References,
		result.References,
		"wrongly calculated references",
	)
}
