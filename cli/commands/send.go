package commands

import (
	"fmt"

	"github.com/raedahgroup/dcrcli/cli"
	rpcclient "github.com/raedahgroup/dcrcli/walletrpcclient"
)

// SendCommand lets the user send DCR.
type SendCommand struct{}

// Execute runs the `send` command.
func (s SendCommand) Execute(args []string) error {
	res, err := send(cli.WalletClient, false)
	if err != nil {
		return err
	}
	cli.PrintResult(cli.StdoutWriter, res)
	return nil
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct{}

// Execute runs the `send-custom` command.
func (s SendCustomCommand) Execute(args []string) error {
	res, err := send(cli.WalletClient, true)
	if err != nil {
		return err
	}
	cli.PrintResult(cli.StdoutWriter, res)
	return nil
}

func send(walletrpcclient *rpcclient.Client, custom bool) (*cli.Response, error) {
	var err error

	sourceAccount, err := cli.GetSendSourceAccount(walletrpcclient)
	if err != nil {
		return nil, err
	}

	// check if account has positive non-zero balance before proceeding
	// if balance is zero, there'd be no unspent outputs to use
	accountBalance, err := walletrpcclient.SingleAccountBalance(sourceAccount, nil)
	if err != nil {
		return nil, err
	}
	if accountBalance.Total == 0 {
		return nil, fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

	destinationAddress, err := cli.GetSendDestinationAddress(walletrpcclient)
	if err != nil {
		return nil, err
	}

	sendAmount, err := cli.GetSendAmount()
	if err != nil {
		return nil, err
	}

	var utxoSelection []string
	if custom {
		// get all utxos in account, pass 0 amount to get all
		utxos, err := walletrpcclient.UnspentOutputs(sourceAccount, 0)
		if err != nil {
			return nil, err
		}

		utxoSelection, err = cli.GetUtxosForNewTransaction(utxos, sendAmount)
		if err != nil {
			return nil, err
		}
	}

	passphrase, err := cli.GetWalletPassphrase()
	if err != nil {
		return nil, err
	}

	var result *rpcclient.SendResult
	if custom {
		result, err = walletrpcclient.SendFromUTXOs(utxoSelection, sendAmount, sourceAccount,
			destinationAddress, passphrase)
	} else {
		result, err = walletrpcclient.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
	}

	if err != nil {
		return nil, err
	}

	res := &cli.Response{
		Columns: []string{
			"Result",
			"Hash",
		},
		Result: [][]interface{}{
			[]interface{}{
				"The transaction was published successfully",
				result.TransactionHash,
			},
		},
	}

	return res, nil
}
