package commands

import (
	"fmt"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// SendCommand lets the user send DCR.
type SendCommand struct {
	commanderStub
}

// Run runs the `send` command.
func (s SendCommand) Run(wallet walletcore.Wallet) error {
	return send(wallet, false)
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct {
	commanderStub
}

// Run runs the `send-custom` command.
func (s SendCustomCommand) Run(wallet walletcore.Wallet) error {
	return send(wallet, true)
}

func send(wallet walletcore.Wallet, custom bool) (err error) {
	sourceAccount, err := selectAccount(wallet)
	if err != nil {
		return err
	}

	// check if account has positive non-zero balance before proceeding
	// if balance is zero, there'd be no unspent outputs to use
	accountBalance, err := wallet.AccountBalance(sourceAccount)
	if err != nil {
		return err
	}
	if accountBalance.Total == 0 {
		return fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

	destinationAddress, err := getSendDestinationAddress(wallet)
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
		utxos, err := wallet.UnspentOutputs(sourceAccount, 0)
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

	fmt.Printf("You are about to send %f DCR to %s\n", sendAmount, destinationAddress)
	sendConfirmed, err := terminalprompt.RequestYesNoConfirmation("Are you sure?", "")
	if err != nil {
		return fmt.Errorf("error reading your response: %s", err.Error())
	}

	if !sendConfirmed {
		fmt.Println("Canceled")
		return nil
	}

	sendDestinations := []txhelper.TransactionDestination{{
		Amount:  sendAmount,
		Address: destinationAddress,
	}}

	var sentTransactionHash string
	if custom {
		sentTransactionHash, err = wallet.SendFromUTXOs(sourceAccount, utxoSelection, sendDestinations, passphrase)
	} else {
		sentTransactionHash, err = wallet.SendFromAccount(sourceAccount, sendDestinations, passphrase)
	}

	if err != nil {
		return err
	}

	fmt.Println("Sent. Txid", sentTransactionHash)
	return nil
}
