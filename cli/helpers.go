package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/raedahgroup/dcrcli/walletrpcclient"

	"github.com/decred/dcrd/dcrutil"
)

func getSendSourceAccount(c *walletrpcclient.Client, promptFn prompter) (uint32, error) {
	// get send  accounts
	accounts, err := c.Balance()
	if err != nil {
		return 0, err
	}
	validateInt := func(v string) error {
		if _, err := strconv.Atoi(v); err != nil {
			return err
		}
		return nil
	}
	promptItems := promptOption{
		Label:    "Select source account",
		Validate: validateInt,
	}
	accountItems := map[int]uint32{}
	for idx, v := range accounts {
		promptItems.Options = append(promptItems.Options, fmt.Sprintf("%3d. %s (%s)",
			idx, v.AccountName, dcrutil.Amount(v.Total).String()))
		accountItems[idx] = v.AccountNumber
	}

	result, err := promptFn(promptItems)
	if err != nil {
		return 0, fmt.Errorf("error getting selected account: %s", err.Error())
	}

	choice, err := strconv.Atoi(result)
	if err != nil {
		return 0, fmt.Errorf("error getting selected account: %s", err.Error())
	}

	account, ok := accountItems[choice]
	if !ok {
		return 0, fmt.Errorf("error selecting account")
	}

	return account, nil
}

func getSendDestinationAddress(c *walletrpcclient.Client, promptFn prompter) (string, error) {
	validate := func(address string) error {
		isValid, err := c.ValidateAddress(address)
		if err != nil {
			return fmt.Errorf("error validating address: %s", err.Error())
		}

		if !isValid {
			return errors.New("invalid address")
		}
		return nil
	}

	prompt := promptOption{
		Label:    "Destination Address",
		Validate: validate,
	}

	result, err := promptFn(prompt)
	if err != nil {
		return "", fmt.Errorf("error receiving input: %s", err.Error())
	}

	return result, nil
}

func getSendAmount(promptFn prompter) (int64, error) {
	var amount int64
	var err error

	validate := func(value string) error {
		amount, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing amount: %s", err.Error())
		}
		return nil
	}

	prompt := promptOption{
		Label:    "Amount (DCR)",
		Validate: validate,
	}

	_, err = promptFn(prompt)
	if err != nil {
		return 0, fmt.Errorf("error receiving input: %s", err.Error())
	}

	return amount, nil
}

func getWalletPassphrase(promptFn prompter) (string, error) {
	prompt := promptOption{
		Label:  "Wallet Passphrase",
		Secure: true,
	}

	result, err := promptFn(prompt)
	if err != nil {
		return "", fmt.Errorf("error receiving input: %s", err.Error())
	}
	return result, nil
}
