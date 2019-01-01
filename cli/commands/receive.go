package commands

import (
	"fmt"
	"os"

	"github.com/mdp/qrterminal"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// ReceiveCommand generates an address for a user to receive DCR.
type ReceiveCommand struct {
	commanderStub
	Args ReceiveCommandArgs `positional-args:"yes"`
}
type ReceiveCommandArgs struct {
	AccountName string `positional-arg-name:"account-name" description:"The name of the account to receive into"`
}

// Run runs the `receive` command.
func (receiveCommand ReceiveCommand) Run(wallet walletcore.Wallet) error {
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
	fmt.Println(receiveAddress)

	// Print out QR code?
	printQR, err := terminalprompt.RequestYesNoConfirmation("Would you like to generate a QR code?", "N")
	if err != nil {
		return fmt.Errorf("error reading your response: %s", err.Error())
	}

	if printQR {
		qrterminal.GenerateHalfBlock(receiveAddress, qrterminal.L, os.Stdout)
	}

	return nil
}
