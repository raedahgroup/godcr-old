package commands

// Commands defines the commands and options available on the cli
type Commands struct {
	CreateWallet    CreateWalletCommand    `command:"create-wallet" description:"Creates a new decred testnet or mainnet wallet" long-description:"Creates a new decred testnet or mainnet wallet. A wallet seed will be generated for the new wallet which must be stored securely. You'll also be asked to set a password for the wallet"`
	Balance         BalanceCommand         `command:"balance" description:"show your balance"`
	Send            SendCommand            `command:"send" description:"send a transaction"`
	SendCustom      SendCustomCommand      `command:"send-custom" description:"send a transaction, manually selecting inputs from unspent outputs"`
	Receive         ReceiveCommand         `command:"receive" description:"show your address to receive funds"`
	History         HistoryCommand         `command:"history" description:"show your transaction history"`
	ShowTransaction ShowTransactionCommand `command:"show-transaction" description:"show details of a transaction"`
}
