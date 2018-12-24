package commands

import "github.com/raedahgroup/dcrcli/config"

type AppCommands struct {
	config.Config
	Balance    BalanceCommand    `command:"balance" description:"show your balance"`
	Send       SendCommand       `command:"send" description:"send a transaction"`
	SendCustom SendCustomCommand `command:"send-custom" description:"send a transaction, manually selecting inputs from unspent outputs"`
	Receive    ReceiveCommand    `command:"receive" description:"show your address to receive funds"`
	History    HistoryCommand    `command:"history" description:"show your transaction history"`
}
