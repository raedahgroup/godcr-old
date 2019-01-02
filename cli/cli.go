package cli

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/cli/commands"
	"github.com/raedahgroup/godcr/cli/help"
	"github.com/raedahgroup/godcr/cli/runner"
	"github.com/raedahgroup/godcr/cli/walletloader"
)

// appConfigWithCliCommands is the entrypoint to the cli application.
// It defines general app options, cli commands with their command-specific options and general cli options
type appConfigWithCliCommands struct {
	commands.AvailableCommands
	commands.ExperimentalCommands
	config.Config
	runner.CliOptions
}

// Run starts the app in cli interface mode
func Run(ctx context.Context, walletMiddleware app.WalletMiddleware, appConfig config.Config) error {
	configWithCommands := &appConfigWithCliCommands{
		Config: appConfig,
	}
	parser := flags.NewParser(configWithCommands, flags.HelpFlag|flags.PassDoubleDash)

	// use command handler wrapper function to provide wallet dependency injection to command handlers at execution time
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		commandRunner := runner.New(ctx, walletMiddleware)
		return commandRunner.Run(parser, command, args, configWithCommands.CliOptions)
	}

	// parser.Parse invokes parser.CommandHandler if a command is provided
	// returns an error of type ErrCommandRequired
	_, err := parser.Parse()
	noCommandPassed := config.IsFlagErrorType(err, flags.ErrCommandRequired)
	helpFlagPassed := config.IsFlagErrorType(err, flags.ErrHelp)

	// if no command is passed but --sync flag was passed, perform sync operation and return
	if noCommandPassed && configWithCommands.CliOptions.SyncBlockchain {
		return syncBlockChain(ctx, walletMiddleware)
	}

	if noCommandPassed {
		displayAvailableCommandsHelpMessage(parser)
	} else if helpFlagPassed {
		displayHelpMessage(parser.Name, parser.Active)
	} else if err != nil {
		fmt.Println(err)
	}

	return err
}

func syncBlockChain(ctx context.Context, walletMiddleware app.WalletMiddleware) error {
	walletExists, err := walletloader.OpenWallet(ctx, walletMiddleware)
	if err != nil || !walletExists {
		return err
	}

	return walletloader.SyncBlockChain(ctx, walletMiddleware)
}

// displayAvailableCommandsHelpMessage prints a simple list of available commands when godcr is run without any command
func displayAvailableCommandsHelpMessage(parser *flags.Parser) {
	registeredCommands := parser.Commands()
	commandNames := make([]string, 0, len(registeredCommands))
	for _, command := range registeredCommands {
		commandNames = append(commandNames, command.Name)
	}
	sort.Strings(commandNames)
	fmt.Fprintln(os.Stderr, "Available Commands: ", strings.Join(commandNames, ", "))
}

func displayHelpMessage(appName string, activeCommand *flags.Command) {
	// help parser should use default app config values, not parsed values
	configWithCommands := &appConfigWithCliCommands{
		Config: config.DefaultConfig(),
	}
	helpParser := flags.NewParser(configWithCommands, flags.HelpFlag|flags.PassDoubleDash)

	if activeCommand == nil {
		help.PrintGeneralHelp(os.Stdout, helpParser, commands.Categories())
	} else {
		help.PrintCommandHelp(os.Stdout, appName, activeCommand)
	}
}
