package commands

import (
	"fmt"
<<<<<<< HEAD

	"github.com/raedahgroup/godcr/cli/termio"
	ws "github.com/raedahgroup/godcr/walletsource"
=======
	"github.com/raedahgroup/dcrcli/cli/utils"
>>>>>>> little refactor
)

// SendCommand lets the user send DCR.
type SendCommand struct {
	CommanderStub
}

<<<<<<< HEAD
// Run runs the `send` command.
func (s SendCommand) Run(walletsource ws.WalletSource, args []string) error {
	return send(walletsource, false)
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct {
	CommanderStub
}

// Run runs the `send-custom` command.
func (s SendCustomCommand) Run(walletsource ws.WalletSource, args []string) error {
	return send(walletsource, true)
}

func send(walletsource ws.WalletSource, custom bool) (err error) {
	sourceAccount, err := selectAccount(walletsource)
=======
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
>>>>>>> little refactor
	if err != nil {
		return err
	}

	// check if account has positive non-zero balance before proceeding
	// if balance is zero, there'd be no unspent outputs to use
<<<<<<< HEAD
	accountBalance, err := walletsource.AccountBalance(sourceAccount)
=======
	accountBalance, err := utils.Wallet.AccountBalance(sourceAccount)
>>>>>>> little refactor
	if err != nil {
		return err
	}
	if accountBalance.Total == 0 {
		return fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

<<<<<<< HEAD
	destinationAddress, err := getSendDestinationAddress(walletsource)
=======
	destinationAddress, err := utils.GetSendDestinationAddress()
>>>>>>> little refactor
	if err != nil {
		return err
	}

<<<<<<< HEAD
	sendAmount, err := getSendAmount()
=======
	sendAmount, err := utils.GetSendAmount()
>>>>>>> little refactor
	if err != nil {
		return err
	}

	var utxoSelection []string
	if custom {
		// get all utxos in account, pass 0 amount to get all
<<<<<<< HEAD
		utxos, err := walletsource.UnspentOutputs(sourceAccount, 0)
=======
		utxos, err := utils.Wallet.UnspentOutputs(sourceAccount, 0)
>>>>>>> little refactor
		if err != nil {
			return err
		}

<<<<<<< HEAD
		utxoSelection, err = getUtxosForNewTransaction(utxos, sendAmount)
=======
		utxoSelection, err = utils.GetUtxosForNewTransaction(utxos, sendAmount)
>>>>>>> little refactor
		if err != nil {
			return err
		}
	}

<<<<<<< HEAD
	passphrase, err := getWalletPassphrase()
=======
	passphrase, err := utils.GetWalletPassphrase()
>>>>>>> little refactor
	if err != nil {
		return err
	}

	var sentTransactionHash string
	if custom {
<<<<<<< HEAD
		sentTransactionHash, err = walletsource.SendFromUTXOs(utxoSelection, sendAmount, sourceAccount, destinationAddress, passphrase)
	} else {
		sentTransactionHash, err = walletsource.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
=======
		sentTransactionHash, err = utils.Wallet.SendFromUTXOs(utxoSelection, sendAmount, sourceAccount, destinationAddress, passphrase)
	} else {
		sentTransactionHash, err = utils.Wallet.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
>>>>>>> little refactor
	}

	if err != nil {
		return err
	}

<<<<<<< HEAD
	columns := []string{
		"Result",
		"Hash",
	}
	rows := [][]interface{}{
		[]interface{}{
			"The transaction was published successfully",
			sentTransactionHash,
=======
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
>>>>>>> little refactor
		},
	}

	termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
	return nil
}
