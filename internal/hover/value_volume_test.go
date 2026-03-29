package hover

import (
	"path"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	"github.com/stretchr/testify/assert"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TestValueVolumePeek tests if hover contains the content of the volume that is
// hovered.
func TestValueVolumePeek(t *testing.T) {
	tmpDir := t.TempDir()
	content := `[Container]
Volume=foo.volume:/app/:U,z
`
	testutils.CreateTempFile(t, tmpDir, "foo.container", content)
	testutils.CreateTempFile(t, tmpDir, "foo.volume", "[Volume]")
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
		expected := "**Volume description**\n\n**Content of file**\n```quadlet\n[Volume]\n```\n\n**Documentation**\n\nSyntax: `source-volume:dest-volume:flag1,flag2,flag3...`\n\nIf a volume source is specified, it must be a path on the host or the name of a named volume. Host paths are allowed to be absolute or relative; relative paths are resolved relative to the directory Podman is run in. If the source does not exist, Podman returns an error. Users must pre-create the source files or directories.\nAny source that does not begin with a `.` or `/` is treated as the name of a named volume. If a volume with that name does not exist, it is created. Volumes created with names are not anonymous, and they are not removed by the `--rm` option and the podman rm `--volumes` command.\n\n`source-volume`: Host directory or source volume.\n\n`dest-volume:`: Container directory. The container-dir must be an absolute path such as `/src/docs`. The volume is mounted into the container at this directory.\n\n**Frequently used flags**\n- `rw|ro`: Add `:ro` or `:rw` option to mount a volume in read-only or read-write mode, respectively. By default, the volumes are mounted read-write.\n- `U`: The `:U` suffix tells Podman to use the correct host UID and GID based on the UID and GID within the container, to change recursively the owner and group of the source volume. Chowning walks the file system under the volume and changes the UID/GID on each file. If the volume has thousands of inodes, this process takes a long time, delaying the start of the container.\n- `z|Z`: To change volume labels in Podman, add the :z suffix for shared content or :Z for private content. Large volumes may delay container startup because Podman must relabel every file, though it optimizes this if :z was previously used. For sensitive system areas or home directories, avoid relabeling and instead disable SELinux separation using --security-opt label=disable to prevent system-wide permission failures.z|Z`: To change volume labels in Podman, add the :z suffix for shared content or :Z for private content. Large volumes may delay container startup because Podman must relabel every file, though it optimizes this if :z was previously used. For sensitive system areas or home directories, avoid relabeling and instead disable SELinux separation using --security-opt label=disable to prevent system-wide permission failures.\n\nFor more flag description check out [Podman run documentation](https://docs.podman.io/en/latest/markdown/podman-run.1.html#volume-v-source-volume-host-dir-container-dir-options)"
		assert.Equal(t, expected, v.Value, "unexpected content")
		assert.Equal(t, expectedStartPos.LineNumber, hoverValue.Range.Start.Line)
		assert.Equal(t, expectedStartPos.Position, hoverValue.Range.Start.Character)
		assert.Equal(t, expectedEndPos.Position, hoverValue.Range.End.Character)
		assert.Equal(t, expectedEndPos.LineNumber, hoverValue.Range.End.Line)
	default:
		assert.Fail(t, "hoverValue content is not protocol.MarkupContent")
	}
}

func TestValueVolumePeekWithoutFile(t *testing.T) {
	tmpDir := t.TempDir()
	content := `[Container]
Volume=foo.volume:/app/:U,z
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
		expected := "**Volume description**\n\n**Documentation**\n\nSyntax: `source-volume:dest-volume:flag1,flag2,flag3...`\n\nIf a volume source is specified, it must be a path on the host or the name of a named volume. Host paths are allowed to be absolute or relative; relative paths are resolved relative to the directory Podman is run in. If the source does not exist, Podman returns an error. Users must pre-create the source files or directories.\nAny source that does not begin with a `.` or `/` is treated as the name of a named volume. If a volume with that name does not exist, it is created. Volumes created with names are not anonymous, and they are not removed by the `--rm` option and the podman rm `--volumes` command.\n\n`source-volume`: Host directory or source volume.\n\n`dest-volume:`: Container directory. The container-dir must be an absolute path such as `/src/docs`. The volume is mounted into the container at this directory.\n\n**Frequently used flags**\n- `rw|ro`: Add `:ro` or `:rw` option to mount a volume in read-only or read-write mode, respectively. By default, the volumes are mounted read-write.\n- `U`: The `:U` suffix tells Podman to use the correct host UID and GID based on the UID and GID within the container, to change recursively the owner and group of the source volume. Chowning walks the file system under the volume and changes the UID/GID on each file. If the volume has thousands of inodes, this process takes a long time, delaying the start of the container.\n- `z|Z`: To change volume labels in Podman, add the :z suffix for shared content or :Z for private content. Large volumes may delay container startup because Podman must relabel every file, though it optimizes this if :z was previously used. For sensitive system areas or home directories, avoid relabeling and instead disable SELinux separation using --security-opt label=disable to prevent system-wide permission failures.z|Z`: To change volume labels in Podman, add the :z suffix for shared content or :Z for private content. Large volumes may delay container startup because Podman must relabel every file, though it optimizes this if :z was previously used. For sensitive system areas or home directories, avoid relabeling and instead disable SELinux separation using --security-opt label=disable to prevent system-wide permission failures.\n\nFor more flag description check out [Podman run documentation](https://docs.podman.io/en/latest/markdown/podman-run.1.html#volume-v-source-volume-host-dir-container-dir-options)"
		assert.Equal(t, expected, v.Value, "unexpected content")
		assert.Equal(t, expectedStartPos.LineNumber, hoverValue.Range.Start.Line)
		assert.Equal(t, expectedStartPos.Position, hoverValue.Range.Start.Character)
		assert.Equal(t, expectedEndPos.Position, hoverValue.Range.End.Character)
		assert.Equal(t, expectedEndPos.LineNumber, hoverValue.Range.End.Line)
	default:
		assert.Fail(t, "hoverValue content is not protocol.MarkupContent")
	}
}
