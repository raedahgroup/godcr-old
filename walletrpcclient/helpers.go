package walletrpcclient

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/manifoldco/promptui"

	"github.com/decred/dcrd/dcrutil"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
)

func getSendSourceAccount(c pb.WalletServiceClient, ctx context.Context) (uint32, error) {
	// get accounts
	accountsRes, err := c.Accounts(ctx, &pb.AccountsRequest{})
	if err != nil {
		return 0, fmt.Errorf("error fetching accounts: %s", err.Error())
	}

	promptItems := []string{}
	accounts := map[string]uint32{}
	for _, v := range accountsRes.Accounts {
		balanceReq := &pb.BalanceRequest{
			AccountNumber:         v.AccountNumber,
			RequiredConfirmations: 0,
		}

		balanceRes, err := c.Balance(ctx, balanceReq)
		if err != nil {
			return 0, fmt.Errorf("error fetching balance for account %d: %s", v.AccountNumber, err.Error())
		}

		item := fmt.Sprintf("%s (%s)", v.AccountName, dcrutil.Amount(balanceRes.Total).String())
		fmt.Println(v.AccountNumber)
		promptItems = append(promptItems, item)
		accounts[item] = v.AccountNumber
	}

	prompt := promptui.Select{
		Label: "Select source account",
		Items: promptItems,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return 0, fmt.Errorf("error getting selected account: %s", err.Error())
	}

	account, ok := accounts[result]
	if !ok {
		return 0, fmt.Errorf("error selecting account")
	}

	return account, nil
}

func getSendDestinationAddress(c pb.WalletServiceClient, ctx context.Context) (string, error) {
	validate := func(address string) error {
		req := &pb.ValidateAddressRequest{
			Address: address,
		}

		r, err := c.ValidateAddress(ctx, req)
		if err != nil {
			return fmt.Errorf("error validating address: %s", err.Error())
		}

		if !r.IsValid {
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

func getSendAmount(c pb.WalletServiceClient, ctx context.Context) (int64, error) {
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

func getWalletPassphrase(c pb.WalletServiceClient, ctx context.Context) (string, error) {
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
