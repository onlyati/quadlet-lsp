package hover

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func handleValueUserNS(info HoverInformation) *protocol.Hover {
	value, ok := info.TokenInfo.CurrentNode.(*parser.ValueNode)
	if !ok {
		return nil
	}

	hoverData := []string{
		"**UserNS description**",
		"",
		"`auto`: Container user: nil (Host User UID is not mapped into container.)",
		"",
		"Podman allocates unique ranges of UIDs and GIDs from the containers subordinate user IDs. The size of the ranges is based on the number of UIDs required in the image. The number of UIDs and GIDs can be overridden with the size option.",
		"",
		"`host`: Container user: 0 (Default User account mapped to root user in container.)",
		"",
		"host or “” (empty string): run in the user namespace of the caller. The processes running in the container have the same privileges on the host as any other process launched by the calling user.",
		"",
		"`keep-id`: creates a user namespace where the current user’s UID:GID are mapped to the same values in the container.",
		"",
		"For containers created by root, the current mapping is created into a new user namespace.",
		"Valid keep-id options:",
		"- uid=UID: override the UID inside the container that is used to map the current user to.",
		"- gid=GID: override the GID inside the container that is used to map the current user to.",
		"- size=SIZE: override the size of the configured user namespace. It is useful to not saturate all the available IDs. Not supported when running as root.",
		"",
		"`nomap`: Container user: nil (Host User UID is not mapped into container.)",
		"",
		"nomap: creates a user namespace where the current rootless user’s UID:GID are not mapped into the container. This option is not allowed for containers created by the root user.",
	}

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: strings.Join(hoverData, "\n"),
		},
		Range: &protocol.Range{
			Start: protocol.Position{
				Line:      value.StartPos.LineNumber,
				Character: value.StartPos.Position,
			},
			End: protocol.Position{
				Line:      value.EndPos.LineNumber,
				Character: value.EndPos.Position,
			},
		},
	}
}
