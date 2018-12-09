package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrcli/cli/terminalprompt"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

func getSendSourceAccount(c *walletrpcclient.Client) (uint32, error) {
	var selection int
	var err error
	// get send  accounts
	accounts, err := c.Balance()
	if err != nil {
		return 0, err
	}
	// Proceed with default account if there's no other account.
	if len(accounts) == 1 {
		return accounts[0].AccountNumber, nil
	}
	// validateAccountSelection  ensures that the input received is a number that corresponds to an account
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
		options[index] = fmt.Sprintf("%s (%s)", account.AccountName, dcrutil.Amount(account.Total).String())
	}

	_, err = terminalprompt.RequestSelection("Select source account", options, validateAccountSelection)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return 0, fmt.Errorf("error getting selected account: %s", err.Error())
	}
	return accounts[selection].AccountNumber, nil
}

func getSendDestinationAddress(c *walletrpcclient.Client) (string, error) {
	validateAddressInput := func(address string) error {
		isValid, err := c.ValidateAddress(address)
		if err != nil {
			return fmt.Errorf("error validating address: %s", err.Error())
		}

		if !isValid {
			return errors.New("invalid address")
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

func getSendAmount() (float64, error) {
	var amount float64
	var err error

	validateAmount := func(input string) error {
		amount, err = strconv.ParseFloat(input, 64)
		if err != nil {
			return fmt.Errorf("error parsing amount: %s", err.Error())
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

func getWalletPassphrase() (string, error) {
	result, err := terminalprompt.RequestInputSecure("Wallet Passphrase", terminalprompt.EmptyValidator)
	if err != nil {
		return "", fmt.Errorf("error receiving input: %s", err.Error())
	}
	return result, nil
}

func getUtxosForNewTransaction(utxos []*walletrpcclient.UnspentOutputsResult) ([]string, error) {
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
		errWrongInput := fmt.Errorf("Error: invalid input. Selection should be numbers between %d - %d",
			minAllowed, maxAllowed)

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
				selection = append(selection, min)
				continue
			}

			max, err = strconv.Atoi(minMax[1])
			if err != nil || max < minAllowed || max > maxAllowed {
				return errWrongInput
			}

			// ensure min is actually smaller than max, swap if otherwise
			if min > max {
				temp := max
				max = min
				min = temp
			}

			for n := min; n <= max; n++ {
				selection = append(selection, n)
			}
		}

		if len(selection) == 0 {
			return errWrongInput
		}

		for _, n := range selection {
			selectedUtxos = append(selectedUtxos, utxos[n-1].OutputKey)
		}

		return nil
	}

	options := make([]string, len(utxos))
	for index, utxo := range utxos {
		options[index] = fmt.Sprintf("%s (%s)", utxo.OutputKey, utxo.Amount)
	}

	_, err = terminalprompt.RequestSelection("Select unspent outputs (e.g 1-4,6)", options, validateUtxoSelection)
	if err != nil {
		// There was an error reading input; we cannot proceed.
		return nil, fmt.Errorf("error reading selection: %s", err.Error())
	}
	return selectedUtxos, nil
}
