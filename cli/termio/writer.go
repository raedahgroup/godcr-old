package termio

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

//TabWriter creates a tabwriter object that writes tab-aligned text.
func TabWriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.TabIndent)
}

// StdoutWriter writes tab-aligned text to os.Stdout
var StdoutWriter = TabWriter(os.Stdout)

// PrintTabularResult formats and prints the content of `res` to `w`
func PrintTabularResult(w *tabwriter.Writer, columnsHeaders []string, rows [][]interface{}) {
	header := ""
	columnLength := len(columnsHeaders)

	for i, columnHeader := range columnsHeaders {
		tab := " \t "
		if columnLength == i+1 {
			tab = " "
		}
		header += strings.ToUpper(columnHeader) + tab
	}
	fmt.Fprintln(w, header)

	for _, row := range rows {
		rowStr := ""
		for range row {
			rowStr += "%v \t "
		}

		// append empty values in columns without data to maintain tabbed formatting
		for i := len(row); i < columnLength; i++ {
			rowStr += " \t "
		}

		rowStr = strings.TrimSuffix(rowStr, "\t ")
		fmt.Fprintln(w, fmt.Sprintf(rowStr, row...))
	}

	w.Flush()
}

// PrintStringResult prints simple string message(s) to a fresh instance of stdOut tabWriter
func PrintStringResult(output ...string) {
	writer := TabWriter(os.Stdout)
	for _, str := range output {
		fmt.Fprintln(writer, str)
	}
	writer.Flush()
}
