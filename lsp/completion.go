package lsp

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func textCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	executor := CommandExecutor{}
	uri := string(params.TextDocument.URI)
	text := documents.read(uri)
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	editorLine := params.Position.Line
	charPos := params.Position.Character

	// This is a [] section, gives options
	if strings.HasPrefix(lines[editorLine], "[") {
		return listSectionCompletions(), nil
	}

	// Check if newSomething macro is typed
	if strings.HasPrefix(lines[editorLine], "new.") {
		return listNewMacros(lines, editorLine), nil
	}

	// Check if property already written and cursor after a '='
	// Then provides some options, depends what is the name of property
	if strings.Contains(lines[editorLine], "=") {
		return listPropertyParameter(executor, lines, editorLine, charPos), nil
	}

	// Looking for possible properties in different sections
	return listPropertyCompletions(lines, editorLine), nil
}

func listSectionCompletions() []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem
	for k := range propertiesMap {
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

		for k, category := range categoryProperty {
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

	for _, prop := range propertiesMap[section] {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: prop.label + "=",
			Documentation: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**" + prop.label + "**\n\n" + strings.Join(prop.hover, "\n"),
			},
		})
	}

	return completionItems
}

func listNewMacros(lines []string, lineNumber protocol.UInteger) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem
	section := findSection(lines, lineNumber)

	insertFormat := protocol.InsertTextFormatSnippet
	itemKind := protocol.CompletionItemKindSnippet

	lineText := lines[lineNumber]

	// Try to find the character range of "new." (if present)
	if !strings.HasPrefix(lineText, "new.") {
		return completionItems
	}

	// Get the rest of the line after "new." prefix
	propName := strings.TrimPrefix(lineText, "new.")

	// We'll replace from position 0 to len("new."+propName)
	startChar := 0
	endChar := len("new." + propName)

	for _, p := range propertiesMap[section] {
		if strings.HasPrefix(p.label, propName) && p.macro != "" {
			textEdit := protocol.TextEdit{
				Range: protocol.Range{
					Start: protocol.Position{Line: lineNumber, Character: uint32(startChar)},
					End:   protocol.Position{Line: lineNumber, Character: uint32(endChar)},
				},
				NewText: p.macro,
			}

			completionItems = append(completionItems, protocol.CompletionItem{
				Label:            "new." + p.label,
				Kind:             &itemKind,
				TextEdit:         &textEdit,
				InsertTextFormat: &insertFormat,
			})
		}
	}

	return completionItems
}

