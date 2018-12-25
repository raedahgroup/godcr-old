package main

import 	(
	"fmt"
	"os"
	"strings"

	flags "github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/cli"
	"github.com/raedahgroup/godcr/config"
	"github.com/raedahgroup/godcr/desktop"
	ws "github.com/raedahgroup/godcr/walletsource"
	"github.com/raedahgroup/godcr/walletsource/dcrwalletrpc"
	"github.com/raedahgroup/godcr/walletsource/mobilewalletlib"
	"github.com/raedahgroup/godcr/web"

	"github.com/raedahgroup/dcrcli/app"
	"github.com/raedahgroup/dcrcli/app/config"
	"github.com/raedahgroup/dcrcli/app/walletmediums/dcrwalletrpc"
	"github.com/raedahgroup/dcrcli/app/walletmediums/mobilewalletlib"
	"github.com/raedahgroup/dcrcli/cli"
	"github.com/raedahgroup/dcrcli/web"
)

func main() {
	args, appConfig, parser, err := config.LoadConfig(true)
	if err != nil {
		handleParseError(err, parser)
	appConfig := config.Default()

	// create parser to parse flags/options from config and commands
	parser := flags.NewParser(&commands.CliCommands{Config: appConfig}, flags.HelpFlag)

	// continueExecution will be false if an error is encountered while parsing or if `-h` or `-v` is encountered
	continueExecution := config.ParseConfig(appConfig, parser)
	if !continueExecution {
	appConfig := config.LoadConfig()
	if appConfig == nil {
		os.Exit(1)
	}

	wallet := connectToWallet(appConfig)

	if appConfig.HTTPMode {
		if len(args) > 0 {
			fmt.Println("unexpected command or flag:", strings.Join(args, " "))
			os.Exit(1)
		}
		enterHttpMode(appConfig.HTTPServerAddress, wallet)
	} else if appConfig.DesktopMode {
		enterDesktopMode(wallet)
		web.StartHttpServer(wallet, appConfig.HTTPServerAddress)
	} else {
		cli.Run(wallet, appConfig)
	}
}

// connectToWallet opens connection to a wallet via any of the available walletmiddleware
// default walletmiddleware is mobilewallet library, alternative is dcrwallet rpc
func connectToWallet(config *config.Config) app.WalletMiddleware {
	var netType string
	if config.UseTestNet {
		netType = "testnet"
	} else {
		netType = "mainnet"
	}

	if !config.UseWalletRPC {
		return mobilewalletlib.New(config.AppDataDir, netType)
	}

	wallet, err := dcrwalletrpc.New(config.WalletRPCServer, config.RPCCert, config.NoDaemonTLS, config.TestNet)
	if err != nil {
		fmt.Println("Connect to dcrwallet rpc failed")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return walletMiddleware
}

func enterHttpMode(serverAddress string, wallet core.Wallet) {
	fmt.Println("Running in http mode")
	web.StartHttpServer(serverAddress, wallet)
}

func enterDesktopMode(walletsource ws.WalletSource) {
	fmt.Println("Running in desktop mode")
	desktop.StartDesktopApp(walletsource)
}

func enterCliMode(appConfig *config.Config, wallet core.Wallet) {
	// todo: correct comment Set the walletrpcclient.Client object that will be used by the command handlers
	cli.Wallet = wallet

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

func enterCliMode(appConfig config.Config, walletsource ws.WalletSource) {
	cli.WalletSource = walletsource

	if appConfig.CreateWallet {
		// perform first blockchain sync after creating wallet
		cli.CreateWallet()
		appConfig.SyncBlockchain = true
	}

	if appConfig.SyncBlockchain {
		// open wallet then sync blockchain, before executing command
		cli.OpenWallet()
		cli.SyncBlockChain()
	}

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
