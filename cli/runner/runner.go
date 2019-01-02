package runner

import (
	"context"
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
	"github.com/raedahgroup/godcr/cli/walletloader"
)

type CommandRunner struct {
	ctx              context.Context
	walletMiddleware app.WalletMiddleware
}

func New(ctx context.Context, walletMiddleware app.WalletMiddleware) *CommandRunner {
	return &CommandRunner{
		ctx,
		walletMiddleware,
	}
}

// Run checks if a command implements `IWalletRunner` and executes the command using the command's Run method
// Other commands are executed using the Execute method implemented by those commands
// If the command does not implement either Run or Execute method, a broken command error is returned
func (runner CommandRunner) Run(parser *flags.Parser, command flags.Commander, args []string, options config.CliOptions) error {
	if command == nil {
		return brokenCommandError(parser.Command)
	}

	// attempt to run the command by injecting walletMiddleware dependency
	if commandRunner, ok := command.(WalletMiddlewareCommandRunner); ok {
		return commandRunner.Run(runner.ctx, runner.walletMiddleware)
	}

	// attempt to run the command by injecting wallet dependencies
	if commandRunner, ok := command.(WalletCommandRunner); ok {
		return runner.processWalletCommand(commandRunner, args, options)
	}

	// try running the command by injecting parser dependency
	if commandRunner, ok := command.(ParserCommandRunner); ok {
		return commandRunner.Run(parser)
	}

	return command.Execute(args)
}

// processWalletCommand handles command execution for commands requiring access to the decred wallet
// Such commands must implement `WalletCommandRunner` by providing a Run function
// The wallet is opened using the provided walletMiddleware, sync operations performed (if requested)
// then, the command is executed using the Run method of the WalletCommandRunner interface
func (runner CommandRunner) processWalletCommand(commandRunner WalletCommandRunner, args []string, options config.CliOptions) error {
	walletExists, err := walletloader.OpenWallet(runner.ctx, runner.walletMiddleware)
	if err != nil || !walletExists {
		return err
	}

	if options.SyncBlockchain {
		err = walletloader.SyncBlockChain(runner.ctx, runner.walletMiddleware)
		if err != nil {
			return err
		}
	}

	return commandRunner.Run(runner.walletMiddleware)
}

func brokenCommandError(command *flags.Command) error {
	return fmt.Errorf("The command %q was not properly setup.\n"+
		"Please report this bug at https://github.com/raedahgroup/godcr/issues",
		commandName(command))
}

func commandName(command *flags.Command) string {
	name := command.Name
	if command.Active != nil {
		return fmt.Sprintf("%s %s", name, commandName(command.Active))
	}
	return name
}
