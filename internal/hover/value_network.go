package hover

import (
	"os"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// handleValuePod is looking for the the pod file and display its content on
// hover output.
func handleValueNetwork(info HoverInformation) *protocol.Hover {
	network := info.value
	if strings.Contains(network, "@") {
		network = utils.ConvertTemplateNameToFile(network)
	}

	files, err := utils.ListQuadletFiles("network", info.RootDir, info.Level)
	if err != nil {
		return nil
	}

	for _, f := range files {
		if network == f.Label {
			switch v := f.Documentation.(type) {
			case string:
				msg := []string{}
				path, _ := strings.CutPrefix(v, "From work directory: ")
				content, err := os.ReadFile(path)
				if err == nil {
					msg = append(msg, "**Content of file**")
					msg = append(msg, "```quadlet")
					msg = append(msg, strings.Split(string(content), "\n")...)
					msg = append(msg, "```")
				}

				return &protocol.Hover{
					Contents: protocol.MarkupContent{
						Kind:  protocol.MarkupKindMarkdown,
						Value: strings.Join(msg, "\n"),
					},
				}
			}
		}
	}

	return nil
}
