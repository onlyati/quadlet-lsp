package hover

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TestValueSecret tests the hover values on lines like 'Secret='.
func TestValueSecret(t *testing.T) {
	tmpDir := t.TempDir()
	content := `[Container]
Secret=foo-secret,type=env,target=FOO
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
		expected := "**Secret description**\n\nSyntax: `secret-name,option1=value1,option2=value2...\n\nThe `secret-name` is the name that has been created by `podman secret create` command.\n\nOptions:\n\n- **type=mount|env**: How the secret is exposed to the container. mount mounts the secret into the container as a file. env exposes the secret as an environment variable. Defaults to mount\n- **target=target-name**: Target of secret. For mounted secrets, this is the path to the secret inside the container. If a fully qualified path is provided, the secret is mounted at that location. Otherwise, the secret is mounted to /run/secrets/target for Linux containers or /var/run/secrets/target for FreeBSD containers. If the target is not set, the secret is mounted to /run/secrets/secretname by default. For env secrets, this is the environment variable key. Defaults to secretname.\n- **uid=n**: UID of secret. Defaults to 0. Mount secret type only.\n- **gid=n**: GID of secret. Defaults to 0. Mount secret type only.\n- **mode=0nnn**: Mode of secret. Defaults to 0444. Mount secret type only."
		assert.Equal(t, expected, v.Value, "unexpected content")
		assert.Equal(t, expectedStartPos.LineNumber, hoverValue.Range.Start.Line)
		assert.Equal(t, expectedStartPos.Position, hoverValue.Range.Start.Character)
		assert.Equal(t, expectedEndPos.Position, hoverValue.Range.End.Character)
		assert.Equal(t, expectedEndPos.LineNumber, hoverValue.Range.End.Line)
	default:
		assert.Fail(t, "hoverValue content is not protocol.MarkupContent")
	}
}
