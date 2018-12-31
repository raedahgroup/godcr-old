package commands

import (
	"fmt"

	flags "github.com/jessevdk/go-flags"
	"github.com/raedahgroup/godcr/app/walletcore"
)

// CliCommands defines the commands and options available on the cli
type Commands struct {
	Balance         BalanceCommand         `command:"balance" description:"show your balance"`
	Send            SendCommand            `command:"send" description:"send a transaction"`
	SendCustom      SendCustomCommand      `command:"send-custom" description:"send a transaction, manually selecting inputs from unspent outputs"`
	Receive         ReceiveCommand         `command:"receive" description:"show your address to receive funds"`
	History         HistoryCommand         `command:"history" description:"show your transaction history"`
	ShowTransaction ShowTransactionCommand `command:"show-transaction" description:"show details of a transaction"`
}

// WalletCommandRunner defines an interface that application commands dependent on
// walletrpcclient.Client can satisfy in order to be provided their dependencies.
type walletCommandRunner interface {
	Run(wallet walletcore.Wallet, args []string) error
	flags.Commander
}

// CommanderStub implements `flags.Commander`, using a noop Execute method.
// Commands embedding this struct would ideally implement `WalletCommandRunner` so that their `Run` method can
// be invoked by the custom command handler which will inject the necessary dependencies to run the command.
type CommanderStub struct{}

func (c CommanderStub) Execute(args []string) error {
	return nil
}

// commandHandler provides a type name for the command handler to register on flags.Parser
type CommandHandler func(flags.Commander, []string) error

// CommandHandlerWrapper provides a command handler that provides walletrpcclient.Client
// to commands.WalletCommandRunner types. Other command that satisfy flags.Commander and do not
// depend on walletrpcclient.Client will be run as well.
// If the command does not satisfy any of these types, ErrNotSupported will be returned.
func CommandHandlerWrapper(parser *flags.Parser, wallet walletcore.Wallet) CommandHandler {
	return func(command flags.Commander, args []string) error {
		if command == nil {
			return brokenCommandError(parser.Command)
		}
		if commandRunner, ok := command.(walletCommandRunner); ok {
			return commandRunner.Run(wallet, args)
		}
		return command.Execute(args)
	}
}

func brokenCommandError(command *flags.Command) error {
	return fmt.Errorf("The command %q was not properly setup.\n" +
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
