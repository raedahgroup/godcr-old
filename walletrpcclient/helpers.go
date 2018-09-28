package walletrpcclient

import (
	"context"
	"errors"
	"fmt"
	"os"

	pb "github.com/decred/dcrwallet/rpc/walletrpc"
)

func getSourceAccount(fromAccount *uint32, c pb.WalletServiceClient, ctx context.Context) error {
	fmt.Println("Source Account: ")
	_, err := fmt.Scanf("%d", fromAccount)
	if err != nil {
		return err
	}

	// validate account number
	r, err := c.Accounts(ctx, &pb.AccountsRequest{})
	if err != nil {
		fmt.Printf("Error validating account; %s", err.Error())
		os.Exit(1)
	}

	for _, v := range r.Accounts {
		if v.AccountNumber == *fromAccount {
			return nil
		}
	}
	return errors.New("invalid account number")
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
