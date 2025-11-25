package parser_test

import (
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
		` [Container]
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

	createTempFile(
		t,
		tmpDir,
		".quadletrc.json",
		`{ "disable": ["qsr001", "qsr002"] }`,
	)

	expected := parser.QuadletDirectory{
		DisabledQSR: []string{"qsr001", "qsr002"},
		Quadlets: map[string]parser.Quadlet{
			"foo.container": {
				DisabledQSR: nil,
				Name:        "foo.container",
				References:  []string{"foo.image"},
				PartOf:      nil,
				Dropins:     nil,
				Header:      nil,
				Properties: map[string][]parser.QuadletProperty{
					"Container": {
						{"Image", "foo.image"},
						{"Exec", "tail -f /dev/null"},
					},
				},
				SourceFile: ` [Container]
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

	result, err := parser.ParseQuadletDir(tmpDir)

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
}
