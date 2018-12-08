package cli

import (
	"bytes"
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
func PrintHelp() {
	usagePrefix := "Usage:\n  dcrcli "
	stderrTabWriter := tabWriter(os.Stderr)
	writeHelpMessage(usagePrefix, stderrTabWriter)
}

func writeHelpMessage(prefix string, w *tabwriter.Writer) {
	res := &response{
		columns: []string{prefix + "[OPTIONS] <command> [<args...>]\n\nAvailable commands:"},
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
