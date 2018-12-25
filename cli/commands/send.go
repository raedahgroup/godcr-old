package commands

import (
	"fmt"
	"github.com/raedahgroup/dcrcli/cli/commands/utils"
)

// SendCommand lets the user send DCR.
type SendCommand struct{}

// Execute runs the `send` command.
func (s SendCommand) Execute(args []string) error {
	res, err := send(false)
	if err != nil {
		return err
	}
	utils.PrintResult(utils.StdoutTabWriter, res)
	return nil
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct{}

// Execute runs the `send-custom` command.
func (s SendCustomCommand) Execute(args []string) error {
	res, err := send(true)
	if err != nil {
		return err
	}
	utils.PrintResult(utils.StdoutTabWriter, res)
	return nil
}

func send(custom bool) (*utils.Response, error) {
	var err error

	sourceAccount, err := utils.SelectAccount()
	if err != nil {
		return nil, err
	}

	// check if account has positive non-zero balance before proceeding
	// if balance is zero, there'd be no unspent outputs to use
	accountBalance, err := utils.Wallet.AccountBalance(sourceAccount)
	if err != nil {
		return nil, err
	}
	if accountBalance.Total == 0 {
		return nil, fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

	destinationAddress, err := utils.GetSendDestinationAddress()
	if err != nil {
		return nil, err
	}

	sendAmount, err := utils.GetSendAmount()
	if err != nil {
		return nil, err
	}

	var utxoSelection []string
	if custom {
		// get all utxos in account, pass 0 amount to get all
		utxos, err := utils.Wallet.UnspentOutputs(sourceAccount, 0)
		if err != nil {
			return nil, err
		}

		utxoSelection, err = utils.GetUtxosForNewTransaction(utxos, sendAmount)
		if err != nil {
			return nil, err
		}
	}

	passphrase, err := utils.GetWalletPassphrase()
	if err != nil {
		return nil, err
	}

	var sentTransactionHash string
	if custom {
		sentTransactionHash, err = utils.Wallet.SendFromUTXOs(utxoSelection, sendAmount, sourceAccount, destinationAddress, passphrase)
	} else {
		sentTransactionHash, err = utils.Wallet.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
	}

	if err != nil {
		return nil, err
	}

	res := &utils.Response{
		Columns: []string{
			"Result",
			"Hash",
		},
		Result: [][]interface{}{
			[]interface{}{
				"The transaction was published successfully",
				sentTransactionHash,
			},
		},
	}

	return res, nil
}
