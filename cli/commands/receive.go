package commands

import (
	"fmt"
	"github.com/raedahgroup/dcrcli/cli/utils"

	qrcode "github.com/skip2/go-qrcode"
)

// ReceiveCommand generates and address for a user to receive DCR.
type ReceiveCommand struct {
	Args struct {
		Account string `positional-arg-name:"account"`
	} `positional-args:"yes"`
}

// Execute runs the `receive` command.
func (receiveCommand ReceiveCommand) Execute(args []string) error {
	var accountNumber uint32
	// if no account name was passed in
	if receiveCommand.Args.Account == "" {
		// display menu options to select account
		var err error
		accountNumber, err = utils.SelectAccount()
		if err != nil {
			return err
		}
	} else {
		// if an account name was passed in e.g. ./dcrcli receive default
		// get the address corresponding to the account name and use it
		var err error
		accountNumber, err = utils.Wallet.AccountNumber(receiveCommand.Args.Account)
		if err != nil {
			return fmt.Errorf("Error fetching account number: %s", err.Error())
		}
	}

	receiveAddress, err := utils.Wallet.GenerateReceiveAddress(accountNumber)
	if err != nil {
		return err
	}

	qr, err := qrcode.New(receiveAddress, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("Error generating QR Code: %s", err.Error())
	}

	res := &utils.Response{
		Columns: []string{
			"Address",
			"QR Code",
		},
		Result: [][]interface{}{
			[]interface{}{
				receiveAddress,
				qr.ToString(true),
			},
		},
	}
	utils.PrintResult(utils.StdoutTabWriter, res)
	return nil
}
