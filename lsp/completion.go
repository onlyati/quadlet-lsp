package lsp

import (
	"fmt"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	uri := string(params.TextDocument.URI)
	text := documents.read(uri)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	editorLine := params.Position.Line
	charPos := params.Position.Character

	// This is a [] section, gives options
	if strings.HasPrefix(lines[editorLine], "[") {
		return listSectionCompletions(), nil
	}

	// Check if property already written and cursor after a '='
	// Then provides some options, depends what is the name of property
	if strings.Contains(lines[editorLine], "=") {
		return listPropertyParameter(lines, editorLine, charPos), nil
	}

	// Looking for possible properties in different sections
	return listPropertyCompletions(lines, editorLine), nil
}

func listSectionCompletions() []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem
	for k := range propertiesMap() {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: k,
		})
	}

	return completionItems
}

func listPropertyCompletions(lines []string, lineNumber protocol.UInteger) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	section := findSection(lines, lineNumber)

	// If no section at all, assume a new empty file
	// advice some basic template
	if section == "" {
		insertFormat := protocol.InsertTextFormatSnippet
		itemKind := protocol.CompletionItemKindSnippet

		for k, category := range defineCategoryProperties() {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label:            k,
				Detail:           category.details,
				InsertText:       category.insertText,
				InsertTextFormat: &insertFormat,
				Kind:             &itemKind,
			})
		}

		return completionItems
	}

	for _, prop := range propertiesMap()[section] {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: prop.label,
			Documentation: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**" + prop.label + "**\n\n" + strings.Join(prop.hover, "\n"),
			},
		})
	}

	return completionItems
}

func listPropertyParameter(lines []string, lineNumber protocol.UInteger, charPos protocol.UInteger) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	section := findSection(lines, lineNumber)

	if section == "" {
		return completionItems
	}

	property := strings.Split(lines[lineNumber], "=")[0]

	if property == "Image" {
		images, err := listImages()
		if err != nil {
			fmt.Printf("failed to execute command: %s", err.Error())
			return completionItems
		}
		return images
	}

	if property == "Secret" {
		secrets, err := listSecrets(
			lines[lineNumber][:charPos],
		)
		if err != nil {
			fmt.Printf("failed to execute command: %s", err.Error())
			return completionItems
		}
		return secrets
	}

	if property == "Volume" {
		volumes, err := listVolumes(
			lines[lineNumber][:charPos],
		)
		if err != nil {
			fmt.Printf("failed to list volmues: %s", err.Error())
			return completionItems
		}
		return volumes
	}

	for _, p := range propertiesMap()[section] {
		if property == p.label {
			for _, parm := range p.parameters {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: parm,
				})
			}
			return completionItems
		}
	}

	return completionItems
}

func listVolumes(line string) ([]protocol.CompletionItem, error) {
	var completionItems []protocol.CompletionItem

	props := strings.Split(line, "=")[1]

	if strings.Count(props, ":") > 0 {
		// Send parameters back
	}

	// List volumes from podman
	output, err := execPodmanCommand(
		[]string{"volume", "ls", "--format", "{{ .Name }}"},
	)
	if err != nil {
		return nil, err
	}
	for _, item := range output {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: item,
		})
	}

	// List *.volume files
	volFiles, err := listQuadletFiles("*.volume")
	if err != nil {
		return nil, err
	}
	completionItems = append(completionItems, volFiles...)

	return completionItems, nil
}

func listSecrets(line string) ([]protocol.CompletionItem, error) {
	var completionItems []protocol.CompletionItem

	props := strings.Split(line, "=")[1]

	if strings.Contains(props, ",") {
		completionItems = []protocol.CompletionItem{
			{
				Label: "type=mount",
			},
			{
				Label: "type=env",
			},
			{
				Label: "target=",
			},
		}

		return completionItems, nil
	}

	output, err := execPodmanCommand(
		[]string{"secret", "ls", "--format", "{{ .Name }}"},
	)
	if err != nil {
		return nil, err
	}
	for _, item := range output {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: item,
		})
	}

	return completionItems, nil
}

func listImages() ([]protocol.CompletionItem, error) {
	var completionItems []protocol.CompletionItem

	output, err := execPodmanCommand(
		[]string{"images", "--format", "{{ .Repository }}:{{ .Tag }}"},
	)
	if err != nil {
		return nil, err
	}
	for _, item := range output {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: item,
		})
	}

	return completionItems, nil
}
