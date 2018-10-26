package walletrpcclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"

	"github.com/decred/dcrd/dcrutil"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
)

func getSourceAccount(c pb.WalletServiceClient, ctx context.Context) (string, error) {
	// get accounts
	accountsRes, err := c.Accounts(ctx, &pb.AccountsRequest{})
	if err != nil {
		return "", fmt.Errorf("error fetching accounts. err: %s", err.Error())
	}

	promptItems := make([]string, len(accountsRes.Accounts))
	for _, v := range accountsRes.Accounts {
		balanceReq := &pb.BalanceRequest{
			AccountNumber:         v.AccountNumber,
			RequiredConfirmations: 0,
		}

		balanceRes, err := c.Balance(ctx, balanceReq)
		if err != nil {
			return "", fmt.Errorf("error fetching balance for account: %d. err: %s", v.AccountNumber, err.Error())
		}

		item := fmt.Sprintf("%s (%s)", v.AccountName, dcrutil.Amount(balanceRes.Total).String())
		promptItems = append(promptItems, item)

	}

	prompt := promptui.Select{
		Label: "Select source account",
		Items: promptItems,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", fmt.Errorf("error getting selected account: %s", err.Error())
	}

	parts := strings.Split(result, "(")
	if len(parts) == 0 {
		return "", fmt.Errorf("error selecting source account")
	}

	return parts[0], nil
}

func getDestinationAddress(destinationAddress *string, c pb.WalletServiceClient, ctx context.Context) error {
	fmt.Println("Destination address: ")
	_, err := fmt.Scanln(destinationAddress)
	if err != nil {
		return err
	}

	// validate address
	req := &pb.ValidateAddressRequest{
		Address: *destinationAddress,
	}
	r, err := c.ValidateAddress(ctx, req)
	if err != nil || !r.IsValid {
		return fmt.Errorf("Invalid address")
	}

	return nil
}

func getAmount(amount *int64, c pb.WalletServiceClient, ctx context.Context) error {
	fmt.Println("Amount; ")
	_, err := fmt.Scanf("%d", amount)
	return err
}

func getPassphrase(passphrase *string) error {
	fmt.Println("Wallet Passphrase: ")
	_, err := fmt.Scanln(passphrase)
	return err
}
