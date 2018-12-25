package commands

import (
	"fmt"
	"github.com/raedahgroup/dcrcli/cli/utils"

	"github.com/raedahgroup/godcr/cli/termio"
	ws "github.com/raedahgroup/godcr/walletsource"
	qrcode "github.com/skip2/go-qrcode"
)

// ReceiveCommand generates an address for a user to receive DCR.
type ReceiveCommand struct {
	CommanderStub
	Args struct {
	Args ReceiveCommandArgs `positional-args:"yes"`
}
type ReceiveCommandArgs struct {
	AccountName string `positional-arg-name:"account-name"`
}

// Run runs the `receive` command.
func (r ReceiveCommand) Run(walletsource ws.WalletSource, args []string) error {
	var accountNumber uint32
	// if no account name was passed in
	if receiveCommand.Args.AccountName == "" {
		// display menu options to select account
		var err error
		accountNumber, err = selectAccount(walletsource)
		if err != nil {
			return err
		}
	} else {
		// if an account name was passed in e.g. ./godcr receive default
		// get the address corresponding to the account name and use it
		var err error
		accountNumber, err = walletsource.AccountNumber(r.Args.Account)
		if err != nil {
			return fmt.Errorf("Error fetching account number: %s", err.Error())
		}
	}

	receiveAddress, err := walletsource.GenerateReceiveAddress(accountNumber)
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
