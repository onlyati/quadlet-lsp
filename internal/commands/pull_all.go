package commands

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
)

func pullAll(command string, e *EditorCommandExecutor, messenger utils.Messenger, executor utils.Commander) {
	defer e.resetRunning(command)

	e.mutex.Lock()
	rootDir := e.rootDir
	e.mutex.Unlock()

	dir, err := os.ReadDir(rootDir)
	if err != nil {
		log.Println("failed to list directory: " + err.Error())
		messenger.SendMessage(
			utils.MessengerError,
			"failed to list directory: "+err.Error(),
		)
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
			messenger.SendMessage(
				utils.MessengerError,
				"failed to read file: "+err.Error(),
			)
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

				messenger.SendMessage(
					utils.MessengerInfo,
					"Start pulling image: "+image,
				)

				output, err := executor.Run("podman", "image", "pull", image)
				if err != nil {
					messenger.SendMessage(
						utils.MessengerError,
						"Failed to pull image: "+image,
					)
				} else {
					messenger.SendMessage(
						utils.MessengerInfo,
						"Image pulled: "+image+"\n"+strings.Join(output, "\n"),
					)
				}
			}
		}
	}
}
