package commands

import (
	"fmt"

	"github.com/raedahgroup/dcrcli/cli/io"
	"github.com/raedahgroup/dcrcli/cli/walletclient"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

// SendCommand lets the user send DCR.
type SendCommand struct{}

// Execute runs the `send` command.
func (s SendCommand) Execute(args []string) error {
	res, err := send(walletclient.WalletClient, false)
	if err != nil {
		return err
	}
	io.PrintResult(io.StdoutWriter, res)
	return nil
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct{}

// Execute runs the `send-custom` command.
func (s SendCustomCommand) Execute(args []string) error {
	res, err := send(walletclient.WalletClient, true)
	if err != nil {
		return err
	}
	io.PrintResult(io.StdoutWriter, res)
	return nil
}

func send(rpcclient *walletrpcclient.Client, custom bool) (*io.Response, error) {
	var err error

	sourceAccount, err := io.GetSendSourceAccount(rpcclient)
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

	destinationAddress, err := io.GetSendDestinationAddress(rpcclient)
	if err != nil {
		return nil, err
	}

	sendAmount, err := io.GetSendAmount()
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

		utxoSelection, err = io.GetUtxosForNewTransaction(utxos, sendAmount)
		if err != nil {
			return nil, err
		}
	}

	passphrase, err := io.GetWalletPassphrase()
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

	res := &io.Response{
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
