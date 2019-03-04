package runner

import (
	"context"

	flags "github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app"
	"github.com/raedahgroup/godcr/app/config"
)

type CommandRunner struct {
	parser           *flags.Parser
	ctx              context.Context
	walletMiddleware app.WalletMiddleware
}

func New(parser *flags.Parser, ctx context.Context, walletMiddleware app.WalletMiddleware) *CommandRunner {
	return &CommandRunner{
		parser:           parser,
		ctx:              ctx,
		walletMiddleware: walletMiddleware,
	}
}

// Run checks if a command requires some form of access to a decred wallet and injects the wallet dependencies needed by the command
// Dependencies for other commands are provided in `runner.RunNoneWalletCommands`
// If the command does not implement the compulsory Execute method, a broken command error is returned
func (runner *CommandRunner) Run(command flags.Commander, args []string, options config.CliOptions) error {
	if command == nil {
		return brokenCommandError(runner.parser.Command)
	}

	// inject walletMiddleware dependency for commands implementing WalletMiddlewareCommandRunner
	if commandRunner, ok := command.(WalletMiddlewareCommandRunner); ok {
		return commandRunner.Run(runner.ctx, runner.walletMiddleware)
	}

	// inject wallet dependency for commands implementing WalletCommandRunner
	// the decred wallet is prepared for use before such commands are executed
	if commandRunner, ok := command.(WalletCommandRunner); ok {
		walletExists, err := prepareWallet(runner.ctx, runner.walletMiddleware, options)
		if err != nil || !walletExists {
			return err
		}
		return commandRunner.Run(runner.ctx, runner.walletMiddleware)
	}

	return runner.RunNoneWalletCommands(command, args)
}

// RunNoneWalletCommands handles command execution for commands that do not require access to the decred wallet
// Such commands may still require some other dependency and those required dependencies are provided by this function
func (runner *CommandRunner) RunNoneWalletCommands(command flags.Commander, args []string) error {
	// inject parser dependency for commands implementing `ParserCommandRunner`
	if commandRunner, ok := command.(ParserCommandRunner); ok {
		return commandRunner.Run(runner.parser)
	}

	// execute other commands directly
	return command.Execute(args)
}
