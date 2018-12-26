package commands

import (
	"fmt"
	"github.com/raedahgroup/dcrcli/cli/utils"
)

// SendCommand lets the user send DCR.
type SendCommand struct{}

// Execute runs the `send` command.
func (s SendCommand) Execute(args []string) error {
	return send(false)
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct{}

// Execute runs the `send-custom` command.
func (s SendCustomCommand) Execute(args []string) error {
	return send(true)
}

func send(custom bool) error {
	sourceAccount, err := utils.SelectAccount()
	if err != nil {
		return err
	}

	// check if account has positive non-zero balance before proceeding
	// if balance is zero, there'd be no unspent outputs to use
	accountBalance, err := utils.Wallet.AccountBalance(sourceAccount)
	if err != nil {
		return err
	}
	if accountBalance.Total == 0 {
		return fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

	destinationAddress, err := utils.GetSendDestinationAddress()
	if err != nil {
		return err
	}

	sendAmount, err := utils.GetSendAmount()
	if err != nil {
		return err
	}

	var utxoSelection []string
	if custom {
		// get all utxos in account, pass 0 amount to get all
		utxos, err := utils.Wallet.UnspentOutputs(sourceAccount, 0)
		if err != nil {
			return err
		}

		utxoSelection, err = utils.GetUtxosForNewTransaction(utxos, sendAmount)
		if err != nil {
			return err
		}
	}

	passphrase, err := utils.GetWalletPassphrase()
	if err != nil {
		return err
	}

	var sentTransactionHash string
	if custom {
		sentTransactionHash, err = utils.Wallet.SendFromUTXOs(utxoSelection, sendAmount, sourceAccount, destinationAddress, passphrase)
	} else {
		sentTransactionHash, err = utils.Wallet.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
	}

	if err != nil {
		return err
	}

	fmt.Printf("Sent. Txid: %s\n", sentTransactionHash)
	return nil
}
