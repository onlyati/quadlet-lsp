package lsp

import (
	"errors"
	"os"
	"path"
	"strings"

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
		"Volume":  "*.volume",
		"Pod":     "*.pod",
		"Network": "*.network",
		"Image":   "*.image",
	}

	if prop, ok := keywords[props[0]]; ok {
		return findQuadlets(prop, props[1])
	}

	return location, nil
}

// Just check if file exists
func findQuadlets(mask, value string) (protocol.Location, error) {
	var location protocol.Location

	if mask == "*.volume" {
		volParts := strings.Split(value, ":")
		value = volParts[0]
	}

	if strings.Contains(value, "@") {
		// If contains '@' then it is a systemd template
		value = convertTemplateNameToFile(value)
	}

	currDir, err := os.Getwd()
	if err != nil {
		return location, err
	}
	defPath := path.Join(currDir, value)

	if _, err := os.Stat(defPath); !errors.Is(err, os.ErrNotExist) {
		return protocol.Location{
			URI: protocol.DocumentUri("file://" + defPath),
		}, nil
	}

	return location, nil
}

// Convert template name like 'web@siteA.container' to 'web@.container'
func convertTemplateNameToFile(s string) string {
	atSign := strings.Index(s, "@")
	dotSign := strings.LastIndex(s, ".")

	return s[:atSign] + "@" + s[dotSign:]
}
