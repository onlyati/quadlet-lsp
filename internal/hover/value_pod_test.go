package hover

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TestValuePodPeek tests when pod file exists that is currently hovered then
// display the content of the file.
func TestValuePodPeek(t *testing.T) {
	tmpDir := t.TempDir()
	content := `[Container]
Pod=foo.pod
`
	testutils.CreateTempFile(t, tmpDir, "foo.container", content)
	testutils.CreateTempFile(t, tmpDir, "foo.pod", "[Pod]\nPublishPort=8080:8080")
	p := parser.NewParser(path.Join(tmpDir, "foo.container"))
	tokenInfo := p.Quadlet.FindToken(parser.NodePosition{
		LineNumber: 1,
		Position:   9,
	})

	info := HoverInformation{
		RootDir:   tmpDir,
		Level:     2,
		TokenInfo: tokenInfo,
	}

	hoverValue := HoverFunction(info)

	assert.NotNil(t, hoverValue, "return nil hover value")

	switch v := hoverValue.Contents.(type) {
	case protocol.MarkupContent:
		expected := "**Content of file**\n```quadlet\n[Pod]\nPublishPort=8080:8080\n```"
		assert.Equal(t, expected, v.Value, "unexpected content")
	default:
		assert.Fail(t, "hoverValue content is not protocol.MarkupContent")
	}
}

func TestValuePodPeekWithTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	content := `[Container]
Pod=foo@hello.pod
`
	testutils.CreateTempFile(t, tmpDir, "foo.container", content)
	testutils.CreateTempFile(t, tmpDir, "foo@.pod", "[Pod]\nPublishPort=8080:8080")
	p := parser.NewParser(path.Join(tmpDir, "foo.container"))
	tokenInfo := p.Quadlet.FindToken(parser.NodePosition{
		LineNumber: 1,
		Position:   9,
	})

	info := HoverInformation{
		RootDir:   tmpDir,
		Level:     2,
		TokenInfo: tokenInfo,
	}

	hoverValue := HoverFunction(info)

	assert.NotNil(t, hoverValue, "return nil hover value")

	switch v := hoverValue.Contents.(type) {
	case protocol.MarkupContent:
		expected := "**Content of file**\n```quadlet\n[Pod]\nPublishPort=8080:8080\n```"
		assert.Equal(t, expected, v.Value, "unexpected content")
	default:
		assert.Fail(t, "hoverValue content is not protocol.MarkupContent")
	}
}

func TestValuePodPeekWithMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	content := `[Container]
Pod=foo.pod
`
	testutils.CreateTempFile(t, tmpDir, "foo.container", content)
	p := parser.NewParser(path.Join(tmpDir, "foo.container"))
	tokenInfo := p.Quadlet.FindToken(parser.NodePosition{
		LineNumber: 1,
		Position:   9,
	})

	info := HoverInformation{
		RootDir:   tmpDir,
		Level:     2,
		TokenInfo: tokenInfo,
	}

	hoverValue := HoverFunction(info)

	assert.Nil(t, hoverValue, "return nil hover value")
}
