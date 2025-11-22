package main

import (
	"fmt"
	"os"

	"github.com/onlyati/quadlet-lsp/internal/cli"
	"github.com/onlyati/quadlet-lsp/internal/lsp"
)

func main() {
	args := os.Args

	if len(args) >= 2 {
		// This is CLI command
		cmd := args[1]
		parms := []string{}
		if len(args) > 2 {
			parms = args[2:]
		}
		program := cli.CliCommand{
			Command: cmd,
			Parms:   parms,
		}
		err := program.Execute()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(8)
		}
		os.Exit(0)
	}

	// Start server
	lsp.Start()
}
