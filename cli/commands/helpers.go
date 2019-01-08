package commands

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrlibwallet/txhelper"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/cli/termio/terminalprompt"
)

// selectAccount lists accounts in wallet and prompts user to select an account, then returns the account number for that account.
// If there is only one account available, it returns the account number for that account.
func selectAccount(wallet walletcore.Wallet) (uint32, error) {
	var selection int
	var err error

	// get send  accounts
	accounts, err := wallet.AccountsOverview()
	if err != nil {
		return 0, err
	}

	// Proceed with default account if there's no other account.
	if len(accounts) == 1 {
		return accounts[0].Number, nil
	}

	// validateAccountSelection ensures that the input received is a number that corresponds to an account
	validateAccountSelection := func(input string) error {
		minAllowed, maxAllowed := 1, len(accounts)
		errWrongInput := fmt.Errorf("Error: input must be between %d and %d", minAllowed, maxAllowed)
		if selection, err = strconv.Atoi(input); err != nil {
			return errWrongInput
		}
		if selection < minAllowed || selection > maxAllowed {
			return errWrongInput
		}
		selection--
		return nil
	}

	options := make([]string, len(accounts))
	for index, account := range accounts {
		options[index] = fmt.Sprintf("%s (%s)", account.Name, account.Balance.Total.String())
	}

	_, err = terminalprompt.RequestSelection("Select account", options, validateAccountSelection)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return 0, fmt.Errorf("error getting selected account: %s", err.Error())
	}

	return accounts[selection].Number, nil
}

// getSendTxDestinations fetches the destinations info to send DCRs to from the user.
func getSendTxDestinations(wallet walletcore.Wallet) (destinations []txhelper.TransactionDestination, sendAmountTotal float64, err error) {
	var index int
	validateAddressInput := func(address string) error {
		if address == "" && index > 0 {
			return nil
		}
		if address == "" {
			return errors.New("You did not specify an address. Try again.")
		}

		isValid, err := wallet.ValidateAddress(address)
		if err != nil {
			return fmt.Errorf("error validating address: %s", err.Error())
		}

		if !isValid {
			return errors.New("That is not a valid address. Try again.")
		}
		return nil
	}

	sendAmountAddressMap := make(map[string]float64)

	for {
		label := "Destination Address"
		if index > 0 {
			label = fmt.Sprintf("Destination Address %d (or blank to continue)", index+1)
		}

		destinationAddress, err := terminalprompt.RequestInput(label, validateAddressInput)
		if err != nil {
			return nil, 0, fmt.Errorf("error receiving input: %s", err.Error())
		}

		if destinationAddress == "" {
			break
		}

		if _, addressExists := sendAmountAddressMap[destinationAddress]; addressExists {
			promptMessage := fmt.Sprintf("The address %s has already been added. Do you want to change the amount?", destinationAddress)
			changeAmountConfirmed, err := terminalprompt.RequestYesNoConfirmation(promptMessage, "N")
			if err != nil {
				return nil, 0, fmt.Errorf("error receiving input: %s", err.Error())
			}
			if !changeAmountConfirmed {
				continue
			}
			index--
		}

		sendAmount, err := getSendAmount()
		if err != nil {
			return nil, 0, fmt.Errorf("error receiving input: %s", err.Error())
		}
		sendAmountAddressMap[destinationAddress] = sendAmount
		index++
	}

	for address, amount := range sendAmountAddressMap {
		destinations = append(destinations, txhelper.TransactionDestination{Address: address, Amount: amount})
		sendAmountTotal += amount
	}
	return
}

// getSendAmount fetches the amout of DCRs to send from the user.
func getSendAmount() (float64, error) {
	var amount float64
	var err error

	validateAmount := func(input string) error {
		if input == "" {
			return errors.New("You did not specify an amount. Try again.")
		}

		amount, err = strconv.ParseFloat(input, 64)
		if err != nil {
			return fmt.Errorf("Invalid amount. Try again")
		}
		return nil
	}

	_, err = terminalprompt.RequestInput("Amount (DCR)", validateAmount)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return 0, fmt.Errorf("error receiving input: %s", err.Error())
	}

	return amount, nil
}

