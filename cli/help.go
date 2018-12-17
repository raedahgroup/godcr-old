package cli

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// HelpMessage returns the cli usage message as a string.
func HelpMessage() string {
	buf := bytes.NewBuffer([]byte{})
	writer := tabWriter(buf)
	writeHelpMessage("", writer)
	return buf.String()
}

// PrintHelp outputs help message to os.Stderr
func PrintHelp(appName string) {
	usagePrefix := fmt.Sprintf("Usage:\n  %s ", appName)
	stderrTabWriter := tabWriter(os.Stderr)
	writeHelpMessage(usagePrefix, stderrTabWriter)
}

func writeSimpleHelpMessage() {
	stderrTabWriter := tabWriter(os.Stderr)
	var availableCommands []string
	var experimentalCommands []string

	for _, command := range supportedCommands() {
		if command.experimental == false {
			availableCommands = append(availableCommands, command.name)
		} else {
			experimentalCommands = append(experimentalCommands, command.name)
		}
	}

	fmt.Fprintln(stderrTabWriter, "available cmds: ", strings.Join(availableCommands, ", "))
	fmt.Fprintln(stderrTabWriter, "experimental: ", strings.Join(experimentalCommands, ", "))

	stderrTabWriter.Flush()
}

func writeHelpMessage(prefix string, w *tabwriter.Writer) {
	res := &response{
		columns: []string{prefix + "dcrcli [OPTIONS] <command> [<args...>]\n\nAvailable commands:"},
	}
	commands := supportedCommands()

	for _, command := range commands {
		item := []interface{}{
			command.name,
			command.description,
		}
		res.result = append(res.result, item)
	}

	printResult(w, res)
}
