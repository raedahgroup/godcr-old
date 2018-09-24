package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

const (
	showHelpMessage = "Specify -h to show available options"
	listCmdMessage  = "Specify -l to list available commands"
)

// usage displays the general usage when the help flag is not displayed and
// and an invalid command was specified.  The commandUsage function is used
// instead when a valid command was specified.
func usage(errorMessage string) {
	appName := filepath.Base(os.Args[0])
	appName = strings.TrimSuffix(appName, filepath.Ext(appName))
	fmt.Fprintln(os.Stderr, errorMessage)
	fmt.Fprintln(os.Stderr, "Usage:")
	fmt.Fprintf(os.Stderr, "  %s [OPTIONS] <command> <args...>\n\n",
		appName)
	fmt.Fprintln(os.Stderr, showHelpMessage)
	fmt.Fprintln(os.Stderr, listCmdMessage)
}

func main() {
	config, args, err := loadConfig()
	if err != nil {
		os.Exit(1)
	}

	// check if arguments were supplied
	// if not, exit
	if len(args) < 1 {
		usage("No command specified")
		os.Exit(1)
	}

	// check if command is supported
	command := args[0]
	if !walletrpcclient.IsCommandSupported(command) {
		fmt.Fprintf(os.Stderr, "Unrecognized command %s'\n", command)
		fmt.Fprintln(os.Stderr, listCmdMessage)
		os.Exit(1)
	}

	// connect to grpc server
	conn, err := walletrpcclient.Connect(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to RPC server %s'\n", err.Error())
		os.Exit(1)
	}

	// remaining arguments are options
	remainingArgs := args[1:]
	opts := make([]string, 0, len(remainingArgs))
	for _, opt := range remainingArgs {
		opts = append(opts, opt)
	}

	res, err := walletrpcclient.RunCommand(conn, command, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s'\n", err.Error())
		os.Exit(1)
	}

	printResult(res)
}

func printResult(res *walletrpcclient.Response) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	header := ""
	columnLength := len(res.Columns)

	for i := range res.Columns {
		tab := " \t "
		if columnLength == i+1 {
			tab = " "
		}
		header += res.Columns[i] + tab
	}

	fmt.Fprintln(w, header)
	for _, v := range res.Result {
		fmt.Fprintln(w, v)
	}

	w.Flush()
}

// list all supported commands
func listCommands() {

}