// getChangeOutputDestinations fetches the amount to be sent to each change address
func getChangeOutputDestinations(wallet walletcore.Wallet, totalInputAmount float64, sourceAccount uint32,
	nUtxoSelection int, sendDestinations []txhelper.TransactionDestination) ([]txhelper.TransactionDestination, error) {

	useRandomChangeAmounts, err := terminalprompt.RequestYesNoConfirmation("Use random amounts for the change outputs?", "y")
	if err != nil {
		return nil, fmt.Errorf("error reading your response: %s", err.Error())
	}

	amountInAtom, err := txhelper.AmountToAtom(totalInputAmount)
	if err != nil {
		return nil, err
	}

	if useRandomChangeAmounts {
		return getChangeDestinationsWithRandomAmounts(wallet, amountInAtom, sourceAccount,
			nUtxoSelection, sendDestinations)
	} else {
		return getChangeDestinationsFromUser(wallet, amountInAtom, sourceAccount,
			nUtxoSelection, sendDestinations)
	}
}

// getChangeDestinationsWithRandomAmounts generates change destination(s) based on the number of change address the user want
func getChangeDestinationsWithRandomAmounts(wallet walletcore.Wallet, amountInAtom int64, sourceAccount uint32,
	nUtxoSelection int, sendDestinations []txhelper.TransactionDestination) (changeOutputDestinations []txhelper.TransactionDestination, err error) {

	nChangeOutputs, err := terminalprompt.RequestNumberInput("How many change outputs would you like to use?", 1)
	if err != nil {
		return
	}

	var changeAddresses []string
	for i := 0; i < nChangeOutputs; i++ {
		address, err := wallet.GenerateReceiveAddress(sourceAccount)
		if err != nil {
			return nil, fmt.Errorf("error generating address: %s", err.Error())
		}
		changeAddresses = append(changeAddresses, address)
	}

	changeAmount, err := txhelper.EstimateChange(nUtxoSelection, amountInAtom, sendDestinations, changeAddresses)
	if err != nil {
		return nil, fmt.Errorf("error in getting change amount: %s", err.Error())
	}
	if changeAmount <= 0 {
		return
	}

	var portionRations []float64
	var rationSum float64
	for i := 0; i < nChangeOutputs; i++ {
		portion := rand.Float64()
		portionRations = append(portionRations, portion)
		rationSum += portion
	}

	for i, portion := range portionRations {
		portionPercentage := portion / rationSum
		amount := portionPercentage * float64(changeAmount)

		changeOutput := txhelper.TransactionDestination{
			Address: changeAddresses[i],
			Amount:  dcrutil.Amount(amount).ToCoin(),
		}
		changeOutputDestinations = append(changeOutputDestinations, changeOutput)
	}
	return
}

// getChangeDestinationsFromUser fetches change destination from the user progressively until the total available change amount is covered
func getChangeDestinationsFromUser(wallet walletcore.Wallet, amountInAtom int64, sourceAccount uint32, nUtxoSelection int, sendDestinations []txhelper.TransactionDestination) ([]txhelper.TransactionDestination, error) {
	var changeOutputDestinations []txhelper.TransactionDestination
	var changeAddresses []string
	var amountAssigned int64

	var index int
	for {
		address, err := wallet.GenerateReceiveAddress(sourceAccount)
		if err != nil {
			return nil, fmt.Errorf("error in generating address: %s", err.Error())
		}
		changeAddresses = append(changeAddresses, address)
		totalChangeAmount, err := txhelper.EstimateChange(nUtxoSelection, amountInAtom, sendDestinations, changeAddresses)
		if err != nil {
			return nil, err
		}
		defaultAmount := dcrutil.Amount(totalChangeAmount - amountAssigned).ToCoin()

		prompt := fmt.Sprintf("[%d]: Change Amount (DCR) (default: %f)", index+1, defaultAmount)
		changeAmount, err := getChangeAmount(prompt)
		if err != nil {
			return nil, err
		}

		if changeAmount > defaultAmount {
			fmt.Println(fmt.Sprintf("Invalid amount. Enter an amount less than %f", defaultAmount))
			changeAddresses[len(changeAddresses)-1] = ""
			changeAddresses = changeAddresses[:len(changeAddresses)-1]
			continue
		}

		if changeAmount == 0 {
			changeAmount = defaultAmount
		}
		changeDestination := txhelper.TransactionDestination{Address: address, Amount: changeAmount}
		changeOutputDestinations = append(changeOutputDestinations, changeDestination)

		if changeAmount == defaultAmount {
			break
		}

		changeAmountInAtom, err := txhelper.AmountToAtom(changeAmount)
		if err != nil {
			return nil, err
		}
		amountAssigned += changeAmountInAtom
		index++
	}
	return changeOutputDestinations, nil
}

