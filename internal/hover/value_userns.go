package hover

import (
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func handleValueUserNS(info HoverInformation) *protocol.Hover {
	parts := strings.SplitN(info.value, ":", 2)
	key := parts[0]

	valueStartPos := protocol.UInteger(strings.Index(info.Line, "=")) + 1
	valueCharacterPos := info.CharacterPosition - valueStartPos

	if valueCharacterPos <= uint32(len(key)) {
		switch key {
		case "auto":
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: "**Container user: nil (Host User UID is not mapped into container.)**\n\nPodman allocates unique ranges of UIDs and GIDs from the containers subordinate user IDs. The size of the ranges is based on the number of UIDs required in the image. The number of UIDs and GIDs can be overridden with the size option.",
				},
				Range: &protocol.Range{
					Start: protocol.Position{Line: info.LineNumber, Character: valueStartPos},
					End:   protocol.Position{Line: info.LineNumber, Character: valueStartPos + 4},
				},
			}
		case "host":
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind: protocol.MarkupKindMarkdown, Value: "**Container user: 0 (Default User account mapped to root user in container.)**\n\nhost or “” (empty string): run in the user namespace of the caller. The processes running in the container have the same privileges on the host as any other process launched by the calling user.",
				},
				Range: &protocol.Range{
					Start: protocol.Position{Line: info.LineNumber, Character: valueStartPos},
					End:   protocol.Position{Line: info.LineNumber, Character: valueStartPos + 4},
				},
			}
		case "keep-id":
			desc := []string{
				"keep-id: creates a user namespace where the current user’s UID:GID are mapped to the same values in the container. For containers created by root, the current mapping is created into a new user namespace.",
				"",
				"Valid keep-id options:",
				"- uid=UID: override the UID inside the container that is used to map the current user to.",
				"- gid=GID: override the GID inside the container that is used to map the current user to.",
				"- size=SIZE: override the size of the configured user namespace. It is useful to not saturate all the available IDs. Not supported when running as root.",
			}
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: "**Container user: $UID (Map user account to same UID within container.)**\n\n" + strings.Join(desc, "\n"),
				},
				Range: &protocol.Range{
					Start: protocol.Position{Line: info.LineNumber, Character: valueStartPos},
					End:   protocol.Position{Line: info.LineNumber, Character: valueStartPos + 7},
				},
			}
		case "nomap":
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: "**Container user: nil (Host User UID is not mapped into container.)**\n\nnomap: creates a user namespace where the current rootless user’s UID:GID are not mapped into the container. This option is not allowed for containers created by the root user.",
				},
				Range: &protocol.Range{
					Start: protocol.Position{Line: info.LineNumber, Character: valueStartPos},
					End:   protocol.Position{Line: info.LineNumber, Character: valueStartPos + 5},
				},
			}
		}
	}

	return nil
}