func listPropertyParameter(c Commander, lines []string, lineNumber protocol.UInteger, charPos protocol.UInteger) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	property := strings.Split(lines[lineNumber], "=")[0]

	if property == "Image" {
		images, err := listImages(c)
		if err != nil {
			fmt.Printf("failed to execute command: %s", err.Error())
			return completionItems
		}
		return images
	}

	if property == "Secret" {
		secrets, err := listSecrets(
			c,
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
			c,
			lines[lineNumber][:charPos],
		)
		if err != nil {
			fmt.Printf("failed to list volmues: %s", err.Error())
			return completionItems
		}
		return volumes
	}

	if property == "Pod" {
		pods, err := listQuadletFiles("*.pod")
		if err != nil {
			fmt.Printf("failed to list pods: %s", err.Error())
			return completionItems
		}
		return pods
	}

	if property == "Network" {
		networks, err := listNetworks(c)
		if err != nil {
			fmt.Printf("failed to list networks: %s", err.Error())
			return completionItems
		}
		return networks
	}

	if property == "PublishPort" {
		ports, err := listPublishedPorts(
			c,
			lines,
			lineNumber,
		)
		if err != nil {
			fmt.Printf("failed to list ports: %s", err.Error())
			return completionItems
		}
		return ports
	}

	section := findSection(lines, lineNumber)

	if section == "" {
		return completionItems
	}

	for _, p := range propertiesMap[section] {
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

// If user at the `PublihsPort=` line, and typting the exposed port number
// provide suggestions based on image inspect what ports can be exposed.
func listPublishedPorts(c Commander, lines []string, lineNumber protocol.UInteger) ([]protocol.CompletionItem, error) {
	var completionItems []protocol.CompletionItem

	// Let's find out that we need to provide any complation at all
	colons := strings.Count(lines[lineNumber], ":")
	tmp := strings.Split(lines[lineNumber], ":")

	// We need complation in two cases:
	// ExposedPorts=127.0.0.1:420:69
	// ExposedPorts=420:69
	if colons == 0 {
		return completionItems, nil
	}
	if colons == 1 {
		// Check if first part is an IP address
		if strings.Count(tmp[0], ".") > 0 {
			return completionItems, nil
		}
	}

	// First looking for `Image=value` value
	// First looing for reverse, people usually define image first then parameters
	imageName := ""
	for i := lineNumber; i > 0; i-- {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "Image=") {
			tmp := strings.Split(line, "=")
			if len(tmp) != 2 {
				return nil, errors.New("seems invalid value at `Image=` line")
			}
			imageName = tmp[1]
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			// We've reached the start of section, try in other direction
			break
		}
	}

	// Check rest of the file for `Image=`
	if imageName == "" {
		for i := lineNumber; int(i) < len(lines); i++ {
			line := strings.TrimSpace(lines[i])
			if strings.HasPrefix(line, "Image=") {
				tmp := strings.Split(line, "=")
				if len(tmp) != 2 {
					return nil, errors.New("seems invalid value at `Image=` line")
				}
				imageName = tmp[1]
			}
			if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
				// We've reached the start of another section
				break
			}
		}
	}

	// We've found something, let's check it
	if imageName != "" {
		output, err := c.Run(
			"podman",
			"image", "inspect", imageName,
		)
		if err != nil {
			return nil, err
		}
		inspectJSON := strings.Join(output, "")
		var data []map[string]any
		json.Unmarshal([]byte(inspectJSON), &data)

		config, ok := data[0]["Config"].(map[string]any)
		if !ok {
			return nil, err
		}

		exposedPorts, ok := config["ExposedPorts"].(map[string]any)
		if !ok {
			return nil, err
		}

		for port := range exposedPorts {
			tmp := strings.Split(port, "/")
			completionItems = append(completionItems, protocol.CompletionItem{
				Label: tmp[0],
			})
		}
	}

	return completionItems, nil
}

func listNetworks(c Commander) ([]protocol.CompletionItem, error) {
	var completionItems []protocol.CompletionItem

	// List networks from podman
	output, err := c.Run(
		"podman",
		"network", "ls", "--format", "{{ .Name }}",
	)
	if err != nil {
		return nil, err
	}
	for _, item := range output {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: item,
		})
	}

	// List *.network files
	volFiles, err := listQuadletFiles("*.network")
	if err != nil {
		return nil, err
	}
	completionItems = append(completionItems, volFiles...)
	return completionItems, nil
}

func listVolumes(c Commander, line string) ([]protocol.CompletionItem, error) {
	var completionItems []protocol.CompletionItem

	props := strings.Split(line, "=")[1]

	if strings.Count(props, ":") == 1 {
		return completionItems, nil
	}

	if strings.Count(props, ":") == 2 {
		// Send volume options back
		opts := []string{"rw", "ro", "z", "Z"}
		for _, opt := range opts {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label: opt,
			})
		}
		return completionItems, nil
	}

	// List volumes from podman
	output, err := c.Run(
		"podman",
		"volume", "ls", "--format", "{{ .Name }}",
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

func listSecrets(c Commander, line string) ([]protocol.CompletionItem, error) {
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

	output, err := c.Run(
		"podman",
		"secret", "ls", "--format", "{{ .Name }}",
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

func listImages(c Commander) ([]protocol.CompletionItem, error) {
	var completionItems []protocol.CompletionItem

	output, err := c.Run(
		"podman",
		"images", "--format", "{{ .Repository }}:{{ .Tag }}",
	)
	if err != nil {
		return nil, err
	}
	for _, item := range output {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: item,
		})
	}

	files, err := listQuadletFiles("*.image")
	if err != nil {
		return completionItems, err
	}
	completionItems = append(completionItems, files...)

	return completionItems, nil
}
