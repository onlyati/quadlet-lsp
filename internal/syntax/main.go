package syntax

import (
	"sync"

	protocol "github.com/tliron/glsp/protocol_3_16"
)

type SyntaxChecker struct {
	documentText string
	uri          string
	hintDiag     protocol.DiagnosticSeverity
	infoDiag     protocol.DiagnosticSeverity
	warnDiag     protocol.DiagnosticSeverity
	errDiag      protocol.DiagnosticSeverity
	checks       []func(SyntaxChecker) []protocol.Diagnostic
}

func NewSyntaxChecker(documentText, uri string) SyntaxChecker {
	return SyntaxChecker{
		documentText: documentText,
		uri:          uri,
		hintDiag:     protocol.DiagnosticSeverityHint,
		infoDiag:     protocol.DiagnosticSeverityInformation,
		warnDiag:     protocol.DiagnosticSeverityWarning,
		errDiag:      protocol.DiagnosticSeverityError,
		checks: []func(SyntaxChecker) []protocol.Diagnostic{
			qsr001,
			qsr002,
			qsr003,
			qsr004,
			qsr005,
			qsr006,
			qsr007,
		},
	}
}

func (s SyntaxChecker) RunAll() []protocol.Diagnostic {
	var wg sync.WaitGroup
	diagChan := make(chan []protocol.Diagnostic, len(s.checks))

	for _, fn := range s.checks {
		wg.Add(1)
		go func(rule func(SyntaxChecker) []protocol.Diagnostic) {
			defer wg.Done()
			result := fn(s)
			if result != nil {
				diagChan <- result
			}
		}(fn)
	}

	wg.Wait()
	close(diagChan)

	var allDiags []protocol.Diagnostic
	for diags := range diagChan {
		allDiags = append(allDiags, diags...)
	}

	return allDiags
}
