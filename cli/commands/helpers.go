package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/decred/dcrd/dcrutil"
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
func getSendDestinationAddress(wallet walletcore.Wallet) (string, error) {
	validateAddressInput := func(address string) error {
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

	address, err := terminalprompt.RequestInput("Destination Address", validateAddressInput)
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

// getUtxosForNewTransaction fetches unspent transaction outputs to be used in a transaction.
func getUtxosForNewTransaction(utxos []*walletcore.UnspentOutput, sendAmount float64) ([]string, error) {
	var selectedUtxos []string
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
			selectedUtxos = append(selectedUtxos, utxo.OutputKey)
		}

		if totalAmountSelected < sendAmount {
			return errors.New("Invalid selection. Total amount from selected outputs is smaller than amount to send")
		}

		return nil
	}

	options := make([]string, len(utxos))
	for index, utxo := range utxos {
		options[index] = fmt.Sprintf("%s (%s)", utxo.OutputKey, utxo.Amount.String())
	}

	_, err = terminalprompt.RequestSelection("Select unspent outputs (e.g 1-4,6)", options, validateUtxoSelection)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return nil, fmt.Errorf("error reading selection: %s", err.Error())
	}
	return selectedUtxos, nil
}
