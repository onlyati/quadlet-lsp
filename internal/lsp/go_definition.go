package lsp

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// If there is a link for a file, e.g.: `Volume=something.volume`,
// it looking for this file in the current work directory then send
// the URI back to the editor.
func textDefinition(
	context *glsp.Context,
	params *protocol.DefinitionParams,
) (any, error) {
	var location protocol.Location

	uri := string(params.TextDocument.URI)
	text := documents.Read(uri)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	editorLine := params.Position.Line
	currentLine := lines[editorLine]

	props := strings.Split(currentLine, "=")

	if len(props) < 2 {
		return nil, nil
	}

	// Depends on which line the cursor stand (what is before '=')
	// looking for different file extension, then find the file.
	// Probably there is a cleaner way, but it works.
	keywords := map[string]string{
		"Volume":  "volume",
		"Pod":     "pod",
		"Network": "network",
		"Image":   "image",
	}

	config.Mu.RLock()
	rootDir := config.WorkspaceRoot
	level := config.Project.DirLevel
	config.Mu.RUnlock()
	if prop, ok := keywords[props[0]]; ok {
		return findQuadlets(prop, props[1], rootDir, *level)
	}

	return location, nil
}

// Just check if file exists
func findQuadlets(mask, value, rootDir string, level int) (protocol.Location, error) {
	var location protocol.Location

	if mask == "volume" {
		volParts := strings.Split(value, ":")
		value = volParts[0]
	}

	if strings.Contains(value, "@") {
		// If contains '@' then it is a systemd template
		value = utils.ConvertTemplateNameToFile(value)
	}

	files, err := utils.ListQuadletFiles(mask, rootDir, level)
	if err != nil {
		return location, err
	}

	for _, f := range files {
		if f.Label == value {
			p := ""
			switch v := f.Documentation.(type) {
			case string:
				p = v
			default:
				return location, nil
			}
			p, _ = strings.CutPrefix(p, "From work directory: ")
			return protocol.Location{
				URI: protocol.DocumentUri("file://" + p),
			}, nil
		}
	}

	return location, nil
}
