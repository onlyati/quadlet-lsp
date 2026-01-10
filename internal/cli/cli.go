// Package cli
//
// Command line interface related actions
package cli

import (
	"errors"
	"fmt"
	"os"

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
	case "version":
		output = []string{data.ProgramVersion}
	case "check":
		output, err = runCheckCLI(os.Args, utils.CommandExecutor{})
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
}
