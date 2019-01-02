package help

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"

	"github.com/jessevdk/go-flags"
)

// printOptionGroups checks if the root parser option group has nested option groups and prints all
func printOptionGroups(output io.Writer, groups []*flags.Group) {
	for _, optionGroup := range groups {
		if len(optionGroup.Groups()) > 0 {
			printOptionGroups(output, optionGroup.Groups())
		} else {
			printOptions(output, optionGroup.ShortDescription, optionGroup.Options())
		}
	}
}

// printOptions adds 2 trailing whitespace for options with short name and 6 for those without
// This is an attempt to stay consistent with the output of parser.WriteHelp
func printOptions(tabWriter io.Writer, optionDescription string, options []*flags.Option) {
	if options != nil && len(options) > 0 {
		fmt.Fprintln(tabWriter, optionDescription)

		// check if there's any option in this group with short and long name
		// this will help to decide whether or not to pad options without short name to maintain readability
		var hasOptionsWithShortName bool
		for _, option := range options {
			if option.ShortName != 0 && option.LongName != "" {
				hasOptionsWithShortName = true
				break
			}
		}

		for _, option := range options {
			optionUsage := parseOptionUsageText(option, hasOptionsWithShortName)
			description := parseOptionDescription(option)
			fmt.Fprintln(tabWriter, fmt.Sprintf("  %s \t %s", optionUsage, description))
		}

		fmt.Fprintln(tabWriter)
	}
}

func parseOptionUsageText(option *flags.Option, hasOptionsWithShortName bool) (optionUsage string) {
	if option.ShortName != 0 && option.LongName != "" {
		optionUsage = fmt.Sprintf("-%c, --%s", option.ShortName, option.LongName)
	} else if option.ShortName != 0 {
		optionUsage = fmt.Sprintf("-%c", option.ShortName)
	} else if hasOptionsWithShortName {
		// pad long name with 4 spaces to align with options having short and long names
		optionUsage = fmt.Sprintf("    --%s", option.LongName)
	} else {
		optionUsage = fmt.Sprintf("--%s", option.LongName)
	}

	if option.Field().Type.Kind() != reflect.Bool {
		optionUsage += "="
	}

	if len(option.Choices) > 0 {
		optionUsage += fmt.Sprintf("[%s]", strings.Join(option.Choices, ","))
	}

	return
}

func parseOptionDescription(option *flags.Option) (description string) {
	description = option.Description
	optionDefaultValue := reflect.ValueOf(option.Value())
	if optionDefaultValue.Kind() == reflect.String && optionDefaultValue.String() != "" {
		description += fmt.Sprintf(" (default: %s)", optionDefaultValue.String())
	}
	return
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

		commands := commandGroups[category]
		for _, command := range commands {
			fmt.Fprintln(tabWriter, fmt.Sprintf("  %s \t %s", command.Name, command.ShortDescription))
		}

		fmt.Fprintln(tabWriter)
	}
}
