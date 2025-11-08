// Package syntax
//
// This package stores everything that is related to syntax checkint.
package syntax

import (
	"slices"
	"strings"
	"sync"

	"github.com/onlyati/quadlet-lsp/internal/utils"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var (
	infoDiag = protocol.DiagnosticSeverityInformation
	warnDiag = protocol.DiagnosticSeverityWarning
	errDiag  = protocol.DiagnosticSeverityError
)

type SyntaxChecker struct {
	documentText string
	uri          string
	checks       []rule
	commander    utils.Commander
	config       *utils.QuadletConfig
}

type rule struct {
	name string
	fn   func(SyntaxChecker) []protocol.Diagnostic
}

func NewSyntaxChecker(documentText, uri string) SyntaxChecker {
	return SyntaxChecker{
		documentText: documentText,
		uri:          uri,
		checks: []rule{
			{"qsr001", qsr001},
			{"qsr002", qsr002},
			{"qsr003", qsr003},
			{"qsr004", qsr004},
			{"qsr005", qsr005},
			{"qsr006", qsr006},
			{"qsr007", qsr007},
			{"qsr008", qsr008},
			{"qsr009", qsr009},
			{"qsr010", qsr010},
			{"qsr011", qsr011},
			{"qsr012", qsr012},
			{"qsr013", qsr013},
			{"qsr014", qsr014},
			{"qsr015", qsr015},
			{"qsr016", qsr016},
			{"qsr017", qsr017},
			{"qsr018", qsr018},
			{"qsr019", qsr019},
			{"qsr020", qsr020},
			{"qsr021", qsr021},
			{"qsr022", qsr022},
			{"qsr023", qsr023},
			{"qsr024", qsr024},
			{"qsr025", qsr025},
			{"qsr026", qsr026},
		},
		commander: utils.CommandExecutor{},
	}
}

func (s SyntaxChecker) RunAll(config *utils.QuadletConfig) []protocol.Diagnostic {
	s.config = config
	var wg sync.WaitGroup
	diagChan := make(chan []protocol.Diagnostic, len(s.checks))

	for _, check := range s.checks {
		// Check if rule is disabled
		if slices.Contains(config.Disable, check.name) {
			continue
		}

		// Check if rule is disabled on file level
		disabledByFile := false
		for line := range strings.SplitSeq(s.documentText, "\n") {
			if !strings.HasPrefix(line, "# disable-qsr:") && !strings.HasPrefix(line, "; disable-qsr:") {
				break
			}

			if slices.Contains(strings.Split(line, " "), check.name) {
				disabledByFile = true
				break
			}
		}

		if disabledByFile {
			continue
		}

		// Perform the syntac check
		wg.Add(1)
		go func(rule func(SyntaxChecker) []protocol.Diagnostic) {
			defer wg.Done()
			result := check.fn(s)
			if result != nil {
				diagChan <- result
			}
		}(check.fn)
	}

	wg.Wait()
	close(diagChan)

	var allDiags []protocol.Diagnostic
	for diags := range diagChan {
		allDiags = append(allDiags, diags...)
	}

	return allDiags
}
