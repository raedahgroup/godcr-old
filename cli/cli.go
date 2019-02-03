package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/app/help"
	"github.com/raedahgroup/godcr/cli/clilog"
	"github.com/raedahgroup/godcr/cli/commands"
	"github.com/raedahgroup/godcr/cli/runner"
	"github.com/raedahgroup/godcr/cli/walletloader"
)

// AppConfigWithCliCommands is the entrypoint to the cli application.
// It defines general app options, cli commands with their command-specific options and general cli options
type AppConfigWithCliCommands struct {
	commands.AvailableCommands
	commands.ExperimentalCommands
	config.Config
}

// Run starts the app in cli interface mode
func Run(ctx context.Context, walletMiddleware app.WalletMiddleware, appConfig config.Config) error {
	configWithCommands := &AppConfigWithCliCommands{
		Config: appConfig,
	}
	parser := flags.NewParser(configWithCommands, flags.None)

	// use command handler wrapper function to provide wallet dependency injection to command handlers at execution time
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		commandRunner := runner.New(parser, ctx, walletMiddleware)
		return commandRunner.Run(command, args, configWithCommands.CliOptions)
	}

	// parser.Parse invokes parser.CommandHandler if a command is provided
	// returns an error of type ErrCommandRequired if no command is passed
	_, err := parser.Parse()
	noCommandPassed := config.IsFlagErrorType(err, flags.ErrCommandRequired)

	// if no command is passed but --sync flag was passed, perform sync operation and return
	if noCommandPassed && configWithCommands.CliOptions.SyncBlockchain {
		return syncBlockChain(ctx, walletMiddleware)
	} else if noCommandPassed {
		listCommands()
	} else if err != nil {
		clilog.LogError(err)
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

// listCommands prints a simple list of available commands when godcr is run without any command
func listCommands() {
	help.PrintOptionsSimple(os.Stdout, commands.HelpParser().Groups())
	for _, category := range commands.Categories() {
		fmt.Fprintf(os.Stderr, "%s: %s\n", category.ShortName, strings.Join(category.CommandNames, ", "))
	}
}
