package commands

import (
	"fmt"

	"github.com/raedahgroup/godcr/cli/termio"
	"github.com/raedahgroup/godcr/walletrpcclient"
)

// SendCommand lets the user send DCR.
type SendCommand struct {
	CommanderStub
}

// Run runs the `send` command.
func (s SendCommand) Run(client *walletrpcclient.Client, args []string) error {
	return send(client, false)
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct {
	CommanderStub
}

// Run runs the `send-custom` command.
func (s SendCustomCommand) Run(client *walletrpcclient.Client, args []string) error {
	return send(client, true)
}

func send(rpcclient *walletrpcclient.Client, custom bool) error {
	var err error

	sourceAccount, err := termio.GetSendSourceAccount(rpcclient)
	if err != nil {
		return err
	}

	// check if account has positive non-zero balance before proceeding
	// if balance is zero, there'd be no unspent outputs to use
	accountBalance, err := rpcclient.SingleAccountBalance(sourceAccount, nil)
	if err != nil {
		return err
	}
	if accountBalance.Total == 0 {
		return fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

	destinationAddress, err := termio.GetSendDestinationAddress(rpcclient)
	if err != nil {
		return err
	}

	sendAmount, err := termio.GetSendAmount()
	if err != nil {
		return err
	}

	var utxoSelection []string
	if custom {
		// get all utxos in account, pass 0 amount to get all
		utxos, err := rpcclient.UnspentOutputs(sourceAccount, 0)
		if err != nil {
			return err
		}

		utxoSelection, err = termio.GetUtxosForNewTransaction(utxos, sendAmount)
		if err != nil {
			return err
		}
	}

	passphrase, err := termio.GetWalletPassphrase()
	if err != nil {
		return err
	}

	var result *walletrpcclient.SendResult
	if custom {
		result, err = rpcclient.SendFromUTXOs(utxoSelection, sendAmount, sourceAccount,
			destinationAddress, passphrase)
	} else {
		result, err = rpcclient.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
	}

	if err != nil {
		return err
	}

	columns := []string{
		"Result",
		"Hash",
	}
	rows := [][]interface{}{
		[]interface{}{
			"The transaction was published successfully",
			result.TransactionHash,
		},
	}

	termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
	return nil
}
