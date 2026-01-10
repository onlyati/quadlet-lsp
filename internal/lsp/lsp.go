// Package lsp
//
// Main entrypoint of the project, this is where the language server is starting.
package lsp

import (
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/commands"
	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/internal/syntax"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "quadlet"

var (
	version   = data.ProgramVersion
	handler   protocol.Handler
	config    *utils.QuadletConfig
	documents = utils.NewDocuments()
	commander commands.EditorCommandExecutor
)

// Start Entry point of the language server
func Start() {
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
			documents.Add(uri, params.TextDocument.Text)

			// Check syntax when file is open
			checker := syntax.NewSyntaxChecker(documents.Read(uri), uri)

			diags := checker.RunAll(config)
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
			if text, ok := documents.CheckURI(uri); ok {
				for _, change := range params.ContentChanges {
					if change_, ok := change.(protocol.TextDocumentContentChangeEvent); ok {
						startIndex, endIndex := change_.Range.IndexesIn(text)
						text = text[:startIndex] + change_.Text + text[endIndex:]
					} else if change_, ok := change.(protocol.TextDocumentContentChangeEventWhole); ok {
						text = change_.Text
					}
				}
				documents.Add(uri, text)
			}

			// Check syntax when file is changed
			checker := syntax.NewSyntaxChecker(documents.Read(uri), uri)

			diags := checker.RunAll(config)
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
			documents.Delete(uri)
			return nil
		},

		// Whenever a save happen, perform a syntax checking
		TextDocumentDidSave: SyntaxCheckOnSave,

		// Handle commands that should be executed
		WorkspaceExecuteCommand: ExecuteCommands,

		// Handle format requests
		TextDocumentFormatting: Format,
	}

	server := server.NewServer(&handler, lsName, false)

	err := server.RunStdio()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	// If no RootURI, it means no directory open, do not initialize lsp
	if params.RootURI == nil {
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Type:    protocol.MessageTypeWarning,
			Message: "Open directory to be able to use quadlet-lsp",
		})
		return nil, errors.New("open directory to use quadlet-lsp")
	}

	// Read and parse configuration
	workspaceDir := *params.RootURI

	if len(params.WorkspaceFolders) > 0 {
		workspaceDir = params.WorkspaceFolders[0].URI
	}

	workspaceDir, _ = strings.CutPrefix(workspaceDir, "file://")

	cfg, err := utils.LoadConfig(workspaceDir, utils.CommandExecutor{})
	if err != nil {
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Type:    protocol.MessageTypeLog,
			Message: "Failed to read .quadletrc.json file, goes with defaults",
		})
	}
	config = cfg
	context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
		Type:    protocol.MessageTypeInfo,
		Message: fmt.Sprintf("Detected Podman target version: %v", config.Podman),
	})
	if !config.Podman.IsSupported() {
		context.Notify(protocol.ServerWindowShowMessage, protocol.ShowMessageParams{
			Type:    protocol.MessageTypeWarning,
			Message: "The specified or found Podman version is not fully supported (>= 5.4.0)",
		})
	}

	commander = commands.NewEditorCommandExecutor(cfg.WorkspaceRoot, *cfg.Project.DirLevel)

	startFileWatcher(
		context,
		path.Join(workspaceDir, ".quadletrc.json"),
		config,
		&documents,
	)

	// Setup server
	capabilities := handler.CreateServerCapabilities()
	capabilities.ExecuteCommandProvider = &protocol.ExecuteCommandOptions{
		Commands: []string{"pullAll", "listJobs"},
	}

	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{"=", "[", "]", ".", ":", ",", "%"},
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
