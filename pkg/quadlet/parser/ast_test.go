package parser

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_QuadletFind(t *testing.T) {
	content := `  # Test container
# Second line

# Section doc
[Container]

# Image doc
Image=foo.image

# Label doc
Label= \
	env=test

# Unit doc
[Unit]
# Description doc
Description= \
	Foo \
	container
`

	tmpDir := t.TempDir()
	testutils.CreateTempFile(t, tmpDir, "foo.container", content)
	parser := NewParser(path.Join(tmpDir, "foo.container"))

	require.Len(t, parser.Errors, 0)
	require.Greater(t, len(parser.Quadlet.Sections), 0)

	quadlet := parser.Quadlet

	// Test outside of the line
	nodeResult := quadlet.FindToken(NodePosition{0, 35})
	require.Nil(t, nodeResult.CurrentNode)
	require.Len(t, nodeResult.ParentNodes, 0)

	nodeResult = quadlet.FindToken(NodePosition{0, 1})
	require.Nil(t, nodeResult.CurrentNode)
	require.Len(t, nodeResult.ParentNodes, 0)

	// First comment line
	nodeResult = quadlet.FindToken(NodePosition{0, 5})
	assert.Len(t, nodeResult.ParentNodes, 0)
	switch v := nodeResult.CurrentNode.(type) {
	case *CommentNode:
		assert.Equal(t, "# Test container", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// Second comment line
	nodeResult = quadlet.FindToken(NodePosition{1, 5})
	assert.Len(t, nodeResult.ParentNodes, 0)
	switch v := nodeResult.CurrentNode.(type) {
	case *CommentNode:
		assert.Equal(t, "# Second line", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// Section's comment
	nodeResult = quadlet.FindToken(NodePosition{3, 2})
	assert.Len(t, nodeResult.ParentNodes, 0)
	switch v := nodeResult.CurrentNode.(type) {
	case *CommentNode:
		assert.Equal(t, "# Section doc", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// [Container] section
	nodeResult = quadlet.FindToken(NodePosition{4, 3})
	require.NotNil(t, nodeResult.CurrentNode)
	assert.Len(t, nodeResult.ParentNodes, 0)
	switch v := nodeResult.CurrentNode.(type) {
	case *SectionNode:
		assert.Equal(t, "[Container]", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// Image's comment
	nodeResult = quadlet.FindToken(NodePosition{6, 2})
	require.NotNil(t, nodeResult.CurrentNode)
	require.Len(t, nodeResult.ParentNodes, 1)

	switch v := nodeResult.CurrentNode.(type) {
	case *CommentNode:
		assert.Equal(t, "# Image doc", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[0].(type) {
	case *SectionNode:
		assert.Equal(t, "[Container]", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// Image=foo.image Key
	nodeResult = quadlet.FindToken(NodePosition{7, 2})
	require.NotNil(t, nodeResult.CurrentNode)
	require.Len(t, nodeResult.ParentNodes, 1)

	switch v := nodeResult.CurrentNode.(type) {
	case *AssignNode:
		assert.Equal(t, "Image", *v.Name)
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[0].(type) {
	case *SectionNode:
		assert.Equal(t, "[Container]", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// Image=foo.image Value
	nodeResult = quadlet.FindToken(NodePosition{7, 8})
	require.NotNil(t, nodeResult.CurrentNode)
	require.Len(t, nodeResult.ParentNodes, 2)

	switch v := nodeResult.CurrentNode.(type) {
	case *ValueNode:
		assert.Equal(t, "foo.image", *v.Value)
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[0].(type) {
	case *AssignNode:
		assert.Equal(t, "Image", *v.Name)
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[1].(type) {
	case *SectionNode:
		assert.Equal(t, "[Container]", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// Label keyword
	nodeResult = quadlet.FindToken(NodePosition{10, 2})
	require.NotNil(t, nodeResult.CurrentNode)
	require.Len(t, nodeResult.ParentNodes, 1)

	switch v := nodeResult.CurrentNode.(type) {
	case *AssignNode:
		assert.Contains(t, "Label", *v.Name)
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[0].(type) {
	case *SectionNode:
		assert.Equal(t, "[Container]", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// Label multi-line value (first line)
	nodeResult = quadlet.FindToken(NodePosition{11, 2})
	require.NotNil(t, nodeResult.CurrentNode)
	require.Len(t, nodeResult.ParentNodes, 2)

	switch v := nodeResult.CurrentNode.(type) {
	case *ValueNode:
		assert.Contains(t, *v.Value, "env=test")
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[0].(type) {
	case *AssignNode:
		assert.Equal(t, "Label", *v.Name)
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[1].(type) {
	case *SectionNode:
		assert.Equal(t, "[Container]", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// [Unit] section
	nodeResult = quadlet.FindToken(NodePosition{14, 3})
	require.NotNil(t, nodeResult.CurrentNode)
	require.Len(t, nodeResult.ParentNodes, 0)

	switch v := nodeResult.CurrentNode.(type) {
	case *SectionNode:
		assert.Equal(t, "[Unit]", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// Description multi-line value (testing the "Foo" line)
	nodeResult = quadlet.FindToken(NodePosition{17, 3})
	require.NotNil(t, nodeResult.CurrentNode)
	require.Len(t, nodeResult.ParentNodes, 2)

	switch v := nodeResult.CurrentNode.(type) {
	case *ValueNode:
		assert.Equal(t, "Foo container", *v.Value)
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[0].(type) {
	case *AssignNode:
		assert.Equal(t, "Description", *v.Name)
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[1].(type) {
	case *SectionNode:
		assert.Equal(t, "[Unit]", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}

	// Description multi-line value (testing the "container" line)
	nodeResult = quadlet.FindToken(NodePosition{18, 3})
	require.NotNil(t, nodeResult.CurrentNode)
	require.Len(t, nodeResult.ParentNodes, 2)

	switch v := nodeResult.CurrentNode.(type) {
	case *ValueNode:
		assert.Equal(t, "Foo container", *v.Value)
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[0].(type) {
	case *AssignNode:
		assert.Equal(t, "Description", *v.Name)
	default:
		require.Fail(t, "invalid type")
	}

	switch v := nodeResult.ParentNodes[1].(type) {
	case *SectionNode:
		assert.Equal(t, "[Unit]", *v.Text)
	default:
		require.Fail(t, "invalid type")
	}
}

func Test_NodeQuadletString(t *testing.T) {
	input := QuadletNode{
		Documents: []*CommentNode{
			{
				Text: utils.AsPtr("# Test container"),
			},
		},
		Sections: []*SectionNode{
			{
				Text: utils.AsPtr("[Container]"),
				Assignments: []*AssignNode{
					{
						Name:      utils.AsPtr("Image"),
						Documents: nil,
						Value: &ValueNode{
							StartPos: NodePosition{2, 0},
							EndPos:   NodePosition{2, 12},
							Value:    utils.AsPtr("foo.image"),
						},
					},
					{
						Name:      utils.AsPtr("Label"),
						Documents: nil,
						Value: &ValueNode{
							StartPos: NodePosition{3, 6},
							EndPos:   NodePosition{3, 13},
							Value:    utils.AsPtr("env=test"),
						},
					},
				},
			},
			{
				Text: utils.AsPtr("[Unit]"),
				Assignments: []*AssignNode{
					{
						StartPos:  NodePosition{5, 0},
						EndPos:    NodePosition{4, 10},
						Name:      utils.AsPtr("Description"),
						Documents: nil,
						Value: &ValueNode{
							StartPos: NodePosition{5, 0},
							EndPos:   NodePosition{4, 12},
							Value:    utils.AsPtr("Foo container"),
						},
					},
				},
			},
		},
	}

	expected := `# Test container

[Container]
Image=foo.image
Label=env=test

[Unit]
Description=Foo container
`

	require.Equal(t, expected, input.String(), "invalid generated code")
}

func Test_NodeDocumentString(t *testing.T) {
	inputs := []CommentNode{
		{
			Text: utils.AsPtr("# Lorem ipsum"),
		},
		{
			Text: utils.AsPtr("Lorem ipsum"),
		},
	}

	expected := "# Lorem ipsum\n"

	for i, input := range inputs {
		require.Equal(t, expected, input.String(), "invalid output at %d", i)
	}
}

func Test_NodeSectionString(t *testing.T) {
	inputs := []SectionNode{
		{
			Text: utils.AsPtr("Container"),
			Documents: []*CommentNode{
				{
					Text: utils.AsPtr("This is container section"),
				},
			},
			Assignments: nil,
		},
		{
			Text: utils.AsPtr("[Container]"),
			Documents: []*CommentNode{
				{
					Text: utils.AsPtr("# This is container section"),
				},
			},
			Assignments: nil,
		},
	}

	expected := `# This is container section
[Container]
`

	for i, input := range inputs {
		require.Equal(t, expected, input.String(), "error at %d", i)
	}
}

func Test_NodeAssignString(t *testing.T) {
	inputs := []AssignNode{
		{
			Name: utils.AsPtr("Image"),
			Documents: []*CommentNode{
				{Text: utils.AsPtr("qsr-disable: qsr016")},
				{Text: utils.AsPtr("Image is in another file")},
			},
			Value: &ValueNode{
				Value: utils.AsPtr("foo.image"),
			},
		},
		{
			Name: utils.AsPtr("Image"),
			Documents: []*CommentNode{
				{Text: utils.AsPtr("# qsr-disable: qsr016")},
				{Text: utils.AsPtr("# Image is in another file")},
			},
			Value: &ValueNode{
				Value: utils.AsPtr("foo.image"),
			},
		},
	}

	expected := `
# qsr-disable: qsr016
# Image is in another file
Image=foo.image
`

	for i, input := range inputs {
		require.Equal(t, expected, input.String(), "error at %d", i)
	}
}

func Test_NodeAssignStringWithLongValue(t *testing.T) {
	input := AssignNode{
		Name: utils.AsPtr("Description"),
		Value: &ValueNode{
			Value: utils.AsPtr("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse egestas augue pretium lectus tristique interdum. Proin consequat rutrum ligula, sit amet pretium lectus feugiat malesuada"),
		},
		Documents: []*CommentNode{
			{Text: utils.AsPtr("# This is a hell long value")},
		},
	}

	expected := `
# This is a hell long value
Description=Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse egestas \
  augue pretium lectus tristique interdum. Proin consequat rutrum ligula, sit amet \
  pretium lectus feugiat malesuada
`

	require.Equal(t, expected, input.String(), "does not split the long lines")
}
