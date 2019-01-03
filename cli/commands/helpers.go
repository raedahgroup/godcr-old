package commands

import (
	"errors"
	"fmt"
	"github.com/raedahgroup/godcr/cli/termio"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/decred/dcrd/chaincfg"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/txscript"

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

// getSendDestinationAddress fetches the destination address to send DCRs to from the user.
func getSendDestinationAddress(wallet walletcore.Wallet, index int) (string, error) {
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

	label := "Destination Address"
	if index > 0 {
		label = fmt.Sprintf("Destination Address %d (or blank to continue)", index+1)
	}
	address, err := terminalprompt.RequestInput(label, validateAddressInput)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return "", fmt.Errorf("error receiving input: %s", err.Error())
	}

	return address, nil
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

// getWalletPassphrase fetches the user's wallet passphrase from the user.
func getWalletPassphrase() (string, error) {
	result, err := terminalprompt.RequestInputSecure("Wallet Passphrase", terminalprompt.EmptyValidator)
	if err != nil {
		return "", fmt.Errorf("error receiving input: %s", err.Error())
	}
	return result, nil
}

func getSendUtxoCount(maxCount int64) (int64, error) {
	var count int64
	var err error

	validateCount := func(input string) error {
		if input == "" {
			count = 1
			return nil
		}
		count, err = strconv.ParseInt(input, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing number: %s", err.Error())
		}
		if count > maxCount {
			return fmt.Errorf("you cannot select more than %d", maxCount)
		}
		return nil
	}

	_, err = terminalprompt.RequestInput("How many change outputs would you like to use? (default: 1)", validateCount)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return 0, fmt.Errorf("error receiving input: %s", err.Error())
	}

	return count, nil
}

func getUseRandomAmount() (bool, error) {
	var yes bool
	var err error

	validate := func(input string) error {
		if input == "" {
			input = "y"
		}
		switch strings.ToLower(input) {
		case "y":
			yes = true
			return nil
		case "n":
			return nil
		default:
			return errors.New("invalid entry")
		}
	}

	_, err = terminalprompt.RequestInput("Use random amounts for the change outputs? (Y/n)", validate)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return false, fmt.Errorf("error receiving input: %s", err.Error())
	}

	return yes, nil
}

// getUtxosForNewTransaction fetches unspent transaction outputs to be used in a transaction.
func getUtxosForNewTransaction(wallet walletcore.Wallet, utxos []*walletcore.UnspentOutput, sendAmount float64, defaultOptions []*walletcore.UnspentOutput) ([]*walletcore.UnspentOutput, error) {
	var selectedUtxos []*walletcore.UnspentOutput
	var err error

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

		var totalAmountSelected float64
		for _, n := range selection {
			utxo := utxos[n]
			totalAmountSelected += dcrutil.Amount(utxo.Amount).ToCoin()
			selectedUtxos = append(selectedUtxos, utxo)
		}

		if totalAmountSelected < sendAmount {
			return errors.New("Invalid selection. Total amount from selected outputs is smaller than amount to send")
		}

		return nil
	}

	if len(defaultOptions) > 0 {
		tabWriter := termio.StdoutWriter
		fmt.Fprintln(tabWriter, "Best sized input for the transaction is ")
		for index, utxo := range defaultOptions {
			address, err := getAddressFromUnspentOutputsResult(utxo)
			if err != nil {
				return nil, fmt.Errorf("error reading address: %s", err.Error())
			}
			date := time.Unix(utxo.ReceiveTime, 0).Format("Mon Jan 2, 2006 3:04PM")
			txn, err := wallet.GetTransaction(utxo.TransactionHash)
			if err != nil {
				return nil, fmt.Errorf("error reading transaction: %s", err.Error())
			}
			fmt.Fprintln(tabWriter, fmt.Sprintf(" [%v]: %s (%s) \t %s \t %v confirmation(s)", (index + 1), address, utxo.Amount.String(), date, txn.Confirmations))
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
			return nil, fmt.Errorf("error in reading choice: %s", err.Error())
		}
		if strings.ToLower(choice) == "a" || choice == "" {
			return defaultOptions, nil
		}
	}

	options := make([]string, len(utxos))
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Amount < utxos[j].Amount
	})
	for index, utxo := range utxos {
		address, err := getAddressFromUnspentOutputsResult(utxo)
		if err != nil {
			return nil, fmt.Errorf("error reading address: %s", err.Error())
		}
		date := time.Unix(utxo.ReceiveTime, 0).Format("Mon Jan 2, 2006 3:04PM")
		txn, err := wallet.GetTransaction(utxo.TransactionHash)
		if err != nil {
			return nil, fmt.Errorf("error reading transaction: %s", err.Error())
		}
		options[index] = fmt.Sprintf("%s (%s) \t %s \t %v confirmation(s)", address, utxo.Amount.String(), date, txn.Confirmations)
	}

	_, err = terminalprompt.RequestSelection("Select unspent outputs (e.g 1-4,6)", options, validateUtxoSelection)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return nil, fmt.Errorf("error reading selection: %s", err.Error())
	}
	return selectedUtxos, nil
}

func getAddressFromUnspentOutputsResult(utxo *walletcore.UnspentOutput) (address string, err error) {
	_, addresses, _, err := txscript.ExtractPkScriptAddrs(txscript.DefaultScriptVersion,
		utxo.PkScript, currentNet())
	if err != nil {
		return
	}
	if len(addresses) < 1 {
		return "", errors.New("Cannot extract any address from output")
	}
	return addresses[0].EncodeAddress(), nil
}

func currentNet() *chaincfg.Params {
	return &chaincfg.TestNet3Params //Todo check config to see if testnet is active
}
