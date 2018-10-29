package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/raedahgroup/dcrcli/server"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

type Version struct {
	Major, Minor, Patch int
	Label               string
	Nick                string
}

var Ver = Version{
	Major: 0,
	Minor: 0,
	Patch: 1,
	Label: "",
}

// CommitHash may be set on the build command line:
// go build -ldflags "-X github.com/decred/dcrdata/version.CommitHash=`git describe --abbrev=8 --long | awk -F "-" '{print $(NF-1)"-"$NF}'`"
var CommitHash string

const AppName string = "dcrcli"

func (v *Version) String() string {
	var hashStr string
	if CommitHash != "" {
		hashStr = "+" + CommitHash
	}
	if v.Label != "" {
		return fmt.Sprintf("%d.%d.%d-%s%s",
			v.Major, v.Minor, v.Patch, v.Label, hashStr)
	}
	return fmt.Sprintf("%d.%d.%d%s",
		v.Major, v.Minor, v.Patch, hashStr)
}

func main() {
	config, args, err := loadConfig()
	if err != nil {
		os.Exit(1)
	}

	if config.Mode == "http" {
		fmt.Println("Running in http mode")
		enterHttpMode(config)
	} else {
		enterCliMode(config, args)
	}
}

func enterHttpMode(config *config) {
	if config.HTTPServerAddress == "" {
		fmt.Println("Cannot start http server. Server address not set")
		os.Exit(1)
	}

	client := walletrpcclient.New()
	err := client.Connect(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to RPC server")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	server.StartHttpServer(config.HTTPServerAddress, client)
}

func enterCliMode(config *config, args []string) {
	client := walletrpcclient.New()

	if len(args) == 0 {
		noCommandReceived(client)
		os.Exit(0)
	}

	command := args[0]
	if command == "-l" {
		showAvailableCommands(client)
	}

	if !client.IsCommandSupported(command) {
		invalidCommandReceived(command)
		os.Exit(1)
	}

	cliExecuteCommand(client, command, config, args)
}

func cliExecuteCommand(client *walletrpcclient.Client, command string, config *config, args []string) {
	// open connection to rpc client
	err := client.Connect(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to RPC server")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// get arguments for this command, where command = args[0]
	commandArgs := args[1:]

	res, err := client.RunCommand(command, commandArgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command '%s %s'", AppName, command)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	printResult(res)
}

func printResult(res *walletrpcclient.Response) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	header := ""
	spaceRow := ""
	columnLength := len(res.Columns)

	for i := range res.Columns {
		tab := " \t "
		if columnLength == i+1 {
			tab = " "
		}
		header += res.Columns[i] + tab
		spaceRow += " " + tab
	}

	fmt.Fprintln(w, header)
	fmt.Fprintln(w, spaceRow)
	for _, row := range res.Result {
		rowStr := ""
		for range row {
			rowStr += "%v \t "
		}

		rowStr = strings.TrimSuffix(rowStr, "\t ")
		fmt.Fprintln(w, fmt.Sprintf(rowStr, row...))
	}

	w.Flush()
}

func showAvailableCommands(client *walletrpcclient.Client) {
	commands := client.ListSupportedCommands()
	printResult(commands)
}

func noCommandReceived(client *walletrpcclient.Client) {
	fmt.Printf("usage: %s [OPTIONS] <command> [<args...>]\n\n", AppName)
	fmt.Printf("available %s commands:\n", AppName)
	showAvailableCommands(client)
	fmt.Printf("\nFor available options, see '%s -h'\n", AppName)
}

func invalidCommandReceived(command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid command. See '%s -h'\n", AppName, command, AppName)
}
