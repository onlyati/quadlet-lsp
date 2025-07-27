package lsp

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/syntax"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "quadlet"

var (
	version   = "0.2.0"
	handler   protocol.Handler
	config    utils.QuadletConfig
	documents = newDocuments()
)

// Entry point of the language server
func Start() {
	// To download automatically during extension install,
	// I've made a simple `version` subcommand, so it is easy
	// to verify which version is downloaded and need to download newer one
	args := os.Args
	if len(args) == 2 {
		if args[1] == "version" {
			fmt.Println(version)
			return
		}
	}

	handler = protocol.Handler{
		// The `hello` and `goodbye` handlers
		Initialize:  initialize,
		Initialized: initialized,
		Shutdown:    shutdown,

		// Make LSP thing, point of this whole thing
		TextDocumentCompletion: textCompletion,
		TextDocumentHover:      textHover,
		TextDocumentDefinition: textDefinition,
		TextDocumentReferences: textReferences,

		// Store document in memory and keep every changes on track
		// It is needed when we want to looking for something in a file,
		// like looking for references.
		TextDocumentDidOpen: func(ctx *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
			uri := string(params.TextDocument.URI)
			documents.add(uri, params.TextDocument.Text)

			// Check syntax when file is open
			checker := syntax.NewSyntaxChecker(documents.read(uri), uri)

			diags := checker.RunAll()
			if len(diags) > 0 {
				ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
					URI:         protocol.DocumentUri(uri),
					Diagnostics: diags,
				})
			} else {
				ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
					URI:         protocol.DocumentUri(uri),
					Diagnostics: []protocol.Diagnostic{},
				})
			}

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

			// Check syntax when file is changed
			checker := syntax.NewSyntaxChecker(documents.read(uri), uri)

			diags := checker.RunAll()
			if len(diags) > 0 {
				ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
					URI:         protocol.DocumentUri(uri),
					Diagnostics: diags,
				})
			} else {
				ctx.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
					URI:         protocol.DocumentUri(uri),
					Diagnostics: []protocol.Diagnostic{},
				})
			}

			return nil
		},
		TextDocumentDidClose: func(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
			uri := string(params.TextDocument.URI)
			documents.delete(uri)
			return nil
		},

		// Whenever a save happen, perform a syntax checking
		TextDocumentDidSave: SyntaxCheckOnSave,
	}

	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	// Read and parse configuration
	workspaceDir := *params.RootURI

	if len(params.WorkspaceFolders) > 0 {
		workspaceDir = params.WorkspaceFolders[0].URI
	}

	workspaceDir, _ = strings.CutPrefix(workspaceDir, "file://")

	cfg, err := utils.LoadConfig(workspaceDir)
	if err != nil {
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Type:    protocol.MessageTypeLog,
			Message: "Failed to read .quadletrc.json file, goes with defaults",
		})
	}
	config = cfg

	startFileWatcher(context, path.Join(workspaceDir, ".quadletrc.json"))

	// Setup server
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
	// Detect Podman version if necesarry
	defPodman := utils.PodmanVersion{}
	if config.Podman != defPodman {
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Type:    protocol.MessageTypeLog,
			Message: fmt.Sprintf("Podman version is overriden from config: %v", config.Podman),
		})
		return nil
	}

	c := utils.CommandExecutor{}
	pVersion, err := utils.NewPodmanVersion(c)
	if err != nil {
		config.Podman = utils.PodmanVersion{Version: 99, Release: 99, Minor: 99}
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Type:    protocol.MessageTypeWarning,
			Message: "Failed to fetch Podman version, assumes it is the latest",
		})
	} else {
		config.Podman = pVersion
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Type:    protocol.MessageTypeLog,
			Message: fmt.Sprintf("Detected Podman version: %v", config.Podman),
		})
	}

	return nil
}

func shutdown(context *glsp.Context) error {
	return nil
}
