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
	writeHelpMessage(writer)
	return buf.String()
}

// PrintHelp outputs help message to os.Stderr
func PrintHelp(appName string) {
	fmt.Printf("Usage:\n  %s [OPTIONS] <command> [<args...>]\n\n", appName)
	stderrTabWriter := tabWriter(os.Stderr)
	writeHelpMessage(stderrTabWriter)
}

func writeHelpMessage(w *tabwriter.Writer) {
	var availableCommands []interface{}
	var experimentalCommands []interface{}

	for _, command := range supportedCommands() {
		if command.experimental == false {
			availableCommands = append(availableCommands, command.name)
		} else {
			experimentalCommands = append(experimentalCommands, command.name)
		}
	}

	availableRowStr := "available cmds: "
	for range availableCommands {
		availableRowStr += "%v, "
	}
	availableRowStr = strings.TrimSuffix(availableRowStr, ", ")
	fmt.Fprintln(w, fmt.Sprintf(availableRowStr, availableCommands...))

	experimentalRowStr := "experimental: "
	for range experimentalCommands {
		experimentalRowStr += "%v, "
	}
	experimentalRowStr = strings.TrimSuffix(experimentalRowStr, ", ")
	fmt.Fprintln(w, fmt.Sprintf(experimentalRowStr, experimentalCommands...))

	w.Flush()
}
