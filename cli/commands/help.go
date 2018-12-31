package commands

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/cli/runner"
	"github.com/raedahgroup/godcr/cli/termio"
	"os"
)

type HelpCommand struct {
	runner.ParserCommand
	Args struct{
		CommandName string `positional-arg-name:"command-name"`
	} `positional-args:"yes"`
}

func (h HelpCommand) Run(ctx context.Context, parser *flags.Parser, args []string) error {
	if h.Args.CommandName == "" {
		active := parser.Active
		parser.Active = nil
		defer func() {parser.Active = active}()
		parser.WriteHelp(termio.StdoutWriter)
		return nil
	}

	var targetCommand = parser.Find(h.Args.CommandName)
	if targetCommand == nil {
		return fmt.Errorf("unknown command %q", h.Args.CommandName)
	}

	helpParser := flags.NewParser(nil, flags.HelpFlag)
	helpParser.Name = parser.Name
	helpParser.Active = targetCommand
	helpParser.WriteHelp(os.Stderr)
	fmt.Printf("To view application options, use '%s help'\n", parser.Name)

	return nil
}