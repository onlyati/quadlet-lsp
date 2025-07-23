package lsp

import (
	"fmt"
	"os"

	"github.com/tliron/commonlog"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "quadlet"

var (
	version   = "0.2.0"
	handler   protocol.Handler
	documents = newDocuments()
)

func Start() {
	args := os.Args
	if len(args) == 2 {
		if args[1] == "version" {
			fmt.Println(version)
			return
		}
	}

	commonlog.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,

		// Make LSP thing
		TextDocumentCompletion: textCompletion,
		TextDocumentHover:      textHover,
		TextDocumentDefinition: textDefinition,
		TextDocumentReferences: textReferences,

		// Document sync
		TextDocumentDidOpen: func(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
			uri := string(params.TextDocument.URI)
			documents.add(uri, params.TextDocument.Text)
			return nil
		},
		TextDocumentDidChange: func(ctx *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
			uri := string(params.TextDocument.URI)
			if text, ok := documents.checkUri(uri); ok {
				for _, change := range params.ContentChanges {
					if change_, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
						startIndex, endIndex := change_.Range.IndexesIn(text)
						text = text[:startIndex] + change_.Text + text[endIndex:]
					} else if change_, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
						text = change_.Text
					}
				}
				documents.add(uri, text)
			}
			return nil
		},
		TextDocumentDidClose: func(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
			uri := string(params.TextDocument.URI)
			documents.delete(uri)
			return nil
		},
		TextDocumentDidSave: func(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
			return nil
		},
	}

	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	commonlog.NewInfoMessage(0, "Initializing server...")

	capabilities := handler.CreateServerCapabilities()

	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{"=", "[", "]", ".", ":", ","},
	}
	capabilities.HoverProvider = &protocol.HoverOptions{}

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	return nil
}
