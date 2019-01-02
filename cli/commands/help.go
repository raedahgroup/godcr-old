package commands

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/cli/help"
	"github.com/raedahgroup/godcr/cli/termio"
	"os"
)

type HelpCommand struct {
	commanderStub
	Args struct {
		CommandName string `positional-arg-name:"command-name"`
	} `positional-args:"yes"`
}

func (h HelpCommand) Run(parser *flags.Parser) error {
	if h.Args.CommandName == "" {
		active := parser.Active
		parser.Active = nil
		defer func() { parser.Active = active }()
		parser.WriteHelp(termio.StdoutWriter)
		return nil
	}

	var targetCommand = parser.Find(h.Args.CommandName)
	if targetCommand == nil {
		return fmt.Errorf("unknown command %q", h.Args.CommandName)
	}

	help.PrintCommandHelp(os.Stdout, parser.Name, targetCommand)

	return nil
}
