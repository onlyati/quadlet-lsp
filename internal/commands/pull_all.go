package commands

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func pullAll(command *protocol.ExecuteCommandParams, e *EditorCommandExecutor, messenger utils.Messenger, executor utils.Commander) {
	defer e.resetRunning(command.Command)

	e.mutex.Lock()
	rootDir := e.rootDir
	e.mutex.Unlock()

	filepath.WalkDir(rootDir, func(p string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		isItContainer := strings.HasSuffix(entry.Name(), ".container")
		isItImage := strings.HasSuffix(entry.Name(), ".image")
		isItVolume := strings.HasSuffix(entry.Name(), ".volume")
		section := ""

		allowed := []string{"container", "image", "volume"}
		isItDoprins := false
		if strings.HasSuffix(entry.Name(), ".conf") {
			tmp := strings.Split(p, string(os.PathSeparator))
			if len(tmp) > 2 {
				parentDirectory := tmp[len(tmp)-2]
				for _, item := range allowed {
					if strings.HasSuffix(parentDirectory, item+".d") {
						isItDoprins = true
						section = "[" + utils.FirstCharacterToUpper(item) + "]"
					}
				}
			}
		}

		if !isItContainer && !isItImage && !isItVolume && !isItDoprins {
			return nil
		}

		tmp := strings.Split(entry.Name(), ".")
		if section == "" {
			section = "[" + utils.FirstCharacterToUpper(tmp[len(tmp)-1]) + "]"
		}

		file, err := os.ReadFile(p)
		if err != nil {
			log.Println("failed to read file: " + err.Error())
			messenger.SendMessage(
				utils.MessengerError,
				"failed to read file: "+err.Error(),
			)
			return nil
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
		return nil
	})
}
