package commands

import (
	"github.com/jessevdk/go-flags"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

// CliCommands defines the commands and options available on the cli
type CliCommands struct {
	Balance    BalanceCommand    `command:"balance" description:"show your balance"`
	Send       SendCommand       `command:"send" description:"send a transaction"`
	SendCustom SendCustomCommand `command:"send-custom" description:"send a transaction, manually selecting inputs from unspent outputs"`
	Receive    ReceiveCommand    `command:"receive" description:"show your address to receive funds"`
	History    HistoryCommand    `command:"history" description:"show your transaction history"`
}

// WalletCommandRunner defines an optional interface that application commands dependent on
// walletrpcclient.Client can satisfy in order to be provided their dependencies.
type WalletCommandRunner interface {
	Run(client *walletrpcclient.Client, args []string) error
	flags.Commander
}
