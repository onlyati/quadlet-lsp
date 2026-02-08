package hover

import (
	"strings"
	"testing"

	"github.com/onlyati/quadlet-lsp/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// TestValueVolumePeek tests if hover contains the content of the volume that is
// hovered.
func TestValueVolumePeek(t *testing.T) {
	tmpDir := t.TempDir()

	testutils.CreateTempFile(t, tmpDir, "foo.volume", "[Volume]\n")

	info := HoverInformation{
		Line:              "Volume=foo.volume:/app:rw",
		URI:               "file://foo.container",
		Section:           "Container",
		CharacterPosition: uint32(7),
		RootDir:           tmpDir,
		Level:             2,
	}
	expectedMessage := []string{
		"**Host directory or source volume**",
		"",
		"If a volume source is specified, it must be a path on the host or the name of a named volume. Host paths are allowed to be absolute or relative; relative paths are resolved relative to the directory Podman is run in. If the source does not exist, Podman returns an error. Users must pre-create the source files or directories.",
		"",
		"Any source that does not begin with a `.` or `/` is treated as the name of a named volume. If a volume with that name does not exist, it is created. Volumes created with names are not anonymous, and they are not removed by the `--rm` option and the podman rm `--volumes` command.",
		"",
		"**Content of file**",
		"```quadlet",
		"[Volume]",
		"",
		"```",
	}
	hoverValue := HoverFunction(info)

	assert.NotNil(t, hoverValue, "return nil hover value")

	switch v := hoverValue.Contents.(type) {
	case protocol.MarkupContent:
		expected := strings.Join(expectedMessage, "\n")
		assert.Equal(t, expected, v.Value, "unexpected content")
	default:
		assert.Fail(t, "hoverValue content is not protocol.MarkupContent")
	}
}

// TestValueVolumeSource tests hover on the source volume.
func TestValueVolumeSource(t *testing.T) {
	cases := []HoverInformation{
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(7),
		},
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(8),
		},
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(20),
		},
	}

	for i, info := range cases {
		hoverValue := HoverFunction(info)

		require.NotNilf(t, hoverValue, "expected hover at %d", i)

		highlight := info.Line[hoverValue.Range.Start.Character:hoverValue.Range.End.Character]
		assert.Equal(t, "/home/ati/tmp", highlight, "unexpected highlight")
	}
}

// TestValueVolumeContainer tests the hover on the container side of volume
// definition.
func TestValueVolumeContainer(t *testing.T) {
	cases := []HoverInformation{
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(22),
		},
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(23),
		},
		{
			Line:              "Volume=/home/ati/tmp:/app/tmp:ro",
			URI:               "file://test.container",
			Section:           "Container",
			CharacterPosition: uint32(29),
		},
	}

	for i, info := range cases {
		hoverValue := HoverFunction(info)

		require.NotNilf(t, hoverValue, "expected hover at %d", i)

		highlight := info.Line[hoverValue.Range.Start.Character:hoverValue.Range.End.Character]
		assert.Equal(t, "/app/tmp", highlight, "unexpected highlight")
	}
}

// TestValueVolumeFlag tests hover on volume flags.
func TestValueVolumeFlag(t *testing.T) {
	info := HoverInformation{
		Line:    "Volume=/home/ati/tmp:/app/tmp:rw,z,U,nocopy,shared",
		URI:     "file://test.container",
		Section: "Container",
	}
	cases := []struct {
		position uint32
		flag     string
	}{
		{uint32(30), "rw"},
		{uint32(31), "rw"},
		{uint32(33), "z"},
		{uint32(35), "U"},
		{uint32(37), "nocopy"},
		{uint32(45), "shared"},
	}

	for i, c := range cases {
		info.CharacterPosition = c.position
		hoverValue := HoverFunction(info)

		require.NotNilf(t, hoverValue, "expected hovere at %d", i)

		highlight := info.Line[hoverValue.Range.Start.Character:hoverValue.Range.End.Character]
		assert.Equal(t, c.flag, highlight, "unexpected highlight")
	}
}
