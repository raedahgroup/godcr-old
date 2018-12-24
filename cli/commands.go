package cli

import "github.com/raedahgroup/dcrcli/config"

type AppCommands struct {
	config.Config
	Balance    BalanceCommand    `command:"balance" description:"show your balance"`
	Send       SendCommand       `command:"send" description:"send a transaction"`
	SendCustom SendCustomCommand `command:"send-custom" description:"send a transaction, manually selecting inputs from unspent outputs"`
	Receive    ReceiveCommand    `command:"receive" description:"show your address to receive funds"`
	History    HistoryCommand    `command:"history" description:"show your transaction history"`
}

type BalanceCommand struct{}

func (b BalanceCommand) Execute(args []string) error {
	res, err := balance(walletClient, args)
	if err != nil {
		return err
	}
	printResult(stdoutWriter, res)
	return nil
}

type SendCommand struct{}

func (s SendCommand) Execute(args []string) error {
	res, err := normalSend(walletClient, args)
	if err != nil {
		return err
	}
	printResult(stdoutWriter, res)
	return nil
}

type SendCustomCommand struct{}

func (s SendCustomCommand) Execute(args []string) error {
	res, err := customSend(walletClient, args)
	if err != nil {
		return err
	}
	printResult(stdoutWriter, res)
	return nil
}

type ReceiveCommand struct {
	Args struct {
		Account string `positional-arg-name:"account"`
	} `positional-args:"yes"`
}

func (r ReceiveCommand) Execute(_ []string) error {
	res, err := receive(walletClient, []string{r.Args.Account})
	if err != nil {
		return err
	}
	printResult(stdoutWriter, res)
	return nil
}

type HistoryCommand struct{}

func (h HistoryCommand) Execute(args []string) error {
	res, err := transactionHistory(walletClient, args)
	if err != nil {
		return err
	}
	printResult(stdoutWriter, res)
	return nil
}
