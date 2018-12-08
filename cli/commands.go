package cli

// command is an action that can be requested at the cli
type command struct {
	name        string
	description string
	handler     handler
}

// supportedCommands provides the commands available from the cli
func supportedCommands() []command {
	return []command{
		{"balance", "show your balance", balance},
		{"send", "send a transaction", send},
		{"receive", "show your address to receive funds", receive},
	}
}
