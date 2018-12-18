package cli

// command is an action that can be requested at the cli
type command struct {
	name        string
	description string
	handler     handler
	experimental bool
}

// supportedCommands provides the commands available from the cli
func supportedCommands() []command {
	return []command{
		{"balance", "show your balance", balance, false},
		{"send", "send a transaction", normalSend, false},
		{"send-custom", "send a transaction, manually selecting inputs from unspent outputs", customSend, true},
		{"receive", "show your address to receive funds", receive, false},
	}
}
