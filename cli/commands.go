package cli

// command is an action that can be requested at the cli
type command struct {
	name        string
	description string
	usage       string
	handler     handler
	experimental bool
}

// supportedCommands provides the commands available from the cli
func supportedCommands() []*command {
	return []*command{
		{"balance", "show your balance", "dcrcli balance", balance, false},
		{"send", "send a transaction", "dcrcli send", normalSend, false},
		{"send-custom", "send a transaction, manually selecting inputs from unspent outputs", "dcrcli send-custom", customSend, false},
		{"receive", "show your address to receive funds", "dcrcli receive [account]", receive, false},
		{"history", "show your transaction history", "dcrcli history", transactionHistory, false},
		{"help", "show the command line help", "dcrcli help [command]", help, false},

	}
}
