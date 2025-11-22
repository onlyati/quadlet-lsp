// Package cli
//
// Command line interface related actions
package cli

import (
	"errors"
	"fmt"

	"github.com/onlyati/quadlet-lsp/internal/data"
	"github.com/onlyati/quadlet-lsp/internal/utils"
)

type CliCommand struct {
	Command string
	Parms   []string
}

// Execute Entrypoint of the command line interface, this is the command router.
func (c CliCommand) Execute() error {
	output := []string{}
	var err error
	switch c.Command {
	case "help":
		help()
	case "check":
		output, err = runCheckCLI(c.Parms, utils.CommandExecutor{})
	case "doc":
		output, err = runDocCLI(c.Parms, utils.CommandExecutor{})
	case "version":
		output = []string{data.ProgramVersion}
	default:
		err = errors.New("invalid command, see 'quadlet-lsp help'")
	}

	for _, l := range output {
		fmt.Println(l)
	}
	if err != nil {
		return err
	}

	return nil
}

func help() {
	fmt.Println("Usage of quadlet-lsp")
	fmt.Println("")
	fmt.Println("Display version: quadlet-lsp version")
	fmt.Println("")
	fmt.Println("Run syntax checks: quadlet-lsp check <dir>")
	fmt.Println("    <dir>: Directory which should be scanned")
	fmt.Println("")
	fmt.Println("Generate documents: quadlet-lsp doc <type>")
	fmt.Println("    <type>: can be: html, md or json")
}
