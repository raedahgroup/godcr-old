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

func (c *Client) sendTransaction(ctx context.Context, opts []string) (*Response, error) {
	var sourceAccount uint32
	var err error
	for {
		err = getSourceAccount(&sourceAccount, c.wc, ctx)
		if err == nil {
			break
		}
		fmt.Printf("error: %s", err.Error())
	}

	var destinationAddress string
	for {
		err = getDestinationAddress(&destinationAddress, c.wc, ctx)
		if err == nil {
			break
		}
		fmt.Printf("error: %s", err.Error())
	}

	var amount int64
	for {
		err = getAmount(&amount, c.wc, ctx)
		if err == nil {
			break
		}
		fmt.Printf("error: %s", err.Error())
	}

	var passphrase string
	err = getPassphrase(&passphrase)
	if err != nil {
		return nil, fmt.Errorf("error taking passphrase: %s", err.Error())
	}

	// construct transaction
	cReq := &pb.ConstructTransactionRequest{
		SourceAccount: sourceAccount,
	}

	constructResponse, err := c.wc.ConstructTransaction(ctx, cReq)
	if err != nil {
		return nil, fmt.Errorf("Error constructing transaction: %s", err.Error())
	}

	// Sign transaction
	sReq := &pb.SignTransactionRequest{
		Passphrase:            []byte(passphrase),
		SerializedTransaction: constructResponse.UnsignedTransaction,
	}

	signResponse, err := c.wc.SignTransaction(ctx, sReq)
	if err != nil {
		return nil, fmt.Errorf("Error signing transaction: %s", err.Error())
	}

	// publish transaction
	pReq := &pb.PublishTransactionRequest{
		SignedTransaction: signResponse.Transaction,
	}
	publishResponse, err := c.wc.PublishTransaction(ctx, pReq)
	if err != nil {
		return nil, fmt.Errorf("Error publishing transaction")
	}

	res := &Response{
		Columns: []string{"Result", "Hash"},
	}

	resultRow := []interface{}{
		"Transaction was published successfully",
		string(publishResponse.TransactionHash),
	}

	res.Result = [][]interface{}{resultRow}

	return res, nil
}
