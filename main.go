package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/cli"
	"github.com/raedahgroup/godcr/config"
	"github.com/raedahgroup/godcr/walletrpcclient"
	"github.com/raedahgroup/godcr/web"
)

func main() {
	args, appConfig, parser, err := config.LoadConfig(true)
	if err != nil {
		handleParseError(err, parser)
		os.Exit(1)
	}
	if appConfig.ShowVersion {
		fmt.Println(config.AppVersion())
		os.Exit(0)
	}

	client, err := walletrpcclient.New(appConfig.WalletRPCServer, appConfig.RPCCert, appConfig.NoDaemonTLS, appConfig.TestNet)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to RPC server")
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if appConfig.HTTPMode {
		if len(args) > 0 {
			fmt.Println("cannot use --http with a command")
			os.Exit(1)
		}
		enterHTTPMode(appConfig.HTTPServerAddress, client)
	} else {
		enterCliMode(appConfig, client)
	}
}

func enterHTTPMode(serverAddress string, client *walletrpcclient.Client) {
	fmt.Println("Running in http mode")
	web.StartHttpServer(serverAddress, client)
}

func enterCliMode(appConfig config.Config, client *walletrpcclient.Client) {
	appRoot := cli.CliRoot{Config: appConfig}

	parser := flags.NewParser(&appRoot, flags.HelpFlag|flags.PassDoubleDash)
	parser.CommandHandler = cli.CommandHandlerWrapper(client)
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
