package cli

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
