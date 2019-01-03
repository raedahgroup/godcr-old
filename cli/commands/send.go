package commands

import (
	"errors"
	"fmt"
	"strings"

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

func send(wallet walletcore.Wallet, custom bool) (err error) {
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
		return fmt.Errorf("Selected account has low balance. Cannot proceed")
	}

	var utxoSelection []*walletcore.UnspentOutput

	if custom {
		// get all utxos in account, pass 0 amount to get all
		utxos, err := wallet.UnspentOutputs(sourceAccount, 0)
		if err != nil {
			return err
		}

		choice, err := terminalprompt.RequestInput("Would you like to (a)utomatically or (m)anually select inputs? (A/m)", func(input string) error {
			switch strings.ToLower(input) {
			case "":
				return nil
			case "a":
				return nil
			case "m":
				return nil
			}
			return errors.New("invalid entry")
		})
		if err != nil {
			return fmt.Errorf("error in reading choice: %s", err.Error())
		}
		if strings.ToLower(choice) == "a" || choice == "" {
			utxoSelection = bestSizedInput(utxos, sendAmountTotal)
		} else {
			utxoSelection, err = getUtxosForNewTransaction(wallet, utxos, sendAmountTotal)
		}
	}

	passphrase, err := getWalletPassphrase()
	if err != nil {
		return err
	}

	if custom {
		fmt.Println("You are about to spend the input")
		for _, output := range utxoSelection {
			fmt.Println(fmt.Sprintf(" %s from %s", output.Amount.String(), output.Address))
		}
		fmt.Println("and send it to")
		for _, destinatoin := range sendDestinations {
			fmt.Println(fmt.Sprintf(" %v DCR to %s", destinatoin.Amount, destinatoin.Address))
		}
	} else {
		if len(sendDestinations) == 1 {
			fmt.Println(fmt.Sprintf("You are about to send %f DCR to %s", sendDestinations[0].Amount, sendDestinations[0].Address))
		} else {
			fmt.Println("You are about to send")
			for _, destination := range sendDestinations {
				fmt.Println(fmt.Sprintf(" %v DCR to %s", destination.Amount, destination.Address))
			}
		}
	}

	sendConfirmed, err := terminalprompt.RequestYesNoConfirmation("Are you sure?", "")
	if err != nil {
		return fmt.Errorf("error reading your response: %s", err.Error())
	}

	if !sendConfirmed {
		fmt.Println("Canceled")
		return nil
	}

	var sentTransactionHash string
	if custom {
		var outputs []string
		for _, utox := range utxoSelection {
			outputs = append(outputs, utox.OutputKey)
		}
		sentTransactionHash, err = wallet.SendFromUTXOs(sourceAccount, outputs, sendDestinations, passphrase)
	} else {
		sentTransactionHash, err = wallet.SendFromAccount(sourceAccount, sendDestinations, passphrase)
	}

	if err != nil {
		return err
	}

	fmt.Println("Sent. Txid", sentTransactionHash)
	return nil
}
