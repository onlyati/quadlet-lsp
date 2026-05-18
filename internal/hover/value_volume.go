package hover

import (
	"os"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func handleValueVolume(info HoverInformation) *protocol.Hover {
	value, ok := info.TokenInfo.CurrentNode.(*parser.ValueNode)
	if !ok {
		return nil
	}

	hoverData := []string{
		"**Volume description**",
		"",
	}

	parts := strings.Split(*value.Value, ":")
	hostVolume := parts[0]
	if strings.HasSuffix(hostVolume, ".volume") {
		// There is a link to another Quadlet, read it and display
		if strings.Contains(hostVolume, "@") {
			hostVolume = utils.ConvertTemplateNameToFile(hostVolume)
		}

		files, err := utils.ListQuadletFiles("volume", info.RootDir, info.Level)
		if err == nil {
			for _, f := range files {
				if hostVolume == f.Label {
					switch v := f.Documentation.(type) {
					case string:
						path, _ := strings.CutPrefix(v, "From work directory: ")
						content, err := os.ReadFile(path)
						if err == nil {
							hoverData = append(hoverData, "**Content of file**")
							hoverData = append(hoverData, "```quadlet")
							hoverData = append(hoverData, strings.Split(string(content), "\n")...)
							hoverData = append(hoverData, "```")
							hoverData = append(hoverData, "")
						}
					}
				}
			}
		}
	}

	// Add the regular documentation
	hoverData = append(hoverData, "**Documentation**")
	hoverData = append(hoverData, "")
	hoverData = append(hoverData, "Syntax: `source-volume:dest-volume:flag1,flag2,flag3...`")
	hoverData = append(hoverData, "")
	hoverData = append(hoverData, "If a volume source is specified, it must be a path on the host or the name of a named volume. Host paths are allowed to be absolute or relative; relative paths are resolved relative to the directory Podman is run in. If the source does not exist, Podman returns an error. Users must pre-create the source files or directories.")
	hoverData = append(hoverData, "Any source that does not begin with a `.` or `/` is treated as the name of a named volume. If a volume with that name does not exist, it is created. Volumes created with names are not anonymous, and they are not removed by the `--rm` option and the podman rm `--volumes` command.")
	hoverData = append(hoverData, "")
	hoverData = append(hoverData, "`source-volume`: Host directory or source volume.")
	hoverData = append(hoverData, "")
	hoverData = append(hoverData, "`dest-volume:`: Container directory. The container-dir must be an absolute path such as `/src/docs`. The volume is mounted into the container at this directory.")
	hoverData = append(hoverData, "")
	hoverData = append(hoverData, "**Frequently used flags**")
	hoverData = append(hoverData, "- `rw|ro`: Add `:ro` or `:rw` option to mount a volume in read-only or read-write mode, respectively. By default, the volumes are mounted read-write.")
	hoverData = append(hoverData, "- `U`: The `:U` suffix tells Podman to use the correct host UID and GID based on the UID and GID within the container, to change recursively the owner and group of the source volume. Chowning walks the file system under the volume and changes the UID/GID on each file. If the volume has thousands of inodes, this process takes a long time, delaying the start of the container.")
	hoverData = append(hoverData, "- `z|Z`: To change volume labels in Podman, add the :z suffix for shared content or :Z for private content. Large volumes may delay container startup because Podman must relabel every file, though it optimizes this if :z was previously used. For sensitive system areas or home directories, avoid relabeling and instead disable SELinux separation using --security-opt label=disable to prevent system-wide permission failures.z|Z`: To change volume labels in Podman, add the :z suffix for shared content or :Z for private content. Large volumes may delay container startup because Podman must relabel every file, though it optimizes this if :z was previously used. For sensitive system areas or home directories, avoid relabeling and instead disable SELinux separation using --security-opt label=disable to prevent system-wide permission failures.")
	hoverData = append(hoverData, "")
	hoverData = append(hoverData, "For more flag description check out [Podman run documentation](https://docs.podman.io/en/latest/markdown/podman-run.1.html#volume-v-source-volume-host-dir-container-dir-options)")

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: strings.Join(hoverData, "\n"),
		},
		Range: &protocol.Range{
			Start: protocol.Position{
				Line:      value.StartPos.LineNumber,
				Character: value.StartPos.Position,
			},
			End: protocol.Position{
				Line:      value.EndPos.LineNumber,
				Character: value.EndPos.Position,
			},
		},
	}
}
