package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/raedahgroup/godcr/walletrpcclient"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
	"github.com/mdp/qrterminal"
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

	// Print out address as string
	fmt.Println(receiveResult.Address)

	// Print out QR code
	validateConfirm := func(address string) error {
		return  nil
	}
	confirm, _ := terminalprompt.RequestInput("Would you like to a generate QR code? (y/n) ", validateConfirm)

	if confirm == "Yes" || confirm == "yes" {
		confirm = "y"
	}

	if confirm == "y" {
		qrterminal.GenerateHalfBlock("https://github.com/mdp/qrterminal", qrterminal.L, os.Stdout)
	}

	return nil
}
