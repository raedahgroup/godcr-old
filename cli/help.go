package cli

import (
	"bytes"
	"fmt"
	"os"
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

func writeHelpMessage(prefix string, w *tabwriter.Writer) {
	res := &response{
		columns: []string{prefix + "[OPTIONS] <command> [<args...>]"},
	}
	commands := supportedCommands()

	for _, command := range commands {
		item := []interface{}{
			command.name,
			command.description,
			command.experimental,
		}
		res.result = append(res.result, item)
	}

	printResult(w, res)
}
