package commands

import (
	"errors"
	"fmt"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"strings"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// SendCommand lets the user send DCR.
type SendCommand struct {
	commanderStub
}

// Run runs the `send` command.
func (s SendCommand) Run(wallet walletcore.Wallet) error {
	return send(wallet, false)
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct {
	commanderStub
}

// Run runs the `send-custom` command.
func (s SendCustomCommand) Run(wallet walletcore.Wallet) error {
	return send(wallet, true)
}

func send(wallet walletcore.Wallet, custom bool) error {
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

	sendDestinations, err := getSendTxDestinations(wallet)
	if err != nil {
		return err
	}
	var sendAmountTotal float64
	for _, destination := range sendDestinations {
		sendAmountTotal += destination.Amount
	}

	if accountBalance.Spendable.ToCoin() < sendAmountTotal {
		return fmt.Errorf("Selected account has insufficient balance. Cannot proceed")
	}

	var changeOutputDestinations []txhelper.TransactionDestination
	var utxoSelection []*walletcore.UnspentOutput
	var totalInputAmount float64

	if custom {
		// get all utxos in account, pass 0 amount to get all
		utxos, err := wallet.UnspentOutputs(sourceAccount, 0)
		if err != nil {
			return err
		}

		choice, err := terminalprompt.RequestInput("Would you like to (a)utomatically or (m)anually select inputs? (A/m)", func(input string) error {
			switch strings.ToLower(input) {
			case "", "a", "m":
				return nil
			}
			return errors.New("invalid entry")
		})
		if err != nil {
			return fmt.Errorf("error in reading choice: %s", err.Error())
		}
		if strings.ToLower(choice) == "a" || choice == "" {
			utxoSelection, totalInputAmount = bestSizedInput(utxos, sendAmountTotal)
		} else {
			utxoSelection, totalInputAmount, err = getUtxosForNewTransaction(utxos, sendAmountTotal)
		}

		changeOutputDestinations, err = getChangeOutputDestinations(wallet, totalInputAmount, sourceAccount,
			len(utxoSelection), sendDestinations)
		if err != nil {
			return
		}
	}

	passphrase, err := getWalletPassphrase()
	if err != nil {
		return err
	}

	if custom {
		fmt.Println("You are about to spend the input(s)")
		for _, utxo := range utxoSelection {
			fmt.Println(fmt.Sprintf(" %s from %s", utxo.Amount.String(), utxo.Address))
		}
		fmt.Println("and send")
		for _, destination := range sendDestinations {
			fmt.Println(fmt.Sprintf(" %f DCR to %s", destination.Amount, destination.Address))
		}
		for _, destination := range changeOutputDestinations {
			fmt.Println(fmt.Sprintf(" %f DCR to %s(change)", destination.Amount, destination.Address))
		}
	} else {
		if len(sendDestinations) == 1 {
			fmt.Println(fmt.Sprintf("You are about to send %f DCR to %s", sendDestinations[0].Amount, sendDestinations[0].Address))
		} else {
			fmt.Println("You are about to send")
			for _, destination := range sendDestinations {
				fmt.Println(fmt.Sprintf(" %f DCR to %s", destination.Amount, destination.Address))
			}
		}
	}

	sendConfirmed, err := terminalprompt.RequestYesNoConfirmation("Do you want to broadcast it?", "")
	if err != nil {
		return fmt.Errorf("error reading your response: %s", err.Error())
	}

	if !sendConfirmed {
		fmt.Println("Canceled")
		return nil
	}

	sendDestinations = append(sendDestinations, changeOutputDestinations...)
	var sentTransactionHash string
	if custom {
		var utxos []string
		var totalInputAmount int64
		for _, utxo := range utxoSelection {
			utxos = append(utxos, utxo.OutputKey)
			totalInputAmount += int64(utxo.Amount)
		}

		changeAddress, err := wallet.GenerateReceiveAddress(sourceAccount)
		if err != nil {
			return fmt.Errorf("error generating change address: %s", err.Error())
		}

		changeAmount, err := txhelper.EstimateChange(len(utxos), int64(totalInputAmount), sendDestinations, []string{changeAddress})
		if err != nil {
			return fmt.Errorf("error estimating change amount: %s", err.Error())
		}

		changeDestinations := []txhelper.TransactionDestination{{
			Amount: dcrutil.Amount(changeAmount).ToCoin(),
			Address: changeAddress,
		}}

		sentTransactionHash, err = wallet.SendFromUTXOs(sourceAccount, utxos, sendDestinations, changeDestinations, passphrase)
	} else {
		sentTransactionHash, err = wallet.SendFromAccount(sourceAccount, sendDestinations, passphrase)
	}

	if err != nil {
		return err
	}

	fmt.Println("Sent txid", sentTransactionHash)
	return nil
}
