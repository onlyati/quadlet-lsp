package hover

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// When hover is displayed on network and file exists, then the file's content
// also displayed.
func TestValueNetworkPeek(t *testing.T) {
	tmpDir := t.TempDir()
	content := `[Container]
Network=foo.network
`
	testutils.CreateTempFile(t, tmpDir, "foo.container", content)
	testutils.CreateTempFile(t, tmpDir, "foo.network", "[Network]\nDNS=1.1.1.1")
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
		expected := "**Content of file**\n```quadlet\n[Network]\nDNS=1.1.1.1\n```"
		assert.Equal(t, expected, v.Value, "unexpected content")
	default:
		assert.Fail(t, "hoverValue content is not protocol.MarkupContent")
	}
}

// When hover is displayed on network and file does not exists, display nothing.
func TestValueNetworkPeekMissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	content := `[Container]
Network=foo.network
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

// Test with template
func TestValueNetworkPeekWithTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	content := `[Container]
Network=foo@hello.network
`
	testutils.CreateTempFile(t, tmpDir, "foo.container", content)
	testutils.CreateTempFile(t, tmpDir, "foo@.network", "[Network]\nDNS=1.1.1.1")
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
		expected := "**Content of file**\n```quadlet\n[Network]\nDNS=1.1.1.1\n```"
		assert.Equal(t, expected, v.Value, "unexpected content")
	default:
		assert.Fail(t, "hoverValue content is not protocol.MarkupContent")
	}
}
