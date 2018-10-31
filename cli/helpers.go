package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"
	"github.com/raedahgroup/dcrcli/walletrpcclient"

	"github.com/decred/dcrd/dcrutil"
)

func getSendSourceAccount(c *walletrpcclient.Client) (uint32, error) {
	// get send  accounts
	accounts, err := c.Balance()
	if err != nil {
		return 0, err
	}

	promptItems := []string{}
	accountItems := map[string]uint32{}
	for _, v := range accounts {
		itemStr := fmt.Sprintf("%s (%s)", v.AccountName, dcrutil.Amount(v.Total).String())
		promptItems = append(promptItems, itemStr)
		accountItems[itemStr] = v.AccountNumber
	}

	prompt := promptui.Select{
		Label: "Select source account",
		Items: promptItems,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("error getting selected account: %s", err.Error())
	}

	account, ok := accountItems[result]
	if !ok {
		return 0, fmt.Errorf("error selecting account")
	}

	return account, nil
}

func getSendDestinationAddress(c *walletrpcclient.Client) (string, error) {
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

	prompt := promptui.Prompt{
		Label:    "Destination Address",
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("error receiving input: %s", err.Error())
	}

	return result, nil
}

func getSendAmount() (int64, error) {
	var amount int64
	var err error

	validate := func(value string) error {
		amount, err = strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("error parsing amount: %s", err.Error())
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Amount (DCR)",
		Validate: validate,
	}

	_, err = prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("error receiving input: %s", err.Error())
	}

	return amount, nil
}

func getWalletPassphrase() (string, error) {
	prompt := promptui.Prompt{
		Label: "Wallet Passphrase",
		Mask:  '*',
	}

	result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("error receiving input: %s", err.Error())
	}
	return result, nil
}
