package commands

import (
	"fmt"
	"github.com/raedahgroup/dcrcli/cli/utils"

<<<<<<< HEAD
	"github.com/raedahgroup/godcr/cli/termio"
	ws "github.com/raedahgroup/godcr/walletsource"
=======
>>>>>>> little refactor
	qrcode "github.com/skip2/go-qrcode"
)

// ReceiveCommand generates and address for a user to receive DCR.
type ReceiveCommand struct {
	CommanderStub
	Args struct {
		Account string `positional-arg-name:"account"`
	} `positional-args:"yes"`
}

<<<<<<< HEAD
// Run runs the `receive` command.
func (r ReceiveCommand) Run(walletsource ws.WalletSource, args []string) error {
=======
// Execute runs the `receive` command.
func (receiveCommand ReceiveCommand) Execute(args []string) error {
>>>>>>> little refactor
	var accountNumber uint32
	// if no account name was passed in
	if receiveCommand.Args.Account == "" {
		// display menu options to select account
		var err error
<<<<<<< HEAD
		accountNumber, err = selectAccount(walletsource)
=======
		accountNumber, err = utils.SelectAccount()
>>>>>>> little refactor
		if err != nil {
			return err
		}
	} else {
		// if an account name was passed in e.g. ./godcr receive default
		// get the address corresponding to the account name and use it
		var err error
<<<<<<< HEAD
		accountNumber, err = walletsource.AccountNumber(r.Args.Account)
=======
		accountNumber, err = utils.Wallet.AccountNumber(receiveCommand.Args.Account)
>>>>>>> little refactor
		if err != nil {
			return fmt.Errorf("Error fetching account number: %s", err.Error())
		}
	}

<<<<<<< HEAD
	receiveAddress, err := walletsource.GenerateReceiveAddress(accountNumber)
=======
	receiveAddress, err := utils.Wallet.GenerateReceiveAddress(accountNumber)
>>>>>>> little refactor
	if err != nil {
		return err
	}

	qr, err := qrcode.New(receiveAddress, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("Error generating QR Code: %s", err.Error())
	}

<<<<<<< HEAD
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
=======
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
>>>>>>> little refactor
	return nil
}
