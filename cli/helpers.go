package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/raedahgroup/dcrcli/walletrpcclient"

	"github.com/decred/dcrd/dcrutil"
)

func getSendSourceAccount(c *walletrpcclient.Client, promptFn prompter) (uint32, error) {
	var choice int
	var err error
	// get send  accounts
	accounts, err := c.Balance()
	if err != nil {
		return 0, err
	}
	// validateInput ensures that the input received is a number that corresponds to an account
	validateInput := func(value string) error {
		if choice, err = strconv.Atoi(value); err != nil {
			return fmt.Errorf("could not recognize input: not an allowed option")
		}
		choiceFloor, choiceCeiling := 0, len(accounts)-1
		if choice < choiceFloor || choice > choiceCeiling {
			return fmt.Errorf("%d is not an allowed option", choice)
		}
		return nil
	}

	promptItems := promptData{Label: "Select source account"}
	accountItems := map[int]uint32{}
	options := make([]string, len(accounts))

	for idx, v := range accounts {
		options[idx] = fmt.Sprintf("%s (%s)", v.AccountName, dcrutil.Amount(v.Total).String())
		accountItems[idx] = v.AccountNumber
	}
	promptItems.Options = options

	for {
		input, err := promptFn(promptItems)
		if err != nil {
			// There was an error reading input; we cannot proceed.
			return 0, fmt.Errorf("error getting selected account: %s", err.Error())
		}
		if err = validateInput(input); err != nil {
			// The user entered incorrect data. Prompt again.
			fmt.Println(err.Error())
			continue
		}
		break
	}
	return accountItems[choice], nil
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

	prompt := promptData{
		Label: "Destination Address",
	}
	var result string
	for {
		result, err := promptFn(prompt)
		if err != nil {
			// There was an error reading input; we cannot proceed.
			return "", fmt.Errorf("error receiving input: %s", err.Error())
		}
		if err = validate(result); err != nil {
			// The user entered incorrect data. Prompt again.
			fmt.Println(err.Error())
			continue
		}
		break
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

	prompt := promptData{Label: "Amount (DCR)"}

	for {
		result, err := promptFn(prompt)
		if err != nil {
			// There was an error reading input; we cannot proceed.
			return 0, fmt.Errorf("error receiving input: %s", err.Error())
		}
		if err = validate(result); err != nil {
			// The user entered incorrect data. Prompt again.
			fmt.Println(err.Error())
			continue
		}
		break
	}
	return amount, nil
}

func getWalletPassphrase(promptFn prompter) (string, error) {
	prompt := promptData{
		Label:  "Wallet Passphrase",
		Secure: true,
	}

	result, err := promptFn(prompt)
	if err != nil {
		return "", fmt.Errorf("error receiving input: %s", err.Error())
	}
	return result, nil
}
