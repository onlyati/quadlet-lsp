package lsp

import (
	"os"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textDefinition(
	context *glsp.Context,
	params *protocol.DefinitionParams,
) (any, error) {
	var location protocol.Location

	uri := string(params.TextDocument.URI)
	text := documents.read(uri)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	editorLine := params.Position.Line
	currentLine := lines[editorLine]

	props := strings.Split(currentLine, "=")

	if len(props) < 2 {
		return nil, nil
	}

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

func findQuadlets(mask, value string) (protocol.Location, error) {
	var location protocol.Location

	currDir, err := os.Getwd()
	if err != nil {
		return location, err
	}

	files, err := listQuadletFiles(mask)
	if err != nil {
		return location, err
	}

	if mask == "*.volume" {
		volParts := strings.Split(value, ":")
		value = volParts[0]
	}

	for _, file := range files {
		if value == file.Label {
			return protocol.Location{
				URI: protocol.DocumentUri("file://" + currDir + string(os.PathSeparator) + file.Label),
			}, nil
		}
	}

	return location, nil
}
