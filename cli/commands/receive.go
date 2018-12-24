package commands

import (
	"fmt"

	"github.com/raedahgroup/dcrcli/cli"

	qrcode "github.com/skip2/go-qrcode"
)

// ReceiveCommand generates and address for a user to receive DCR.
type ReceiveCommand struct {
	Args struct {
		Account string `positional-arg-name:"account"`
	} `positional-args:"yes"`
}

// Execute runs the `receive` command.
func (r ReceiveCommand) Execute(_ []string) error {
	var accountNumber uint32
	walletrpcclient := cli.WalletClient

	// if no account name was passed in
	if r.Args.Account == "" {
		// display menu options to select account
		var err error
		accountNumber, err = cli.GetSendSourceAccount(walletrpcclient)
		if err != nil {
			return err
		}
	} else {
		// if an account name was passed in e.g. ./dcrcli receive default
		// get the address corresponding to the account name and use it
		var err error
		accountNumber, err = walletrpcclient.AccountNumber(r.Args.Account)
		if err != nil {
			return fmt.Errorf("Error fetching account number: %s", err.Error())
		}
	}

	receiveResult, err := walletrpcclient.Receive(accountNumber)
	if err != nil {
		return err
	}

	qr, err := qrcode.New(receiveResult.Address, qrcode.Medium)
	if err != nil {
		return fmt.Errorf("Error generating QR Code: %s", err.Error())
	}

	res := &cli.Response{
		Columns: []string{
			"Address",
			"QR Code",
		},
		Result: [][]interface{}{
			[]interface{}{
				receiveResult.Address,
				qr.ToString(true),
			},
		},
	}
	cli.PrintResult(cli.StdoutWriter, res)
	return nil
}
