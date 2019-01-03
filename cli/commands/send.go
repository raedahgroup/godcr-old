package commands

import (
	"fmt"
	"sort"

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
	return send(wallet, false, false)
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct {
	commanderStub
	ManuallySelectInputs bool `short:"l" long:"list" description:"Display a list of unspent output to select for the transaction"`
}

// Run runs the `send-custom` command.
func (s SendCustomCommand) Run(wallet walletcore.Wallet) error {
	return send(wallet, true, s.ManuallySelectInputs)
}

func send(wallet walletcore.Wallet, custom bool, manuallySelectInputs bool) (err error) {
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

	var destinationAddresses []string
	sendAmounts := make(map[string]float64)
	var sendAmountTotal float64

	for {
		destinationAddress, err := getSendDestinationAddress(wallet, len(destinationAddresses))
		if err != nil {
			return err
		}
		if destinationAddress == "" {
			break
		}

		if _, addressExists := sendAmounts[destinationAddress]; addressExists {
			promptMessage := fmt.Sprintf("The address %s has already been added. Do you want to change the amount", destinationAddress)
			changeAmountConfirmed, err := terminalprompt.RequestYesNoConfirmation(promptMessage, "Y")
			if err != nil {
				return err
			}
			if !changeAmountConfirmed {
				continue
			}
		} else {
			destinationAddresses = append(destinationAddresses, destinationAddress)
		}

		sendAmount, err := getSendAmount()
		if err != nil {
			return err
		}
		sendAmounts[destinationAddress] = sendAmount
		sendAmountTotal += sendAmount
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

		bestSizedInput := func() (result []*walletcore.UnspentOutput) {
			sort.Slice(utxos, func(i, j int) bool {
				return utxos[i].Amount < utxos[j].Amount
			})
			var accumulatedAmount float64
			for _, utxo := range utxos {
				if utxo.Amount.ToCoin() > sendAmountTotal {
					result = []*walletcore.UnspentOutput{utxo}
					return
				}
				if accumulatedAmount <= sendAmountTotal {
					result = append(result, utxo)
					accumulatedAmount += utxo.Amount.ToCoin()
				}
			}

			return
		}

		utxoSelection = bestSizedInput()
		if manuallySelectInputs {
			utxoSelection, err = getUtxosForNewTransaction(wallet, utxos, sendAmountTotal, utxoSelection)
			if err != nil {
				return err
			}
		}
	}

	passphrase, err := getWalletPassphrase()
	if err != nil {
		return err
	}

	if custom {
		fmt.Println("You are about to spend the input")
		for _, output := range utxoSelection {
			address, err := getAddressFromUnspentOutputsResult(output)
			if err != nil {
				fmt.Println(fmt.Sprintf("Cannot extract address from output: %v", err))
			}
			fmt.Println(fmt.Sprintf(" %s from %s (%s)", output.Amount.String(), address, output.Amount.String()))
		}
		fmt.Println("and send it to")
		for _, address := range destinationAddresses {
			fmt.Println(fmt.Sprintf(" %v DCR to %s", sendAmounts[address], address))
		}
	} else {
		if len(destinationAddresses) == 1 {
			fmt.Println(fmt.Sprintf("You are about to send %f DCR to %s", sendAmountTotal, destinationAddresses[0]))
		} else {
			fmt.Println("You are about to send")
			for _, address := range destinationAddresses {
				fmt.Println(fmt.Sprintf(" %v DCR to %s", sendAmounts[address], address))
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

	var sendDestinations []txhelper.TransactionDestination
	for _, address := range destinationAddresses {
		sendDestinations = append(sendDestinations, txhelper.TransactionDestination{Amount: sendAmounts[address], Address: address})
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
