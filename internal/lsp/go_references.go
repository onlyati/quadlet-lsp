package lsp

import (
	"os"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// The go reference feature scan all file in the current directory
// and looking for when the current file is used.
func textReferences(
	context *glsp.Context,
	params *protocol.ReferenceParams,
) ([]protocol.Location, error) {
	var locations []protocol.Location

	uri := string(params.TextDocument.URI)
	text := documents.read(uri)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	editorLine := params.Position.Line
	currentLine := strings.TrimSpace(lines[editorLine])

	pathParts := strings.Split(uri, string(os.PathSeparator))
	currentFileName := pathParts[len(pathParts)-1]

	// It only work when user cursor in the line of the section title
	// Probably there is a cleaner way, but it works
	keywords := map[string]string{
		"[Volume]":  "Volume",
		"[Network]": "Network",
		"[Pod]":     "Pod",
		"[Image]":   "Image",
	}

	if prop, ok := keywords[currentLine]; ok {
		return findReferences(prop, currentFileName)
	}

	return locations, nil
}

func findReferences(prop, currentFileName string) ([]protocol.Location, error) {
	var locations []protocol.Location

	targetLine := prop + "=" + currentFileName
	locations, err := findLineStartWith(targetLine)
	if err != nil {
		return nil, err
	}
	return locations, nil
}
