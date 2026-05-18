package hover

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestHandlePropertyHover(t *testing.T) {
	content := `[Container]
AppArmor=foo
`
	tmpDir := t.TempDir()
	testutils.CreateTempFile(t, tmpDir, "foo.container", content)
	p := parser.NewParser(path.Join(tmpDir, "foo.container"))
	tokenInfo := p.Quadlet.FindToken(parser.NodePosition{
		LineNumber: 1,
		Position:   3,
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
		expected := "**AppArmor**\n\nSets the apparmor confinement profile for the container. A value of `unconfined` turns off apparmor confinement."
		assert.Equal(t, expected, v.Value, "unexpected content")
	default:
		assert.Fail(t, "hoverValue content is not protocol.MarkupContent")
	}
}
