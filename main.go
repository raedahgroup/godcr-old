package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/dcrcli/cli"
	"github.com/raedahgroup/dcrcli/config"
	ws "github.com/raedahgroup/dcrcli/walletsource"
	"github.com/raedahgroup/dcrcli/walletsource/dcrwalletrpc"
	"github.com/raedahgroup/dcrcli/walletsource/mobilewalletlib"
	"github.com/raedahgroup/dcrcli/web"
)

func main() {
	appConfig, parser, err := config.LoadConfig()
	if err != nil {
		handleParseError(err, parser)
		os.Exit(1)
	}

	walletSource := makeWalletSource(appConfig)

	if config.HTTPMode {
		enterHttpMode(appConfig.HTTPServerAddress, walletSource)
	} else {
<<<<<<< HEAD
		enterCliMode(appConfig, walletSource)
=======
		enterCliMode(appName, walletSource, args, config.SyncBlockchain)
>>>>>>> cli and web interface functionailty restored
	}
}

// makeWalletSource opens connection to a wallet via the selected source/medium
// default is mobile wallet library, alternative is dcrwallet rpc
func makeWalletSource(config *config.Config) ws.WalletSource {
	var walletSource ws.WalletSource
	var err error

	if config.UseWalletRPC {
		walletSource, err = dcrwalletrpc.New(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS)
		if err != nil {
			fmt.Println("Connect to dcrwallet rpc failed")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else {
		var netType string
		if config.TestNet {
			netType = "testnet"
		} else {
			netType = "mainnet"
		}

		walletSource = mobilewalletlib.New(config.AppDataDir, netType)
	}

	return walletSource
}

func enterHttpMode(serverAddress string, walletsource ws.WalletSource) {
	fmt.Println("Running in http mode")
	web.StartHttpServer(serverAddress, walletsource)
}

<<<<<<< HEAD
func enterCliMode(appConfig *config.Config, walletsource ws.WalletSource) {
	// todo: correct comment Set the walletrpcclient.Client object that will be used by the command handlers
	cli.WalletSource = walletsource

	parser := flags.NewParser(appConfig, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		if config.IsFlagErrorType(err, flags.ErrCommandRequired) {
			// No command was specified, print the available commands.
			availableCommands := supportedCommands(parser)
			fmt.Fprintln(os.Stderr, "Available Commands: ", strings.Join(availableCommands, ", "))
		} else {
			handleParseError(err, parser)
		}
		os.Exit(1)
	}
}

func supportedCommands(parser *flags.Parser) []string {
	registeredCommands := parser.Commands()
	commandNames := make([]string, 0, len(registeredCommands))
	for _, command := range registeredCommands {
		commandNames = append(commandNames, command.Name)
	}
	sort.Strings(commandNames)
	return commandNames
}

func handleParseError(err error, parser *flags.Parser) {
	if err == nil {
		return
	}
	if (parser.Options & flags.PrintErrors) != flags.None {
		// error printing is already handled by go-flags.
		return
	}
	if !config.IsFlagErrorType(err, flags.ErrHelp) {
		fmt.Println(err)
	} else if parser.Active == nil {
		// Print help for the root command (general help with all the options and commands).
		parser.WriteHelp(os.Stderr)
	} else {
		// Print a concise command-specific help.
		printCommandHelp(parser.Name, parser.Active)
	}
}

func printCommandHelp(appName string, command *flags.Command) {
	helpParser := flags.NewParser(nil, flags.HelpFlag)
	helpParser.Name = appName
	helpParser.Active = command
	helpParser.WriteHelp(os.Stderr)
	fmt.Printf("To view application options, use '%s -h'\n", appName)
=======
func enterCliMode(appName string, walletsource ws.WalletSource, args []string, shouldSyncBlockchain bool) {
	c := cli.New(walletsource, appName)
	c.RunCommand(args, shouldSyncBlockchain)
>>>>>>> cli and web interface functionailty restored
}
