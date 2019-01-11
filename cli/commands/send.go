package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// SendCommand lets the user send DCR.
type SendCommand struct {
	commanderStub
	Spendunconfirmed bool `short:"u" long:"spendunconfirmed" description:"allow users chose to use 0 confirmations for send transactions."`
}

// Run runs the `send` command.
func (s SendCommand) Run(ctx context.Context, wallet walletcore.Wallet) error {
	var requireConfirmation int32 = 0
	if s.Spendunconfirmed{
		send(wallet, requireConfirmation, false)	
	}else{
		send(wallet, walletcore.DefaultRequiredConfirmations, false)
	}
	return nil
}

// SendCustomCommand sends DCR using coin control.
type SendCustomCommand struct {
	commanderStub
	Spendunconfirmed bool `short:"u" long:"spendunconfirmed" description:"Use unconfirmed outputs for send transactions."`
}

// Run runs the `send-custom` command.
func (s SendCustomCommand) Run(ctx context.Context, wallet walletcore.Wallet) error {
	var requireConfirmation int32 
	if s.Spendunconfirmed{
		requireConfirmation = 0
	}else{
		requireConfirmation = walletcore.DefaultRequiredConfirmations
	}

	return send(wallet, requireConfirmation, false)
}

func send(wallet walletcore.Wallet, requireConfirmation int32, custom bool) error {
	sourceAccount, err := selectAccount(wallet)
	if err != nil {
		return err
	}


	// check if account has positive non-zero balance before proceeding
	// if balance is zero, there'd be no unspent outputs to use

	accountBalance, err := wallet.AccountBalance(sourceAccount, requireConfirmation)
	if err != nil {
		return err
	}
			fmt.Println(accountBalance)

	if accountBalance.Total == 0 {
		return fmt.Errorf("Selected account has 0 balance. Cannot proceed")
	}

	sendDestinations, sendAmountTotal, err := getSendTxDestinations(wallet)
	if err != nil {
		return err
	}

	if accountBalance.Spendable.ToCoin() < sendAmountTotal {
		return fmt.Errorf("Selected account has insufficient balance. Cannot proceed")
	}

	var sentTxHash string
	if custom {
		sentTxHash, err = completeCustomSend(wallet, sourceAccount, sendDestinations, sendAmountTotal)
	} else {
		sentTxHash, err = completeNormalSend(wallet, sourceAccount, sendDestinations)
	}

	if err != nil {
		return err
	}

	fmt.Println("Sent txid", sentTxHash)
	return nil
}

func completeCustomSend(wallet walletcore.Wallet, sourceAccount uint32, sendDestinations []txhelper.TransactionDestination, sendAmountTotal float64) (string, error) {
	var changeOutputDestinations []txhelper.TransactionDestination
	var utxoSelection []*walletcore.UnspentOutput
	var totalInputAmount float64

	// get all utxos in account, pass 0 amount to get all
	utxos, err := wallet.UnspentOutputs(sourceAccount, 0, walletcore.DefaultRequiredConfirmations)
	if err != nil {
		return "", err
	}

	choice, err := terminalprompt.RequestInput("Would you like to (a)utomatically or (m)anually select inputs? (A/m)", func(input string) error {
		switch strings.ToLower(input) {
		case "", "a", "m":
			return nil
		}
		return errors.New("invalid entry")
	})
	if err != nil {
		return "", fmt.Errorf("error in reading choice: %s", err.Error())
	}
	if strings.ToLower(choice) == "a" || choice == "" {
		utxoSelection, totalInputAmount = bestSizedInput(utxos, sendAmountTotal)
	} else {
		utxoSelection, totalInputAmount, err = getUtxosForNewTransaction(utxos, sendAmountTotal)
		if err != nil {
			return "", err
		}
	}

	changeOutputDestinations, err = getChangeOutputDestinations(wallet, totalInputAmount, sourceAccount,
		len(utxoSelection), sendDestinations)
	if err != nil {
		return "", err
	}

	passphrase, err := getWalletPassphrase()
	if err != nil {
		return "", err
	}

	fmt.Println("You are about to spend the input(s)")
	for _, utxo := range utxoSelection {
		fmt.Println(fmt.Sprintf(" %s \t from %s", utxo.Amount.String(), utxo.Address))
	}
	fmt.Println("and send")
	for _, destination := range sendDestinations {
		fmt.Println(fmt.Sprintf(" %f DCR \t to %s", destination.Amount, destination.Address))
	}
	for _, destination := range changeOutputDestinations {
		fmt.Println(fmt.Sprintf(" %f DCR \t to %s (change)", destination.Amount, destination.Address))
	}

	sendConfirmed, err := terminalprompt.RequestYesNoConfirmation("Do you want to broadcast it?", "")
	if err != nil {
		return "", fmt.Errorf("error reading your response: %s", err.Error())
	}

	if !sendConfirmed {
		return "", errors.New("transaction canceled")
	}

	var outputKeys []string
	for _, utxo := range utxoSelection {
		outputKeys = append(outputKeys, utxo.OutputKey)
	}
	return wallet.SendFromUTXOs(sourceAccount, walletcore.DefaultRequiredConfirmations, outputKeys, sendDestinations, changeOutputDestinations, passphrase)
}

func completeNormalSend(wallet walletcore.Wallet, sourceAccount uint32, sendDestinations []txhelper.TransactionDestination) (string, error) {
	passphrase, err := getWalletPassphrase()
	if err != nil {
		return "", err
	}

	if len(sendDestinations) == 1 {
		fmt.Println(fmt.Sprintf("You are about to send %f DCR to %s", sendDestinations[0].Amount, sendDestinations[0].Address))
	} else {
		fmt.Println("You are about to send")
		for _, destination := range sendDestinations {
			fmt.Println(fmt.Sprintf(" %f DCR \t to %s", destination.Amount, destination.Address))
		}
	}

	sendConfirmed, err := terminalprompt.RequestYesNoConfirmation("Do you want to broadcast it?", "")
	if err != nil {
		return "", fmt.Errorf("error reading your response: %s", err.Error())
	}

	if !sendConfirmed {
		return "", errors.New("transaction cancelled")
	}

	return wallet.SendFromAccount(sourceAccount, walletcore.DefaultRequiredConfirmations, sendDestinations, passphrase)
}
