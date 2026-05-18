package hover

import (
	"strings"

	"github.com/onlyati/quadlet-lsp/pkg/quadlet/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func handleValueSecret(info HoverInformation) *protocol.Hover {
	value, ok := info.TokenInfo.CurrentNode.(*parser.ValueNode)
	if !ok {
		return nil
	}

	hoverData := []string{
		"**Secret description**",
		"",
		"Syntax: `secret-name,option1=value1,option2=value2...",
		"",
		"The `secret-name` is the name that has been created by `podman secret create` command.",
		"",
		"Options:",
		"",
		"- **type=mount|env**: How the secret is exposed to the container. mount mounts the secret into the container as a file. env exposes the secret as an environment variable. Defaults to mount",
		"- **target=target-name**: Target of secret. For mounted secrets, this is the path to the secret inside the container. If a fully qualified path is provided, the secret is mounted at that location. Otherwise, the secret is mounted to /run/secrets/target for Linux containers or /var/run/secrets/target for FreeBSD containers. If the target is not set, the secret is mounted to /run/secrets/secretname by default. For env secrets, this is the environment variable key. Defaults to secretname.",
		"- **uid=n**: UID of secret. Defaults to 0. Mount secret type only.",
		"- **gid=n**: GID of secret. Defaults to 0. Mount secret type only.",
		"- **mode=0nnn**: Mode of secret. Defaults to 0444. Mount secret type only.",
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
