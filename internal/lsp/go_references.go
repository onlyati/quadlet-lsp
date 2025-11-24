package lsp

import (
	"errors"
	"os"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type goReferenceProperty struct {
	property string
	searchIn []string
}

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
	keywords := map[string]goReferenceProperty{
		"[Volume]": {
			property: "Volume",
			searchIn: []string{"container", "pod"},
		},
		"[Network]": {
			property: "Network",
			searchIn: []string{"container", "pod", "kube"},
		},
		"[Pod]": {
			property: "Pod",
			searchIn: []string{"container", "kube", "volume", "network", "image", "build"},
		},
		"[Image]": {
			property: "Image",
			searchIn: []string{"container"},
		},
	}

	config.Mu.RLock()
	rootDir := config.WorkspaceRoot
	config.Mu.RUnlock()
	if prop, ok := keywords[currentLine]; ok {
		return findReferences(prop, currentFileName, rootDir)
	}

	return locations, nil
}

func findReferences(prop goReferenceProperty, currentFileName, rootDir string) ([]protocol.Location, error) {
	var locations []protocol.Location
	files := []protocol.CompletionItem{}

	for _, d := range prop.searchIn {
		filesTmp, err := utils.ListQuadletFiles(d, rootDir)
		if err != nil {
			return nil, err
		}
		files = append(files, filesTmp...)
	}

	for _, f := range files {
		p := ""
		switch v := f.Documentation.(type) {
		case string:
			p = v
		default:
			return nil, errors.New("unexpected error: documentation is not string")
		}

		p, _ = strings.CutPrefix(p, "From work directory: ")
		content, err := os.ReadFile(p)
		if err != nil {
			return nil, err
		}

		fname := p[strings.LastIndex(p, string(os.PathSeparator))+1:]
		section := utils.FirstCharacterToUpper(fname[strings.LastIndex(fname, ".")+1:])

		items := utils.FindItems(utils.FindItemProperty{
			URI:           p,
			RootDirectory: rootDir,
			Text:          string(content),
			Section:       "[" + section + "]",
			Property:      prop.property,
		})
		for _, item := range items {
			if prop.property == "Volume" {
				volParts := strings.Split(item.Value, ":")
				item.Value = volParts[0]
			}
			if strings.Contains(item.Value, "@") {
				// If contains '@' then it is a systemd template
				item.Value = utils.ConvertTemplateNameToFile(item.Value)
			}
			if item.Value == currentFileName {
				uri := p
				if item.FilePath != "" {
					uri = item.FilePath
				}
				locations = append(locations, protocol.Location{
					URI: protocol.DocumentUri("file://" + uri),
					Range: protocol.Range{
						Start: protocol.Position{Line: item.LineNumber, Character: 0},
						End:   protocol.Position{Line: item.LineNumber, Character: item.Length},
					},
				})
			}
		}
	}

	return locations, nil
}
