package lsp

import (
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func startFileWatcher(ctx *glsp.Context, path string, cfg *utils.QuadletConfig, docs *utils.Documents) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("failed to create watcher:", err)
		return
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// React to Write or Rename
				if event.Op&(fsnotify.Write|fsnotify.Rename) != 0 {
					// In case of Rename, the file may have been replaced -> re-add watcher
					if event.Op&fsnotify.Rename != 0 {

						// Wait a moment to let the file be recreated
						time.Sleep(200 * time.Millisecond)

						err := watcher.Remove(path)
						if err != nil {
							log.Println("error removing watcher:", err)
						}
						err = watcher.Add(path)
						if err != nil {
							log.Println("error re-adding watcher:", err)
						}

						tmpCfg, err := utils.LoadConfig(path, utils.CommandExecutor{})
						if err != nil {
							ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
								Type:    protocol.MessageTypeError,
								Message: fmt.Sprintf("Failed to load config: %v", err),
							})
							continue
						}
						cfg.Mu.Lock()
						needSyntaxCheck := slices.Compare(cfg.Disable, tmpCfg.Disable) != 0 || cfg.Podman != tmpCfg.Podman

						cfg.Disable = tmpCfg.Disable
						cfg.PodmanVersion = tmpCfg.PodmanVersion
						cfg.Podman = tmpCfg.Podman

						cfg.Mu.Unlock()

						ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
							Type:    protocol.MessageTypeInfo,
							Message: fmt.Sprintf("config has been reloaded, Podman version target: %v", tmpCfg.Podman),
						})

						if !tmpCfg.Podman.IsSupported() {
							ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
								Type:    protocol.MessageTypeWarning,
								Message: "The specified or found Podman version is not fully supported (>= 5.4.0)",
							})
						}

						if needSyntaxCheck {
							CheckAllOpenFileForSyntax(ctx, docs)
						}

					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("watcher error:", err)
			}
		}
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Println("failed to add ", path, " file to watcher:", err)
	}
}
