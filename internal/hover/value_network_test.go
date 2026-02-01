package hover

import (
	"testing"

	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestValueNetworkPeek(t *testing.T) {
	tmpDir := t.TempDir()

	createTempFile(t, tmpDir, "foo.network", "[Network]\nDNS=1.1.1.1")
	info := HoverInformation{
		Line:              "Network=foo.network",
		URI:               "file://foo.container",
		CharacterPosition: uint32(8),
		RootDir:           tmpDir,
		Level:             2,
	}

	hoverValue := HoverFunction(info)

	assert.NotNil(t, hoverValue, "return nil hover value")

	switch v := hoverValue.Contents.(type) {
	case protocol.MarkupContent:
		expected := "**Content of file**\n```quadlet\n[Network]\nDNS=1.1.1.1\n```"
		assert.Equal(t, expected, v.Value, "unexpected content")
	default:
		t.Fatal("hoverValue content is not protocol.MarkupContent")
	}
}
