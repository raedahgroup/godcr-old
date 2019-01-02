package commands

import (
	"reflect"

	"github.com/raedahgroup/godcr/cli/help"
)

// AvailableCommands defines thoroughly-tested commands and options available on the cli
type AvailableCommands struct {
	CreateWallet    CreateWalletCommand    `command:"create-wallet" description:"Creates a new decred testnet or mainnet wallet" long-description:"Creates a new decred testnet or mainnet wallet. A wallet seed will be generated for the new wallet which must be stored securely. You'll also be asked to set a password for the wallet"`
	Balance         BalanceCommand         `command:"balance" description:"Show total balance for each account in wallet" long-description:"Also shows spendable balance if different from total balance"`
	Send            SendCommand            `command:"send" description:"Send a transaction"`
	Receive         ReceiveCommand         `command:"receive" description:"Show your address to receive funds"`
	History         HistoryCommand         `command:"history" description:"Show your transaction history"`
	ShowTransaction ShowTransactionCommand `command:"show-transaction" description:"Show details of a transaction"`
	Help            HelpCommand            `command:"help" description:"Show general application help. Run help <command-name> to get help message for a specific command"`
}

// AvailableCommands defines experimental commands and options available on the cli
type ExperimentalCommands struct {
	SendCustom SendCustomCommand `command:"send-custom" description:"Send a transaction, manually selecting inputs from unspent outputs"`
}

// Categories return information for the different categories of commands defined in this file
func Categories() []*help.CommandCategory {
	parseCommandNames := func(commandCategory interface{}) (commandNames []string) {
		commandData := reflect.ValueOf(commandCategory).Elem()
		dataType := commandData.Type()

		for i := 0; i < commandData.NumField(); i++ {
			commandNames = append(commandNames, dataType.Field(i).Tag.Get("command"))
		}
		return
	}

	return []*help.CommandCategory{
		{Name: "Available Commands", ShortName: "Commands", CommandNames: parseCommandNames(&AvailableCommands{})},
		{Name: "Experimental Commands", ShortName: "Experimental", CommandNames: parseCommandNames(&ExperimentalCommands{})},
	}
}

// commanderStub implements `flags.Commander`, using a noop Execute method to satisfy `flags.Commander` interface
// Commands embedding this struct should ideally implement any of the interfaces in `runners.go`
// The `Run` method of such commands will be invoked by `CommandRunner.Run`, providing specific dependencies required by the command
type commanderStub struct{}

// Noop Execute method added to satisfy `flags.Commander` interface
func (w commanderStub) Execute(args []string) error {
	return nil
}
