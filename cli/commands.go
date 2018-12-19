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

type ReceiveCommand struct{}

func (r ReceiveCommand) Execute(args []string) error {
	res, err := receive(walletClient, args)
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
