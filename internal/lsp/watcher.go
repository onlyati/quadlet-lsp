package lsp

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func startFileWatcher(ctx *glsp.Context, path string, cfg *utils.QuadletConfig) {
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
						time.Sleep(100 * time.Millisecond)

						err := watcher.Remove(path)
						if err != nil {
							log.Println("error removing watcher:", err)
						}
						err = watcher.Add(path)
						if err != nil {
							log.Println("error re-adding watcher:", err)
						}

						tmpCfg, err := utils.LoadConfig(path)
						cfg.Mu.Lock()
						cfg.Disable = tmpCfg.Disable
						cfg.PodmanVersion = tmpCfg.PodmanVersion
						cfg.Podman = tmpCfg.Podman

						ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
							Type:    protocol.MessageTypeInfo,
							Message: fmt.Sprintf("config has been reloaded, Podman version target: %v", config.Podman),
						})

						if !cfg.Podman.IsSupported() {
							ctx.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
								Type:    protocol.MessageTypeWarning,
								Message: "The specified or found Podman version is not fully supported (>= 5.4.0)",
							})
						}

						cfg.Mu.Unlock()
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
