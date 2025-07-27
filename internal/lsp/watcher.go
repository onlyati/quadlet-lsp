package lsp

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func startFileWatcher(ctx *glsp.Context, path string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("failed to create watcher:", err)
		return
	}

	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Only react to Write or Rename
				if event.Op&(fsnotify.Write|fsnotify.Rename) != 0 {
					log.Printf("File changed: %s", event.Name)

					// Send info-level window/showMessage
					msg := &protocol.ShowMessageParams{
						Type:    protocol.MessageTypeInfo,
						Message: "File changed: " + event.Name,
					}
					ctx.Notify(protocol.ServerWindowShowMessage, msg)
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
