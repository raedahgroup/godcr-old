package commands

import (
	"fmt"

	"github.com/raedahgroup/dcrcli/cli/termio"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

// SendCommand lets the user send DCR.
type SendCommand struct{}

// Execute is a stub method to satisfy the commander interface, so that
// it can be passed to the custom command handler which will inject the
// necessary dependencies to run the command.
func (h SendCommand) Execute(args []string) error {
	return nil
}

// Execute runs the `send` command.
func (s SendCommand) Run(client *walletrpcclient.Client, args []string) error {
	res, err := send(client, false)
	if err != nil {
		return err
	}
	termio.PrintResult(termio.StdoutWriter, res)
	return nil
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct{}

// Execute is a stub method to satisfy the commander interface, so that
// it can be passed to the custom command handler which will inject the
// necessary dependencies to run the command.
func (h SendCustomCommand) Execute(args []string) error {
	return nil
}

// Execute runs the `send-custom` command.
func (s SendCustomCommand) Run(client *walletrpcclient.Client, args []string) error {
	res, err := send(client, true)
	if err != nil {
		return err
	}
	termio.PrintResult(termio.StdoutWriter, res)
	return nil
}

func send(rpcclient *walletrpcclient.Client, custom bool) (*termio.Response, error) {
	var err error

	sourceAccount, err := termio.GetSendSourceAccount(rpcclient)
	if err != nil {
		return nil, err
	}

	// check if account has positive non-zero balance before proceeding
	// if balance is zero, there'd be no unspent outputs to use
	accountBalance, err := rpcclient.SingleAccountBalance(sourceAccount, nil)
	if err != nil {
		return nil, err
	}
	if accountBalance.Total == 0 {
		return nil, fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

	destinationAddress, err := termio.GetSendDestinationAddress(rpcclient)
	if err != nil {
		return nil, err
	}

	sendAmount, err := termio.GetSendAmount()
	if err != nil {
		return nil, err
	}

	var utxoSelection []string
	if custom {
		// get all utxos in account, pass 0 amount to get all
		utxos, err := rpcclient.UnspentOutputs(sourceAccount, 0)
		if err != nil {
			return nil, err
		}

		utxoSelection, err = termio.GetUtxosForNewTransaction(utxos, sendAmount)
		if err != nil {
			return nil, err
		}
	}

	passphrase, err := termio.GetWalletPassphrase()
	if err != nil {
		return nil, err
	}

	var result *walletrpcclient.SendResult
	if custom {
		result, err = rpcclient.SendFromUTXOs(utxoSelection, sendAmount, sourceAccount,
			destinationAddress, passphrase)
	} else {
		result, err = rpcclient.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
	}

	if err != nil {
		return nil, err
	}

	res := &termio.Response{
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
