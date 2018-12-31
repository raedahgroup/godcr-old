package commands

import (
	"context"
	"fmt"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/runner"
	"github.com/raedahgroup/godcr/cli/termio"
	qrcode "github.com/skip2/go-qrcode"
)

// ReceiveCommand generates an address for a user to receive DCR.
type ReceiveCommand struct {
	runner.WalletCommand
	Args ReceiveCommandArgs `positional-args:"yes"`
}
type ReceiveCommandArgs struct {
	AccountName string `positional-arg-name:"account-name"`
}

// Run runs the `receive` command.
func (receiveCommand ReceiveCommand) Run(ctx context.Context, wallet walletcore.Wallet, args []string) error {
	var accountNumber uint32
	// if no account name was passed in
	if receiveCommand.Args.AccountName == "" {
		// display menu options to select account
		var err error
		accountNumber, err = selectAccount(wallet)
		if err != nil {
			return err
		}
	} else {
		// if an account name was passed in e.g. ./godcr receive default
		// get the address corresponding to the account name and use it
		var err error
		accountNumber, err = wallet.AccountNumber(receiveCommand.Args.AccountName)
		if err != nil {
			return fmt.Errorf("Error fetching account number: %s", err.Error())
		}
	}

	receiveAddress, err := wallet.GenerateReceiveAddress(accountNumber)
	if err != nil {
		return err
	}

	qr, err := qrcode.New(receiveAddress, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("Error generating QR Code: %s", err.Error())
	}

	columns := []string{
		"Address",
		"QR Code",
	}
	rows := [][]interface{}{
		[]interface{}{
			receiveAddress,
			qr.ToString(true),
		},
	}
	termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
	return nil
}
