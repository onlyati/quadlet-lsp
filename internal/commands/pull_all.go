package commands

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func pullAll(command string, e *EditorCommandExecutor, ctx glsp.Context, executor utils.Commander) {
	defer e.resetRunning(command)

	e.mutex.Lock()
	rootDir := e.rootDir
	e.mutex.Unlock()

	dir, err := os.ReadDir(rootDir)
	if err != nil {
		log.Println("failed to list directory: " + err.Error())
		ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Type:    protocol.MessageTypeError,
			Message: "failed to list directory: " + err.Error(),
		})
		return
	}

	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		isItContainer := strings.HasSuffix(entry.Name(), ".container")
		isItImage := strings.HasSuffix(entry.Name(), ".image")
		isItVolume := strings.HasSuffix(entry.Name(), ".volume")
		if !isItContainer && !isItImage && !isItVolume {
			continue
		}

		tmp := strings.Split(entry.Name(), ".")
		section := "[" + utils.FirstCharacterToUpper(tmp[len(tmp)-1]) + "]"

		file, err := os.ReadFile(path.Join(rootDir, entry.Name()))
		if err != nil {
			log.Println("failed to read file: " + err.Error())
			ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
				Type:    protocol.MessageTypeError,
				Message: "failed to read file: " + err.Error(),
			})
			continue
		}

		inSection := false
		for line := range strings.SplitSeq(string(file), "\n") {
			line := strings.TrimSpace(line)

			if line == section && !inSection {
				inSection = true
				continue
			}

			if strings.HasPrefix(line, "[") {
				inSection = false
				continue
			}

			if inSection {
				tmp := strings.SplitN(line, "=", 2)
				if len(tmp) != 2 {
					continue
				}

				if tmp[0] != "Image" {
					continue
				}
				image := tmp[1]

				_, err := executor.Run("podman", "image", "exists", image)
				if err == nil {
					continue
				}

				ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
					Type:    protocol.MessageTypeInfo,
					Message: "Start pulling image: " + image,
				})

				output, err := executor.Run("podman", "image", "pull", image)
				if err != nil {
					ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
						Type:    protocol.MessageTypeError,
						Message: "Failed to pull image: " + image,
					})
				} else {
					ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
						Type:    protocol.MessageTypeInfo,
						Message: "Image pulled: " + image + "\n" + strings.Join(output, "\n"),
					})
				}
			}
		}
	}
}
