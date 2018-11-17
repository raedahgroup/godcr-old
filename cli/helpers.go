package cli

import (
	"errors"
	"fmt"
	"strconv"

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

func getSendAmount() (int64, error) {
	var amount int64
	var err error

	validateAmount := func(input string) error {
		amount, err = strconv.ParseInt(input, 10, 64)
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
