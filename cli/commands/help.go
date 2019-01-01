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
	if len(command.LongDescription) > 0 {
		fmt.Println(command.LongDescription)
	}else{
		fmt.Println(command.ShortDescription)
	}
	fmt.Println()

	args := command.Args()
	options := command.Options()
	usageText := fmt.Sprintf("Usage: %s %s", appName, command.Name)
	if args != nil && len(args) > 0 {
		usageText += " [args]"
	}
	if options != nil && len(options) > 0 {
		usageText += " [options]"
	}
	fmt.Println(usageText)
	fmt.Println()

	generateWhiteSpace := func(inputLength, inputLengthCap int) string {
		makeSpace := func(count int) (output string) {
			for i := 0; i < count; i++ {
				output += " "
			}
			return
		}
		maxSpaces := 2
		if inputLength == inputLengthCap {
			return makeSpace(maxSpaces)
		}

		spacesToMake := maxSpaces - (inputLength - inputLengthCap)
		if spacesToMake < 1 {
			spacesToMake = 1
		}
		return makeSpace(spacesToMake)
	}

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

		for _, arg := range args {
			required := ""
			if arg.Required == 1 {
				required = "(required)"
			}
			fmt.Println(fmt.Sprintf("%s %s%s %s", arg.Name, required, generateWhiteSpace(len(arg.Name), longestNameLength), arg.Description))
		}
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
		for _, option := range options {
			output := "  "
			if option.ShortName != 0 {
				output += "-" + string(option.ShortName) + ", "
			}else {
				output += " " + " " + " " + " "
			}
			// the length of the name determines the number of extrapaces
			output += fmt.Sprintf("--%s %s %s", option.LongName, generateWhiteSpace(len(option.LongName), longestNameLength), option.Description)
			fmt.Println(output)
		}
		fmt.Println()
	}

	fmt.Println(fmt.Sprintf("Use %s -h to view application options", appName))

}
