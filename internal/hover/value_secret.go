package hover

import (
	"strings"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

func handleValueSecret(info HoverInformation) *protocol.Hover {
	valueOffset := protocol.UInteger(strings.Index(info.Line, "=")) + 1

	for part := range strings.SplitSeq(info.value, ",") {
		if info.CharacterPosition >= valueOffset && info.CharacterPosition < uint32(len(part))+valueOffset {
			msg := getSecretDescription(part)
			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: strings.Join(msg, "\n"),
				},
				Range: &protocol.Range{
					Start: protocol.Position{Line: info.LineNumber, Character: valueOffset},
					End:   protocol.Position{Line: info.LineNumber, Character: valueOffset + uint32(len(part))},
				},
			}
		}
		valueOffset += uint32(len(part) + 1)

	}

	return nil
}

func getSecretDescription(part string) []string {
	if strings.HasPrefix(part, "type=") {
		return []string{
			"**type=mount|env**",
			"",
			"",
			"How the secret is exposed to the container. mount mounts the secret into the container as a file. env exposes the secret as an environment variable. Defaults to mount",
		}
	}

	if strings.HasPrefix(part, "target=") {
		return []string{
			"**target=target**",
			"",
			"",
			"Target of secret. For mounted secrets, this is the path to the secret inside the container. If a fully qualified path is provided, the secret is mounted at that location. Otherwise, the secret is mounted to /run/secrets/target for Linux containers or /var/run/secrets/target for FreeBSD containers. If the target is not set, the secret is mounted to /run/secrets/secretname by default. For env secrets, this is the environment variable key. Defaults to secretname.",
		}
	}

	if strings.HasPrefix(part, "uid=") {
		return []string{
			"**uid=0**",
			"",
			"",
			"UID of secret. Defaults to 0. Mount secret type only.",
		}
	}

	if strings.HasPrefix(part, "gid=") {
		return []string{
			"**gid=0**",
			"",
			"",
			"GID of secret. Defaults to 0. Mount secret type only.",
		}
	}

	if strings.HasPrefix(part, "mode=") {
		return []string{
			"**mode=0**",
			"",
			"",
			"Mode of secret. Defaults to 0444. Mount secret type only.",
		}
	}

	// It is probably the secret value
	return []string{
		"Secret name that is used and previously created by `podman secret create` command.",
	}
}
