package commands

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/cli/termio"
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

	PrintCommandHelp(parser.Name, targetCommand)

	return nil
}

func PrintCommandHelp(appName string, command *flags.Command) {
	tabWriter := termio.StdoutWriter
	fmt.Fprintln(tabWriter ,fmt.Sprintf("%s. %s\n", command.ShortDescription, command.LongDescription))

	usageText := fmt.Sprintf("Usage:\n  %s %s", appName, command.Name)
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
			fmt.Fprintln(tabWriter,fmt.Sprintf("  %s %s \t %s", arg.Name, required, arg.Description))
		}
		fmt.Fprintln(tabWriter)
	}

	options := command.Options()
	if options != nil && len(options) > 0 {
		fmt.Fprintln(tabWriter,"Options:")
		// option printout attempts to add 2 whitespace for options with short name and 6 for those without
		// This is an attempt to stay consistent with the output of parser.WriteHelp
		for _, option := range options {
			var optionUsage string

			if option.ShortName != 0 && option.LongName != "" {
				optionUsage = fmt.Sprintf("  -%c, --%s", option.ShortName, option.LongName)
			} else if option.ShortName != 0 {
				optionUsage = fmt.Sprintf("  -%c", option.ShortName)
			} else {
				optionUsage = fmt.Sprintf("      --%s", option.LongName)
			}

			fmt.Fprintln(tabWriter,fmt.Sprintf("%s \t %s", optionUsage, option.Description))
		}
		fmt.Fprintln(tabWriter)
	}

	fmt.Fprintln(tabWriter, fmt.Sprintf("Use %s -h to view application options", appName))
	tabWriter.Flush()
}
