package walletrpcclient

import (
	"context"
	"errors"
	"fmt"
	"os"

	pb "github.com/decred/dcrwallet/rpc/walletrpc"
	"google.golang.org/grpc"
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

func sendTransaction(conn *grpc.ClientConn, ctx context.Context, opts []string) (*Response, error) {
	c := pb.NewWalletServiceClient(conn)
	var sourceAccount uint32
	var err error
	for {
		err = getSourceAccount(&sourceAccount, c, ctx)
		if err == nil {
			break
		}
		fmt.Printf("error: %s", err.Error())
	}

	var destinationAddress string
	for {
		err = getDestinationAddress(&destinationAddress, c, ctx)
		if err == nil {
			break
		}
		fmt.Printf("error: %s", err.Error())
	}

	var amount int64
	for {
		err = getAmount(&amount, c, ctx)
		if err == nil {
			break
		}
		fmt.Printf("error: %s", err.Error())
	}

	var passphrase string

	// construct transaction
	cReq := &pb.ConstructTransactionRequest{
		SourceAccount: sourceAccount,
	}

	constructResponse, err := c.ConstructTransaction(ctx, cReq)
	if err != nil {
		return nil, fmt.Errorf("Error constructing transaction: %s", err.Error())
	}

	// Sign transaction
	sReq := &pb.SignTransactionRequest{
		Passphrase:            []byte(passphrase),
		SerializedTransaction: constructResponse.UnsignedTransaction,
	}

	signResponse, err := c.SignTransaction(ctx, sReq)
	if err != nil {
		return nil, fmt.Errorf("Error signing transaction: %s", err.Error())
	}

	// publish transaction
	pReq := &pb.PublishTransactionRequest{
		SignedTransaction: signResponse.Transaction,
	}
	publishResponse, err := c.PublishTransaction(ctx, pReq)
	if err != nil {
		return nil, fmt.Errorf("Error publishing transaction")
	}

	res := &Response{
		Columns: []string{"Result", "Hash"},
	}

	resultRow := fmt.Sprintf("%s \t %s",
		"Transaction was published successfully",
		publishResponse.TransactionHash,
	)

	res.Result = []string{resultRow}

	return res, nil
}

/**
func getTransactions(conn *grpc.ClientConn, ctx context.Context, opts []string) (*Response, error) {
	c := pb.NewWalletServiceClient(conn)

	// check if passed options are complete
	if len(opts) < 2 {
		return nil, fmt.Errorf("command 'transactions' requires 2 params. %d found", len(opts))
	}

	// get block height
	startingBlockheight, err := strconv.ParseInt(opts[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Error getting starting block height from options: %s", err.Error())
	}

	limit, err := strconv.ParseInt(opts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Error getting limit from options: %s", err.Error())
	}

	req := &pb.GetTransactionsrequest{}
}
**/

func init() {
	RegisterHandler("send", "Send DCR to address", sendTransaction)
}
