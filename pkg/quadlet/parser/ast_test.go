package parser

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/stretchr/testify/require"
)

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
