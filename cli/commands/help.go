package commands

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/cli/help"
)

type HelpCommand struct {
	commanderStub
	Args struct {
		CommandName string `positional-arg-name:"command-name"`
	} `positional-args:"yes"`
}

func (h HelpCommand) Run(parser *flags.Parser) error {
	if h.Args.CommandName == "" {
		help.PrintGeneralHelp(os.Stdout, HelpParser(), Categories())
		return nil
	}

	var targetCommand = parser.Find(h.Args.CommandName)
	if targetCommand == nil {
		return fmt.Errorf("unknown command %q", h.Args.CommandName)
	}

	help.PrintCommandHelp(os.Stdout, parser.Name, targetCommand)

	return nil
}

type GeneralHelpData struct {
	config.CommandLineOptions `group:"Options:"`
	AvailableCommands
	ExperimentalCommands
}

func HelpParser() *flags.Parser {
	helpData := GeneralHelpData{}
	return flags.NewParser(&helpData, flags.HelpFlag|flags.PassDoubleDash)
}
