package parser

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	inputContent := `# Test container

[Container]
Image=foo.image
Label= \
  env=test

[Unit]
Description=Foo \
  container
`
	tmpDir := t.TempDir()
	testutils.CreateTempFile(t, tmpDir, "foo.container", inputContent)

	parser := NewParser(path.Join(tmpDir, "foo.container"))
	parser.Run()

	require.Len(t, parser.Errors, 0)

	// # Test container
	require.Len(t, parser.Quadlet.Documents, 1, "missing parsed document")
	assert.Equal(t, "# Test container", *parser.Quadlet.Documents[0].Text)
	assert.Equal(t, NodePosition{0, 0}, parser.Quadlet.Documents[0].StartPos)
	assert.Equal(t, NodePosition{0, 16}, parser.Quadlet.Documents[0].EndPos)

	require.Len(t, parser.Quadlet.Sections, 2, "missing secion")

	// [Container]
	containerSection := parser.Quadlet.Sections[0]
	assert.Equal(t, "[Container]", *containerSection.Text)
	assert.Equal(t, NodePosition{2, 0}, containerSection.StartPos)
	assert.Equal(t, NodePosition{2, 11}, containerSection.EndPos)
	require.Len(t, parser.Quadlet.Sections[0].Assignments, 2)

	// [Container] => Image=foo.image
	imageKeyword := containerSection.Assignments[0]
	assert.Equal(t, "Image", *imageKeyword.Name)
	assert.Equal(t, NodePosition{3, 0}, imageKeyword.StartPos)
	assert.Equal(t, NodePosition{3, 15}, imageKeyword.EndPos)

	imageKeywordValue := imageKeyword.Value
	assert.Equal(t, "foo.image", *imageKeywordValue.Value)
	assert.Equal(t, NodePosition{3, 6}, imageKeywordValue.StartPos)
	assert.Equal(t, NodePosition{3, 15}, imageKeywordValue.EndPos)

	// [Container] => Label=env=test
	labelKeyword := containerSection.Assignments[1]
	assert.Equal(t, "Label", *labelKeyword.Name)
	assert.Equal(t, NodePosition{4, 0}, labelKeyword.StartPos)
	assert.Equal(t, NodePosition{5, 10}, labelKeyword.EndPos)

	labelKeywordValue := containerSection.Assignments[1].Value
	assert.Equal(t, "env=test", *labelKeywordValue.Value)
	assert.Equal(t, NodePosition{5, 2}, labelKeywordValue.StartPos)
	assert.Equal(t, NodePosition{5, 10}, labelKeywordValue.EndPos)

	// [Unit]
	unitSection := parser.Quadlet.Sections[1]
	assert.Equal(t, "[Unit]", *unitSection.Text)
	assert.Equal(t, NodePosition{7, 0}, unitSection.StartPos)
	assert.Equal(t, NodePosition{7, 6}, unitSection.EndPos)
	require.Len(t, unitSection.Assignments, 1)

	// [Unit] => Description=Foo container
	descKeyword := unitSection.Assignments[0]
	assert.Equal(t, "Description", *descKeyword.Name)
	assert.Equal(t, NodePosition{8, 0}, descKeyword.StartPos)
	assert.Equal(t, NodePosition{9, 11}, descKeyword.EndPos)

	descKeywordValue := unitSection.Assignments[0].Value
	assert.Equal(t, "Foo container", *descKeywordValue.Value)
	assert.Equal(t, NodePosition{8, 12}, descKeywordValue.StartPos)
	assert.Equal(t, NodePosition{9, 11}, descKeywordValue.EndPos)
}
