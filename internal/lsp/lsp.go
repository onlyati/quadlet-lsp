// Package lsp
//
// Main entrypoint of the project, this is where the language server is starting.
package lsp

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/commands"
	"github.com/onlyati/quadlet-lsp/internal/syntax"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	_ "github.com/tliron/commonlog/simple"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
)

const lsName = "quadlet"

var (
	version   = "0.6.0"
	handler   protocol.Handler
	config    *utils.QuadletConfig
	documents = utils.NewDocuments()
	commander commands.EditorCommandExecutor
)

// Entry point of the language server
func Start() {
	// To download automatically during extension install,
	// I've made a simple `version` subcommand, so it is easy
	// to verify which version is downloaded and need to download newer one
	args := os.Args
	if len(args) >= 2 {
		if args[1] == "version" {
			fmt.Println(version)
			return
		}

		if args[1] == "check" {
			rc, output := runCLI(args, utils.CommandExecutor{})
			for _, l := range output {
				fmt.Println(l)
			}
			os.Exit(rc)
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

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
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

	commander = commands.NewEditorCommandExecutor(cfg.WorkspaceRoot)

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

func runCLI(args []string, commander utils.Commander) (int, []string) {
	log.SetOutput(io.Discard)

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to read current working directory: %s", err.Error())
		os.Exit(1)
	}

	checkCfg, err := utils.LoadConfig(cwd, commander)
	if err != nil {
		fmt.Printf("failed to load config: %s", err.Error())
		os.Exit(1)
	}

	workEntity := cwd
	if len(args) == 3 {
		workEntity = args[2]
	}

	stat, err := os.Stat(workEntity)
	if err != nil {
		fmt.Printf("failed to stat info: %s", err.Error())
		os.Exit(1)
	}
	diags := map[string][]protocol.Diagnostic{}

	if stat.IsDir() {
		filepath.WalkDir(workEntity, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			tmp := strings.Split(d.Name(), ".")
			ext := tmp[len(tmp)-1]

			allowed := []string{
				"container",
				"network",
				"pod",
				"kube",
				"volume",
				"build",
				"conf",
			}
			if !slices.Contains(allowed, ext) {
				return nil
			}

			f, err := os.ReadFile(p)
			if err != nil {
				fmt.Printf("failed to read file: %s\n", err.Error())
				return nil
			}

			uri := p
			isStartWithWorkEntity := strings.HasPrefix(uri, workEntity)
			if !isStartWithWorkEntity {
				uri = workEntity + string(os.PathSeparator) + uri
			}
			s := syntax.NewSyntaxChecker(string(f), uri)
			tmpDiags := s.RunAll(checkCfg)

			key, _ := strings.CutPrefix(p, workEntity+string(os.PathSeparator))
			diags[key] = tmpDiags
			return nil
		})
	} else {
		f, err := os.ReadFile(workEntity)
		if err != nil {
			fmt.Printf("failed to read file: %s", err.Error())
			os.Exit(1)
		}
		s := syntax.NewSyntaxChecker(string(f), workEntity)
		tmpDiags := s.RunAll(checkCfg)
		diags[workEntity] = tmpDiags
	}

	output := []string{}
	line := fmt.Sprintf(
		"%-40s, %-18s, %-13s, %s",
		"File",
		"QSR number",
		"Range",
		"Message",
	)
	output = append(output, line)

	found := false
	for f, fDiags := range diags {
		for _, diag := range fDiags {
			if *diag.Severity != protocol.DiagnosticSeverityInformation {
				found = true
			}
			line := fmt.Sprintf(
				"%-40s, %-18s, %02d.%03d-%02d.%03d, %s",
				f,
				*diag.Source,
				diag.Range.Start.Line, diag.Range.Start.Character,
				diag.Range.End.Line, diag.Range.End.Character,
				diag.Message,
			)
			output = append(output, line)
		}
	}

	if found {
		return 4, output
	}
	return 0, output
}
