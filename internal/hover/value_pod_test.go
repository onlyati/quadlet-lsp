package hover

import (
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TestValuePodPeek tests when pod file exists that is currently hovered then
// display the content of the file.
func TestValuePodPeek(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.pod", "[Pod]\nPublishPort=8080:8080")
	info := HoverInformation{
		Line:              "Pod=foo.pod",
		URI:               "file://foo.container",
		CharacterPosition: uint32(4),
		RootDir:           tmpDir,
		Level:             2,
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
