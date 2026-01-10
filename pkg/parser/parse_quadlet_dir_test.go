package parser_test

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseQuadletDir(t *testing.T) {
	tmpDir := t.TempDir()

	createTempFile(
		t,
		tmpDir,
		"foo.container",
		`[Container]
Image=foo.image
Exec=tail \
  -f /dev/null
`)
	createTempFile(
		t,
		tmpDir,
		"foo.image",
		`[Image]
Image=docker.io/library/debian
`)

	// Dropins for foo.container
	createTempDir(t, tmpDir, "foo.container.d")
	createTempFile(t, path.Join(tmpDir, "foo.container.d"), "labels.conf", `
[Container]
Label="app=foo"
`)

	createTempFile(
		t,
		tmpDir,
		".quadletrc.json",
		`{ "disable": ["qsr001", "qsr002"] }`,
	)

	// Test application in nested directory
	createTempDir(t, tmpDir, "app")
	createTempFile(t, path.Join(tmpDir, "app"), "bar.container", `
[Container]
Image=bar.image
`)
	createTempDir(t, path.Join(tmpDir, "app"), "bar.container.d")
	createTempFile(t, path.Join(tmpDir, "app", "bar.container.d"), "labels.conf", `
[Container]
Label="app=bar"
`)

	// Still should work even outside of the nested direcrory
	createTempDir(t, tmpDir, "bar.container.d")
	createTempFile(t, path.Join(tmpDir, "bar.container.d"), "labels.conf", `
[Container]
Label="app=bar"
`)

	// An oprhan dropins, no 'baz.container' exists
	createTempDir(t, tmpDir, "baz.container.d")
	createTempFile(t, path.Join(tmpDir, "baz.container.d"), "labels.conf", `
[Container]
Label="app=baz"
`)

	expected := parser.QuadletDirectory{
		DisabledQSR: []string{"qsr001", "qsr002"},
		OprhanDropins: []parser.Dropin{
			{
				Directory: "baz.container.d",
				FileName:  "baz.container.d/labels.conf",
				Properties: map[string][]parser.QuadletProperty{
					"Container": {
						{"Label", "\"app=baz\""},
					},
				},
				SourceFile: `
[Container]
Label="app=baz"
`,
			},
		},
		Quadlets: map[string]parser.Quadlet{
			"app/bar.container": {
				DisabledQSR: nil,
				Name:        "app/bar.container",
				References:  []string{"bar.image"},
				PartOf:      nil,
				Dropins: []parser.Dropin{
					{
						Directory: "bar.container.d",
						FileName:  "app/bar.container.d/labels.conf",
						Properties: map[string][]parser.QuadletProperty{
							"Container": {
								{"Label", "\"app=bar\""},
							},
						},
						SourceFile: `
[Container]
Label="app=bar"
`,
					},
					{
						Directory: "bar.container.d",
						FileName:  "bar.container.d/labels.conf",
						Properties: map[string][]parser.QuadletProperty{
							"Container": {
								{"Label", "\"app=bar\""},
							},
						},
						SourceFile: `
[Container]
Label="app=bar"
`,
					},
				},
				Header: nil,
				Properties: map[string][]parser.QuadletProperty{
					"Container": {
						{"Image", "bar.image"},
					},
				},
				SourceFile: `
[Container]
Image=bar.image
`,
			},
			"foo.container": {
				DisabledQSR: nil,
				Name:        "foo.container",
				References:  []string{"foo.image"},
				PartOf:      nil,
				Dropins: []parser.Dropin{
					{
						Directory: "foo.container.d",
						FileName:  "foo.container.d/labels.conf",
						Properties: map[string][]parser.QuadletProperty{
							"Container": {
								{"Label", "\"app=foo\""},
							},
						},
						SourceFile: `
[Container]
Label="app=foo"
`,
					},
				},
				Header: nil,
				Properties: map[string][]parser.QuadletProperty{
					"Container": {
						{"Image", "foo.image"},
						{"Exec", "tail -f /dev/null"},
					},
				},
				SourceFile: `[Container]
Image=foo.image
Exec=tail \
  -f /dev/null
`,
			},
			"foo.image": {
				DisabledQSR: nil,
				Name:        "foo.image",
				References:  nil,
				PartOf:      []string{"foo.container"},
				Dropins:     nil,
				Header:      nil,
				Properties: map[string][]parser.QuadletProperty{
					"Image": {
						{"Image", "docker.io/library/debian"},
					},
				},
				SourceFile: `[Image]
Image=docker.io/library/debian
`,
			},
		},
	}

	result, err := parser.ParseQuadletDir(tmpDir, 2)

	require.NoError(
		t,
		err,
		"failed to parse directory",
	)
	assert.Equal(
		t,
		expected.DisabledQSR,
		result.DisabledQSR,
		"wrongly parsed disable rules",
	)
	assert.Equal(
		t,
		expected.Quadlets,
		result.Quadlets,
		"wrongly parsed quadlets",
	)
	assert.Equal(
		t,
		expected.OprhanDropins,
		result.OprhanDropins,
	)
}
