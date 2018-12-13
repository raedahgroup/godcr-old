package main

import 	(
	"fmt"
	"os"
	"sort"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/cli"
	"github.com/raedahgroup/godcr/config"
	"github.com/raedahgroup/godcr/desktop"
	"github.com/raedahgroup/godcr/walletsource"
	"github.com/raedahgroup/godcr/walletsource/dcrwalletrpc"
	"github.com/raedahgroup/godcr/web"
)

func main() {
	args, appConfig, parser, err := config.LoadConfig(true)
	if err != nil {
		handleParseError(err, parser)
		os.Exit(1)
	}

	// open connection to a wallet via the selected source/medium
	// default is mobile wallet
	var walletSource ws.WalletSource

	if appConfig.UseWalletRPC {
		walletSource, err = dcrwalletrpc.New(appConfig.WalletRPCServer, appConfig.RPCCert, appConfig.NoDaemonTLS, appConfig.TestNet)
		if err != nil {
			fmt.Println("Failed to connect to wallet")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	} else {
		err = fmt.Errorf("Only dcrwallet daemon is supported currently")
	}

	if appConfig.HTTPMode {
		if len(args) > 0 {
			fmt.Println("unexpected command or flag:", strings.Join(args, " "))
			os.Exit(1)
		}
		enterHttpMode(appConfig.HTTPServerAddress, walletSource)
	} else if appConfig.DesktopMode {
		enterDesktopMode(walletSource)
	} else {
		enterCliMode(appConfig, walletSource)
	}
}

func enterHttpMode(serverAddress string, walletsource ws.WalletSource) {
	fmt.Println("Running in http mode")
	web.StartHttpServer(serverAddress, walletsource)
}

func enterDesktopMode(walletsource ws.WalletSource) {
	fmt.Println("Running in desktop mode")
	desktop.StartDesktopApp(walletsource)
}

func enterCliMode(appConfig config.Config, walletsource ws.WalletSource) {
	cli.WalletSource = walletsource

	appRoot := cli.Root{Config: appConfig}
	parser := flags.NewParser(&appRoot, flags.HelpFlag|flags.PassDoubleDash)
	parser.CommandHandler = cli.CommandHandlerWrapper(parser, client)
	if _, err := parser.Parse(); err != nil {
		if config.IsFlagErrorType(err, flags.ErrCommandRequired) {
			// No command was specified, print the available commands.
			var availableCommands []string
			if parser.Active != nil {
				availableCommands = supportedCommands(parser.Active)
			} else {
				availableCommands = supportedCommands(parser.Command)
			}
			fmt.Fprintln(os.Stderr, "Available Commands: ", strings.Join(availableCommands, ", "))
		} else {
			handleParseError(err, parser)
		}
		os.Exit(1)
	}
}

func supportedCommands(parser *flags.Command) []string {
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
