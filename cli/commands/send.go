package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/runner"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// SendCommand lets the user send DCR.
type SendCommand struct {
	runner.WalletCommand
}

// Run runs the `send` command.
func (s SendCommand) Run(ctx context.Context, wallet walletcore.Wallet, args []string) error {
	return send(wallet, false)
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct {
	runner.WalletCommand
}

// Run runs the `send-custom` command.
func (s SendCustomCommand) Run(ctx context.Context, wallet walletcore.Wallet, args []string) error {
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

	validateConfirm := func(userResponse string) error {
		userResponse = strings.TrimSpace(userResponse)
		userResponse = strings.Trim(userResponse, `"`)
		if userResponse == "" || strings.EqualFold("Y", userResponse) || strings.EqualFold("n", userResponse) {
			return nil
		} else {
			return fmt.Errorf("invalid option, try again")
		}
	}
	confirm, err := terminalprompt.RequestInput("Are you sure? (y/N) ", validateConfirm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading your response: %s", err.Error())
		return err
	}

	if strings.EqualFold(confirm, "yes") || strings.EqualFold(confirm, "y") {
		confirm = "y"
	}

	if confirm != "y" {
		fmt.Printf("Operation cancelled\n")
		return nil
	}

	var sentTransactionHash string

	if custom {
		sentTransactionHash, err = wallet.SendFromUTXOs(utxoSelection, sendAmount, sourceAccount, destinationAddress, passphrase)
	} else {
		sentTransactionHash, err = wallet.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
	}

	if err != nil {
		return err
	}

	fmt.Println("Sent. Txid", sentTransactionHash)
	return nil
}