// getChangeAmount fetches the amount of DCRs to send from the user.
func getChangeAmount(prompt string) (amount float64, err error) {
	validateAmount := func(input string) error {
		if input == "" {
			amount = 0
			return nil
		}

		amount, err = strconv.ParseFloat(input, 64)
		if err != nil || amount <= 0 {
			return fmt.Errorf("Invalid amount. Try again")
		}
		return nil
	}

	_, err = terminalprompt.RequestInput(prompt, validateAmount)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return 0, fmt.Errorf("error receiving input: %s", err.Error())
	}

	return
}

// getWalletPassphrase fetches the user's wallet passphrase from the user.
func getWalletPassphrase() (string, error) {
	result, err := terminalprompt.RequestInputSecure("Spending Passphrase", terminalprompt.EmptyValidator)
	if err != nil {
		return "", fmt.Errorf("error receiving input: %s", err.Error())
	}
	return result, nil
}

// getUtxosForNewTransaction fetches unspent transaction outputs to be used in a transaction.
func getUtxosForNewTransaction(utxos []*walletcore.UnspentOutput, sendAmount float64) (selectedUtxos []*walletcore.UnspentOutput, totalAmountSelected float64, err error) {
	var removeWhiteSpace = func(str string) string {
		return strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, str)
	}

	// validateAccountSelection  ensures that the input received is a number that corresponds to an account
	validateUtxoSelection := func(selectedOptions string) error {
		minAllowed, maxAllowed := 1, len(utxos)
		errWrongInput := errors.New("your selection does not match any available option")

		// remove white space and split user input into comma-delimited selection ranges
		selectionRanges := strings.Split(removeWhiteSpace(selectedOptions), ",")
		var selection []int

		for _, minMaxRange := range selectionRanges {
			minMax := strings.Split(minMaxRange, "-")
			var min, max int
			var err error

			min, err = strconv.Atoi(minMax[0])
			if err != nil || min < minAllowed || min > maxAllowed {
				return errWrongInput
			}

			if len(minMax) == 1 {
				selection = append(selection, min-1)
				continue
			}

			max, err = strconv.Atoi(minMax[1])
			if err != nil || max < minAllowed || max > maxAllowed {
				return errWrongInput
			}

			// ensure min is actually smaller than max, swap if otherwise
			if min > max {
				min, max = max, min
			}

			for n := min; n <= max; n++ {
				selection = append(selection, n-1)
			}
		}

		if len(selection) == 0 {
			return errWrongInput
		}

		selectedUtxos = selectedUtxos[:0]
		totalAmountSelected = 0
		for _, n := range selection {
			utxo := utxos[n]
			totalAmountSelected += dcrutil.Amount(utxo.Amount).ToCoin()
			selectedUtxos = append(selectedUtxos, utxo)
		}

		if totalAmountSelected < sendAmount {
			return errors.New("Invalid selection. Total amount from selected inputs is smaller than amount to send")
		}

		return nil
	}

	options := make([]string, len(utxos))
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Amount < utxos[j].Amount
	})
	for index, utxo := range utxos {
		date := time.Unix(utxo.ReceiveTime, 0).Format("Mon Jan 2, 2006 3:04PM")
		options[index] = fmt.Sprintf("%s (%s) \t %s \t %d confirmation(s)", utxo.Address, utxo.Amount.String(), date, utxo.Confirmations)
	}

	_, err = terminalprompt.RequestSelection("Select input(s) (e.g 1-4,6)", options, validateUtxoSelection)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return nil, 0, fmt.Errorf("error reading selection: %s", err.Error())
	}
	return selectedUtxos, totalAmountSelected, nil
}

// bestSizedInput returns the smallest output or the least consecutive combination of
// outputs that can handle a transaction of the supplied sendAmountTotal from the utxos
func bestSizedInput(utxos []*walletcore.UnspentOutput, sendAmountTotal float64) ([]*walletcore.UnspentOutput, float64) {
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Amount < utxos[j].Amount
	})
	for _, utxo := range utxos {
		if utxo.Amount.ToCoin() > sendAmountTotal {
			return []*walletcore.UnspentOutput{utxo}, utxo.Amount.ToCoin()
		}
	}
	for noOfPairs := 2; noOfPairs <= len(utxos); noOfPairs++ {
		for i := 0; i < len(utxos); i++ {
			var accumulatedAmount float64
			var result []*walletcore.UnspentOutput
			for j := i; j < i+noOfPairs && j < len(utxos); j++ {
				result = append(result, utxos[j])
				accumulatedAmount += utxos[j].Amount.ToCoin()
				if accumulatedAmount >= sendAmountTotal {
					return result, accumulatedAmount
				}
			}
		}
	}
	var totalInputAmount float64
	for _, utxo := range utxos {
		totalInputAmount += utxo.Amount.ToCoin()
	}
	return utxos, sendAmountTotal
}
