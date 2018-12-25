package commands

// CliCommands defines the commands and options available on the cli
type CliCommands struct {
	Balance    BalanceCommand    `command:"balance" description:"show your balance"`
	Receive    ReceiveCommand    `command:"receive" description:"show your address to receive funds"`
	Send       SendCommand       `command:"send" description:"send a transaction"`
	SendCustom SendCustomCommand `command:"send-custom" description:"send a transaction, manually selecting inputs from unspent outputs"`
	History    HistoryCommand    `command:"history" description:"show your transaction history"`
}
