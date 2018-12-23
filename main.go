package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/raedahgroup/dcrcli/cli/core"

	"github.com/jessevdk/go-flags"

	"github.com/raedahgroup/dcrcli/cli"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
	"github.com/raedahgroup/dcrcli/web"
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

var appVersion = fmt.Sprintf("%s version: %s", core.AppName(), Ver.String())

func main() {
	config, parser, err := loadConfig()
	if err != nil {
		handleParseError(err, parser)
		os.Exit(1)
	}
	if config == nil {
		os.Exit(0)
	}

	// Show version and exit if the version flag was specified
	if config.ShowVersion {
		fmt.Println(appVersion)
		os.Exit(0)
	}

	client, err := walletrpcclient.New(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to RPC server")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if config.HTTPMode {
		fmt.Println("Running in http mode")
		enterHTTPMode(config, client)
	} else {
		enterCliMode(config, client)
	}
}

func enterHTTPMode(config *core.Config, client *walletrpcclient.Client) {
	web.StartHttpServer(config.HTTPServerAddress, client)
}

func enterCliMode(config *core.Config, client *walletrpcclient.Client) {
	cli.Setup(client)
	parser := flags.NewParser(&cli.DcrcliCommands, flags.Default&(^flags.PrintErrors))
	_, err := parser.Parse()
	if isFlagErrorType(err, flags.ErrCommandRequired) {
		commands := supportedCommands(parser)
		fmt.Fprintln(os.Stderr, "Available Commands: ", strings.Join(commands, ", "))
	} else {
		handleParseError(err, parser)
	}
	os.Exit(1)
}

func supportedCommands(parser *flags.Parser) []string {
	registeredCommands := parser.Commands()
	commands := make([]string, 0, len(registeredCommands))
	for _, command := range registeredCommands {
		commands = append(commands, command.Name)
	}
	sort.Strings(commands)
	return commands
}

func loadConfig() (*core.Config, *flags.Parser, error) {
	// load defaults first
	commands := core.DefaultConfig()

	parser := flags.NewParser(&commands, flags.HelpFlag)

	_, err := parser.Parse()
	if err != nil && !isFlagErrorType(err, flags.ErrHelp) {
		return nil, parser, err
	}

	if commands.ShowVersion {
		return nil, parser, fmt.Errorf(appVersion)
	}

	// Load additional config from file
	err = flags.NewIniParser(parser).ParseFile(commands.ConfigFile)
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return nil, parser, fmt.Errorf("Error parsing configuration file: %v", err.Error())
		}
		return nil, parser, err
	}

	// Parse command line options again to ensure they take precedence.
	_, err = parser.Parse()
	if err != nil && !isFlagErrorType(err, flags.ErrHelp) {
		return nil, parser, err
	}

	return &commands, parser, nil
}

func isFlagErrorType(err error, errorType flags.ErrorType) bool {
	if err == nil {
		return false
	}
	if flagErr, ok := err.(*flags.Error); ok && flagErr.Type == errorType {
		return true
	}
	return false
}

func handleParseError(err error, parser *flags.Parser) {
	if err == nil {
		return
	}
	if (parser.Options & flags.PrintErrors) != flags.None {
		// error printing is already handled by go-flags.
		return
	}
	if isFlagErrorType(err, flags.ErrHelp) {
		parser.WriteHelp(os.Stderr)
	} else {
		fmt.Println(err)
	}
}
