package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	// get an instance of walletrpcclient without connecting
	// because there are commands we can run at this stage without needing to connect to dcrwallet
	client := walletrpcclient.New()

	serveHTTP := false
	if config.HTTPServerAddress != "" {
		serveHTTP = true
	}

	// check if arguments were supplied
	// if not, run ./dcrcli -l, previously used to exit and show usage info
	if len(args) < 1 && !serveHTTP {
		//Go to usage page
		//usage("No command specified")
		res, err := client.RunCommand("listcommands", nil)
		if err != nil {
			// can never happen at this stage
			fmt.Fprintf(os.Stderr, "Error running command %s'\n", err.Error())
			os.Exit(1)
		}
		printResult(res)
		os.Exit(0)
	}

	var command string

	if !serveHTTP {
		command = args[0]
		if args[0] == "listcommands" {
			res, err := client.RunCommand("listcommands", nil)
			if err != nil {
				// can never happen at this stage
				fmt.Fprintf(os.Stderr, "Error running command %s'\n", err.Error())
				os.Exit(1)
			}
			printResult(res)
			os.Exit(0)
		}
	}

	err = client.Connect(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error connecting to RPC server %s'\n", err.Error())
		os.Exit(1)
	}

	if serveHTTP {
		server := server.New(config.HTTPServerAddress, client)
		server.Serve()
		// control can never go beyond this ????
	}

	// check if command is supported
	if !client.IsCommandSupported(command) {
		fmt.Fprintf(os.Stderr, "Unrecognized command %s'\n", command)
		fmt.Fprintln(os.Stderr, listCmdMessage)
		os.Exit(1)
	}

	// remaining arguments are options
	remainingArgs := args[1:]
	opts := make([]string, 0, len(remainingArgs))
	for _, opt := range remainingArgs {
		opts = append(opts, opt)
	}

	res, err := client.RunCommand(command, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running command %s'\n", err.Error())
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
