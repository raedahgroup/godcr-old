package help

import (
	"fmt"
	"io"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/cli/termio"
)

type CommandCategory struct {
	Name         string
	ShortName    string
	CommandNames []string
}

func commandCategoryName(commandName string, commandCategories []*CommandCategory) string {
	for _, category := range commandCategories {
		for _, command := range category.CommandNames {
			if commandName == command {
				return category.Name
			}
		}
	}
	return "Other commands"
}

func PrintGeneralHelp(output io.Writer, parser *flags.Parser, commandCategories []*CommandCategory) {
	tabWriter := termio.TabWriter(output)

	// print version
	fmt.Fprintf(tabWriter, "%s v%s\n", app.Name(), app.Version())
	fmt.Fprintln(tabWriter)

	// print general app options
	printOptionGroups(tabWriter, parser.Groups())

	// loop through all commands registered on parser and separate into groups
	commandGroups := map[string][]*flags.Command{}
	for _, command := range parser.Commands() {
		commandCategory := commandCategoryName(command.Name, commandCategories)
		commandGroups[commandCategory] = append(commandGroups[commandCategory], command)
	}
	printCommands(tabWriter, commandGroups)
}

func PrintCommandHelp(output io.Writer, appName string, command *flags.Command) {
	tabWriter := termio.TabWriter(output)

	// command description
	fmt.Fprintln(tabWriter, fmt.Sprintf("%s. %s\n", command.ShortDescription, command.LongDescription))

	usageText := fmt.Sprintf("Usage:\n  %s %s", appName, command.Name)
	args := command.Args()
	if args != nil && len(args) > 0 {
		usageText += " [args]"
	}
	usageText += " [options]"
	fmt.Fprintln(tabWriter, usageText)
	fmt.Fprintln(tabWriter)

	if args != nil && len(args) > 0 {
		fmt.Fprintln(tabWriter, "Arguments:")
		for _, arg := range args {
			required := ""
			if arg.Required == 1 {
				required = "(required)"
			}
			fmt.Fprintln(tabWriter, fmt.Sprintf("  %s %s \t %s", arg.Name, required, arg.Description))
		}
		fmt.Fprintln(tabWriter)
	}

	printOptions(tabWriter, "Command options:", command.Options())

	fmt.Fprintln(tabWriter, fmt.Sprintf("Use `%s -h` to view application options", appName))
	tabWriter.Flush()
}

// printOptionsSimple prints options in one line per option group, without the description for each option
func PrintOptionsSimple(output io.Writer, groups []*flags.Group) {
	for _, optionGroup := range groups {
		if len(optionGroup.Groups()) > 0 {
			PrintOptionsSimple(output, optionGroup.Groups())
			continue
		}

		options := optionGroup.Options()
		if options != nil && len(options) > 0 {
			optionUsages := make([]string, len(options))
			for i, option := range options {
				optionUsages[i] = parseOptionUsageText(option, false)
			}

			fmt.Fprintf(output, "%s %s\n", optionGroup.ShortDescription, strings.Join(optionUsages, ", "))
		}
	}
}
