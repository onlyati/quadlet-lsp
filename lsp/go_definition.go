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

	prop := strings.Split(currentLine, "=")

	if len(prop) < 2 {
		return nil, nil
	}

	currDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if prop[0] == "Volume" {
		volumes, err := listQuadletFiles("*.volume")
		if err != nil {
			return nil, err
		}

		volParts := strings.Split(prop[1], ":")
		volName := volParts[0]

		for _, vol := range volumes {
			if volName == vol.Label {
				return protocol.Location{
					URI: protocol.DocumentUri("file://" + currDir + string(os.PathSeparator) + vol.Label),
				}, nil
			}
		}

		return location, nil
	}

	if prop[0] == "Pod" {
		pods, err := listQuadletFiles("*.pod")
		if err != nil {
			return nil, err
		}

		for _, pod := range pods {
			if prop[1] == pod.Label {
				return protocol.Location{
					URI: protocol.DocumentUri("file://" + currDir + string(os.PathSeparator) + pod.Label),
				}, nil
			}
		}

		return location, nil
	}

	if prop[0] == "Network" {
		networks, err := listQuadletFiles("*.network")
		if err != nil {
			return nil, err
		}

		for _, network := range networks {
			if prop[1] == network.Label {
				return protocol.Location{
					URI: protocol.DocumentUri("file://" + currDir + string(os.PathSeparator) + network.Label),
				}, nil
			}
		}

		return location, nil
	}

	return location, nil
}
