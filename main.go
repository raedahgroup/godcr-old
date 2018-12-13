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
	"github.com/raedahgroup/dcrcli/web"
)

func main() {
	appConfig, parser, err := config.LoadConfig()
	if err != nil {
		handleParseError(err, parser)
		os.Exit(1)
	}

	// open connection to a wallet via the selected source/medium
	// default is mobile wallet
	var walletSource ws.WalletSource

	if appConfig.UseWalletRPC {
		walletSource, err = dcrwalletrpc.New(appConfig.WalletRPCServer, appConfig.RPCCert, appConfig.NoDaemonTLS)
		if err != nil {
			fmt.Println("Failed to connect to wallet")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else {
		err = fmt.Errorf("Only dcrwallet daemon is supported currently")
	}

	if appConfig.HTTPMode {
		enterHttpMode(appConfig.HTTPServerAddress, walletSource)
	} else {
		enterCliMode(appConfig, walletSource)
	}
}

func enterHttpMode(serverAddress string, walletsource ws.WalletSource) {
	fmt.Println("Running in http mode")
	web.StartHttpServer(serverAddress, walletsource)
}

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
}
