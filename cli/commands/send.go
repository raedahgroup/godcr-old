package commands

import (
	"fmt"

	"github.com/raedahgroup/godcr/cli/termio"
	ws "github.com/raedahgroup/godcr/walletsource"
)

// SendCommand lets the user send DCR.
type SendCommand struct {
	CommanderStub
}

// Run runs the `send` command.
func (s SendCommand) Run(walletsource ws.WalletSource, args []string) error {
	return send(walletsource, false)
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct {
	CommanderStub
}

// Run runs the `send-custom` command.
func (s SendCustomCommand) Run(walletsource ws.WalletSource, args []string) error {
	return send(walletsource, true)
}

func send(walletsource ws.WalletSource, custom bool) (err error) {
	sourceAccount, err := selectAccount(walletsource)
	if err != nil {
		return err
	}

	// check if account has positive non-zero balance before proceeding
	// if balance is zero, there'd be no unspent outputs to use
	accountBalance, err := walletsource.AccountBalance(sourceAccount)
	if err != nil {
		return err
	}
	if accountBalance.Total == 0 {
		return fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

	destinationAddress, err := getSendDestinationAddress(walletsource)
	if err != nil {
		return err
	}

	sendAmount, err := getSendAmount()
	if err != nil {
		return err
	}

	var utxoSelection []string
	if custom {
		// get all utxos in account, pass 0 amount to get all
		utxos, err := walletsource.UnspentOutputs(sourceAccount, 0)
		if err != nil {
			return err
		}

		utxoSelection, err = getUtxosForNewTransaction(utxos, sendAmount)
		if err != nil {
			return err
		}
	}

	passphrase, err := getWalletPassphrase()
	if err != nil {
		return err
	}

	var sentTransactionHash string
	if custom {
		sentTransactionHash, err = walletsource.SendFromUTXOs(utxoSelection, sendAmount, sourceAccount, destinationAddress, passphrase)
	} else {
		sentTransactionHash, err = walletsource.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
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
			sentTransactionHash,
		},
	}

	termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
	return nil
}
