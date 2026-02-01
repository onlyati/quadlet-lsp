package hover

import (
	"testing"

	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func TestValuePodPeek(t *testing.T) {
	tmpDir := t.TempDir()

	createTempFile(t, tmpDir, "foo.pod", "[Pod]\nPublishPort=8080:8080")
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
		t.Fatal("hoverValue content is not protocol.MarkupContent")
	}
}
