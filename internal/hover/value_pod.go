package hover

import (
	"os"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// handleValuePod is looking for the the pod file and display its content on
// hover output.
func handleValuePod(info HoverInformation) *protocol.Hover {
	value, ok := info.TokenInfo.CurrentNode.(*parser.ValueNode)
	if !ok {
		return nil
	}
	pod := *value.Value
	if strings.Contains(pod, "@") {
		pod = utils.ConvertTemplateNameToFile(pod)
	}

	files, err := utils.ListQuadletFiles("pod", info.RootDir, info.Level)
	if err != nil {
		return nil
	}

	for _, f := range files {
		if pod == f.Label {
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
