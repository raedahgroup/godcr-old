package help

import (
	"fmt"
	"io"
	"reflect"
	"sort"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/cli/termio"
)

type CommandCategory struct {
	Name string
	ShortName string
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

	// print usage
	fmt.Fprintf(tabWriter, "Usage:\n  %s [options] <command> [args]\n", parser.Name)
	fmt.Fprintln(tabWriter)

	// print general app options
	for _, optionGroup := range parser.Groups() {
		printOptions(tabWriter, optionGroup.ShortDescription, optionGroup.Options())
	}

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

// option printout attempts to add 2 whitespace for options with short name and 6 for those without
// This is an attempt to stay consistent with the output of parser.WriteHelp
func printOptions(tabWriter io.Writer, optionDescription string, options []*flags.Option) {
	if options != nil && len(options) > 0 {
		fmt.Fprintln(tabWriter, optionDescription)

		for _, option := range options {
			var optionUsage string

			if option.ShortName != 0 && option.LongName != "" {
				optionUsage = fmt.Sprintf("  -%c, --%s", option.ShortName, option.LongName)
			} else if option.ShortName != 0 {
				optionUsage = fmt.Sprintf("  -%c", option.ShortName)
			} else {
				optionUsage = fmt.Sprintf("      --%s", option.LongName)
			}

			if option.Field().Type.Kind() != reflect.Bool {
				optionUsage += "="
			}

			description := option.Description
			optionDefaultValue := reflect.ValueOf(option.Value())
			if optionDefaultValue.Kind() == reflect.String {
				description += fmt.Sprintf(" (default: %s)", optionDefaultValue.String())
			}

			fmt.Fprintln(tabWriter, fmt.Sprintf("%s \t %s", optionUsage, description))
		}

		fmt.Fprintln(tabWriter)
	}
}

func printCommands(tabWriter io.Writer, commandGroups map[string][]*flags.Command) {
	// sort first to ensure consistent display order
	categories := make([]string, 0, len(commandGroups))
	for category := range commandGroups {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	for _, category := range categories {
		fmt.Fprintf(tabWriter, "%s:\n", category)

		for _, command := range commandGroups[category] {
			fmt.Fprintln(tabWriter, fmt.Sprintf("  %s \t %s", command.Name, command.ShortDescription))
		}

		fmt.Fprintln(tabWriter)
	}
}
