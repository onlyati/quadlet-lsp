package cli

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"github.com/onlyati/quadlet-lsp/internal/syntax"
	"github.com/onlyati/quadlet-lsp/internal/utils"
	"github.com/onlyati/quadlet-lsp/pkg/parser"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

// runCheckCLI Run syntax checks on specific directory from command line.
// This funtion walk through the specified directory and perform all registers
// syntax check. It respect disabled checks normally.
func runCheckCLI(args []string, commander utils.Commander) ([]string, error) {
	log.SetOutput(io.Discard)

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	checkCfg, err := utils.LoadConfig(cwd, commander)
	if err != nil {
		return nil, err
	}

	workEntity := cwd
	if len(args) == 1 {
		workEntity = args[0]
	}
	if workEntity == "." {
		workEntity = cwd
	}

	stat, err := os.Stat(workEntity)
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, errors.New("parameter must be a directory")
	}

	quadletDir, err := parser.ParseQuadletDir(workEntity)
	if err != nil {
		return nil, err
	}

	diags := map[string][]protocol.Diagnostic{}
	for _, q := range quadletDir.Quadlets {
		s := syntax.NewSyntaxChecker(
			q.SourceFile,
			path.Join(workEntity, q.Name),
		)
		tmpDiags := s.RunAll(checkCfg)
		diags[q.Name] = tmpDiags

		for _, d := range q.Dropins {
			s := syntax.NewSyntaxChecker(
				d.SourceFile,
				path.Join(workEntity, d.Directory, d.FileName),
			)
			tmpDiags := s.RunAll(checkCfg)
			diags[path.Join(d.Directory, d.FileName)] = tmpDiags
		}
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
		return output, nil
	}
	return output, nil
}
