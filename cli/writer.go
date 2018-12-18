package cli

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

func tabWriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
}

func printResult(w *tabwriter.Writer, res *response) {
	header := ""
	spaceRow := ""
	columnLength := len(res.columns)

	for i := range res.columns {
		tab := " \t "
		if columnLength == i+1 {
			tab = " "
		}
		header += res.columns[i] + tab
		spaceRow += " " + tab
	}
	fmt.Fprintln(w, header)
	fmt.Fprintln(w, spaceRow)


	var availableCommands []interface{}
	var experimentalCommands []interface{}

	for _, row := range res.result {
		if row[2] == false {
			availableCommands = append(availableCommands, row[0])
		} else {
			experimentalCommands = append(experimentalCommands, row[0])
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
