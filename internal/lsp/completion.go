package lsp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// This function handles the completion event that is received.
func textCompletion(context *glsp.Context, params *protocol.CompletionParams) (any, error) {
	executor := utils.CommandExecutor{}
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

// Cursor in a line that start with '[' character. Probably want to type
// section header like `[Network]`.
func listSectionCompletions() []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem
	for k := range data.PropertiesMap {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: k,
		})
	}

	return completionItems
}

// No special handling, just advice suggestion based on static
// data in `properties.go` file.
func listPropertyCompletions(lines []string, lineNumber protocol.UInteger) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	section := findSection(lines, lineNumber)

	// If no section at all, assume a new empty file
	// advice some basic template
	if section == "" {
		insertFormat := protocol.InsertTextFormatSnippet
		itemKind := protocol.CompletionItemKindSnippet

		for k, category := range data.CategoryProperty {
			completionItems = append(completionItems, protocol.CompletionItem{
				Label:            k,
				Detail:           category.Details,
				InsertText:       category.InsertText,
				InsertTextFormat: &insertFormat,
				Kind:             &itemKind,
			})
		}

		return completionItems
	}

	// It is a line where the '=' is not present, so probably just
	// want to type something, let give a hint.
	for _, prop := range data.PropertiesMap[section] {
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: prop.Label + "=",
			Documentation: protocol.MarkupContent{
				Kind:  protocol.MarkupKindMarkdown,
				Value: "**" + prop.Label + "**\n\n" + strings.Join(prop.Hover, "\n"),
			},
		})
	}

	return completionItems
}

// In the editor typed something like `new.Something`, when `Something`
// can be different thing like `Environment`, `Secret`, etc.
// This provide new templates based on `properties.go` file at `mask` attribute.
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

	for _, p := range data.PropertiesMap[section] {
		if strings.HasPrefix(p.Label, propName) && p.Macro != "" {
			textEdit := protocol.TextEdit{
				Range: protocol.Range{
					Start: protocol.Position{Line: lineNumber, Character: uint32(startChar)},
					End:   protocol.Position{Line: lineNumber, Character: uint32(endChar)},
				},
				NewText: p.Macro,
			}

			completionItems = append(completionItems, protocol.CompletionItem{
				Label:            "new." + p.Label,
				Kind:             &itemKind,
				TextEdit:         &textEdit,
				InsertTextFormat: &insertFormat,
			})
		}
	}

	return completionItems
}

// This function run when there is already an '=' sign in the
// line of the cursor. This provide dynamic and static completions.
func listPropertyParameter(c utils.Commander, lines []string, lineNumber protocol.UInteger, charPos protocol.UInteger) []protocol.CompletionItem {
	var completionItems []protocol.CompletionItem

	property := strings.Split(lines[lineNumber], "=")[0]

	// Looking for *.image files and check `podman images` command
	if property == "Image" {
		images, err := listImages(c)
		if err != nil {
			fmt.Printf("failed to execute command: %s", err.Error())
			return completionItems
		}
		return images
	}

	// Looking for `podman secret ls`
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

	// Looking for *.volume files and `podman volume ls`
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

	// Looking for *.pod files
	if property == "Pod" {
		pods, err := utils.ListQuadletFiles("*.pod")
		if err != nil {
			fmt.Printf("failed to list pods: %s", err.Error())
			return completionItems
		}
		return pods
	}

	// Looking for *.network files and `podman network ls`
	if property == "Network" {
		networks, err := listNetworks(c)
		if err != nil {
			fmt.Printf("failed to list networks: %s", err.Error())
			return completionItems
		}
		return networks
	}

	// Check what ports are exposed in the image
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

	// Check what is the specified user in the image
	if strings.HasPrefix(lines[lineNumber], "UserNS=keep-id:") {
		id, err := listUserIdFromImage(
			c,
			lines,
			lineNumber,
		)
		if err != nil {
			fmt.Printf("failed to list user from image: %s", err.Error())
			return completionItems
		}
		return id
	}

	section := findSection(lines, lineNumber)

	if section == "" {
		return completionItems
	}

	// Generic static suggestions based on `properties.go` file
	for _, p := range data.PropertiesMap[section] {
		if property == p.Label {
			for _, parm := range p.Parameters {
				completionItems = append(completionItems, protocol.CompletionItem{
					Label: parm,
				})
			}
			return completionItems
		}
	}

	return completionItems
}

// If somebody type `UserNS=keep-id:`, then check if image has any user
// defined, and provide its id for uid and gid as well
func listUserIdFromImage(c utils.Commander, lines []string, lineNumber protocol.UInteger) ([]protocol.CompletionItem, error) {
	var completionItems []protocol.CompletionItem

	imageName := findImageName(lines, lineNumber)

	// We've found something, looking for User
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

		user, ok := config["User"].(string)
		if !ok {
			return nil, err
		}

		completionItems = append(completionItems, protocol.CompletionItem{
			Label: "uid=" + user,
		})
		completionItems = append(completionItems, protocol.CompletionItem{
			Label: "gid=" + user,
		})
	}

	return completionItems, nil
}

// If user at the `PublihsPort=` line, and typting the exposed port number
// provide suggestions based on image inspect what ports can be exposed.
func listPublishedPorts(c utils.Commander, lines []string, lineNumber protocol.UInteger) ([]protocol.CompletionItem, error) {
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
	imageName := findImageName(lines, lineNumber)

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

// List *.network files and looking for out put `podman network ls`.
func listNetworks(c utils.Commander) ([]protocol.CompletionItem, error) {
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
	volFiles, err := utils.ListQuadletFiles("*.network")
	if err != nil {
		return nil, err
	}
	completionItems = append(completionItems, volFiles...)
	return completionItems, nil
}

// List *.volume files and looking for `podman volume ls`.
func listVolumes(c utils.Commander, line string) ([]protocol.CompletionItem, error) {
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
	volFiles, err := utils.ListQuadletFiles("*.volume")
	if err != nil {
		return nil, err
	}
	completionItems = append(completionItems, volFiles...)

	return completionItems, nil
}

// Looking for `podman secret ls`
func listSecrets(c utils.Commander, line string) ([]protocol.CompletionItem, error) {
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

// List *.images files and `podman images`
func listImages(c utils.Commander) ([]protocol.CompletionItem, error) {
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

	files, err := utils.ListQuadletFiles("*.image")
	if err != nil {
		return completionItems, err
	}
	completionItems = append(completionItems, files...)

	return completionItems, nil
}
