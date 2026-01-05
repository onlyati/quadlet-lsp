package lsp

import (
	"os"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
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
	text := documents.Read(uri)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	editorLine := params.Position.Line
	currentLine := strings.TrimSpace(lines[editorLine])

	pathParts := strings.Split(uri, string(os.PathSeparator))
	currentFileName := pathParts[len(pathParts)-1]

	// It only work when user cursor in the line of the section title
	// Probably there is a cleaner way, but it works
	keywords := map[string]utils.GoReferenceProperty{
		"[Volume]": {
			Property: "Volume",
			SearchIn: []string{"container", "pod"},
			DirLevel: *config.Project.DirLevel,
		},
		"[Network]": {
			Property: "Network",
			SearchIn: []string{"container", "pod", "kube"},
			DirLevel: *config.Project.DirLevel,
		},
		"[Pod]": {
			Property: "Pod",
			SearchIn: []string{"container", "kube", "volume", "network", "image", "build"},
			DirLevel: *config.Project.DirLevel,
		},
		"[Image]": {
			Property: "Image",
			SearchIn: []string{"container"},
			DirLevel: *config.Project.DirLevel,
		},
	}

	config.Mu.RLock()
	rootDir := config.WorkspaceRoot
	config.Mu.RUnlock()
	if prop, ok := keywords[currentLine]; ok {
		return utils.FindReferences(prop, currentFileName, rootDir)
	}

	return locations, nil
}
