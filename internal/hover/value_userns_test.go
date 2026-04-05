package hover

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TestValueUserNS tests hover output for the 'UserNS=' lines.
func TestValueUserNS(t *testing.T) {
	tmpDir := t.TempDir()
	content := `[Container]
UserNS=keep-id:uid=101,gid=101
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
	expectedStartPos := tokenInfo.CurrentNode.(*parser.ValueNode).StartPos
	expectedEndPos := tokenInfo.CurrentNode.(*parser.ValueNode).EndPos

	hoverValue := HoverFunction(info)

	assert.NotNil(t, hoverValue, "return nil hover value")

	switch v := hoverValue.Contents.(type) {
	case protocol.MarkupContent:
		expected := "**UserNS description**\n\n`auto`: Container user: nil (Host User UID is not mapped into container.)\n\nPodman allocates unique ranges of UIDs and GIDs from the containers subordinate user IDs. The size of the ranges is based on the number of UIDs required in the image. The number of UIDs and GIDs can be overridden with the size option.\n\n`host`: Container user: 0 (Default User account mapped to root user in container.)\n\nhost or “” (empty string): run in the user namespace of the caller. The processes running in the container have the same privileges on the host as any other process launched by the calling user.\n\n`keep-id`: creates a user namespace where the current user’s UID:GID are mapped to the same values in the container.\n\nFor containers created by root, the current mapping is created into a new user namespace.\nValid keep-id options:\n- uid=UID: override the UID inside the container that is used to map the current user to.\n- gid=GID: override the GID inside the container that is used to map the current user to.\n- size=SIZE: override the size of the configured user namespace. It is useful to not saturate all the available IDs. Not supported when running as root.\n\n`nomap`: Container user: nil (Host User UID is not mapped into container.)\n\nnomap: creates a user namespace where the current rootless user’s UID:GID are not mapped into the container. This option is not allowed for containers created by the root user."
		assert.Equal(t, expected, v.Value, "unexpected content")
		assert.Equal(t, expectedStartPos.LineNumber, hoverValue.Range.Start.Line)
		assert.Equal(t, expectedStartPos.Position, hoverValue.Range.Start.Character)
		assert.Equal(t, expectedEndPos.Position, hoverValue.Range.End.Character)
		assert.Equal(t, expectedEndPos.LineNumber, hoverValue.Range.End.Line)
	default:
		assert.Fail(t, "hoverValue content is not protocol.MarkupContent")
	}
}
