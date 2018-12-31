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

	PrintCommandHelp(parser.Name, targetCommand)

	return nil
}

func PrintCommandHelp(appName string, command *flags.Command) {
	helpParser := flags.NewParser(nil, flags.HelpFlag)
	helpParser.Name = appName
	helpParser.Active = command
	helpParser.WriteHelp(os.Stdout)
	fmt.Printf("To view application options, use '%s help'\n", appName)
}

/*
I have looked for how to get command-specific options but did
find a way around it and so went back to this. I will be happy to have a pointer to this
*/