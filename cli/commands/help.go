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
	fmt.Printf("%s. %s\n", command.ShortDescription, command.LongDescription)
	fmt.Println()

	args := command.Args()
	options := command.Options()
	usageText := fmt.Sprintf("Usage: %s %s", appName, command.Name)
	if args != nil && len(args) > 0 {
		usageText += " [args]"
	}
	usageText += " [options]"
	fmt.Println(usageText)
	fmt.Println()

	if args != nil && len(args) > 0 {
		fmt.Println("Arguments:")
		longestNameLength := 0
		for _, arg := range args {
			argNameLength := len(arg.Name)
			// required args takes extra 10 space for printing the text '(required)'
			if arg.Required == 1 {
				argNameLength += 10
			}
			if argNameLength > longestNameLength {
				longestNameLength = argNameLength
			}
		}

		writer := termio.StdoutWriter
		for _, arg := range args {
			required := ""
			if arg.Required == 1 {
				required = "(required)"
			}
			fmt.Fprintln(writer,fmt.Sprintf("%s %s \t %s", arg.Name, required, arg.Description))
		}
		writer.Flush()
		fmt.Println()
	}

	if options != nil && len(options) > 0 {
		fmt.Println("Options:")
		longestNameLength := 0
		for _, option := range options {
			if len(option.LongName) > longestNameLength {
				longestNameLength = len(option.LongName)
			}
		}
		writer := termio.StdoutWriter
		for _, option := range options {
			name := "  "
			if option.ShortName != 0 {
				name += "-" + string(option.ShortName) + ", "
			}else {
				name += " " + " " + " " + " "
			}
			fmt.Fprintln(writer,fmt.Sprintf("%s--%s \t %s", name, option.LongName, option.Description))
		}
		writer.Flush()
		fmt.Println()
	}

	fmt.Println(fmt.Sprintf("Use %s -h to view application options", appName))

}
