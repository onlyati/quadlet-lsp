package lsp

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/onlyati/quadlet-lsp/internal/commands"
	"github.com/onlyati/quadlet-lsp/internal/syntax"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type CliMessenger struct {
	messages []struct {
		level utils.MessengerLevel
		text  string
	}
}

func (c *CliMessenger) SendMessage(level utils.MessengerLevel, text string) {
	c.messages = append(c.messages, struct {
		level utils.MessengerLevel
		text  string
	}{level, text})
}

func runDocCLI(args []string, commander utils.Commander) (int, []string) {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("failed to read current working directory: %s", err.Error())
		return 1, nil
	}

	parm := ""
	if len(args) >= 3 {
		parm = args[2]
	}
	e := commands.NewEditorCommandExecutor(cwd)
	messenger := CliMessenger{}
	commands.GenerateDoc(
		&protocol.ExecuteCommandParams{
			Command:   "generateDoc",
			Arguments: []any{parm},
		},
		&e,
		&messenger,
		commander,
	)

	rc := 0
	for _, m := range messenger.messages {
		if m.level == utils.MessengerError || m.level == utils.MessengerWarning {
			rc = 4
		}
		fmt.Println(m.text)
	}

	return rc, nil
}

func runCheckCLI(args []string, commander utils.Commander) (int, []string) {
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
