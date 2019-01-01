package commands

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/cli/runner"
	"github.com/raedahgroup/godcr/cli/termio"
)

type HelpCommand struct {
	runner.ParserCommand
	Args struct {
		CommandName string `positional-arg-name:"command-name"`
	} `positional-args:"yes"`
}

func (h HelpCommand) Run(parser *flags.Parser, args []string) error {
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

	PrintCommandHelp(parser.Name, targetCommand)

	return nil
}

func PrintCommandHelp(appName string, command *flags.Command) {
	tabWriter := termio.StdoutWriter
	fmt.Fprintln(tabWriter ,fmt.Sprintf("%s. %s\n", command.ShortDescription, command.LongDescription))

	usageText := fmt.Sprintf("Usage: %s %s", appName, command.Name)
	args := command.Args()
	if args != nil && len(args) > 0 {
		usageText += " [args]"
	}
	usageText += " [options]"
	fmt.Fprintln(tabWriter, usageText)
	fmt.Fprintln(tabWriter)

	if args != nil && len(args) > 0 {
		fmt.Fprintln(tabWriter,"Arguments:")
		for _, arg := range args {
			required := ""
			if arg.Required == 1 {
				required = "(required)"
			}
			fmt.Fprintln(tabWriter,fmt.Sprintf("%s %s \t %s", arg.Name, required, arg.Description))
		}
		fmt.Fprintln(tabWriter)
	}

	options := command.Options()
	if options != nil && len(options) > 0 {
		fmt.Fprintln(tabWriter,"Options:")
		// option printout attempts to add 2 whitespace for options with short name and 4 for those without in order to
		// This is an attempt to stay consistent with the output of parser.WriteHelp
		for _, option := range options {
			optionUsage := " " + " "
			if option.ShortName != 0 {
				optionUsage += "-" + string(option.ShortName) + ", "
			}else {
				optionUsage += " " + " " + " " + " "
			}
			fmt.Fprintln(tabWriter,fmt.Sprintf("%s--%s \t %s", optionUsage, option.LongName, option.Description))
		}
		fmt.Fprintln(tabWriter)
	}

	fmt.Fprintln(tabWriter, fmt.Sprintf("Use %s -h to view application options", appName))
	tabWriter.Flush()
}
