package walletrpcclient

import (
	"context"
	"fmt"

	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/txscript"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	walletServiceClient pb.WalletServiceClient
}

func New(address, cert string, noTLS bool) (*Client, error) {
	c := &Client{}
	conn, err := c.connect(address, cert, noTLS)
	if err != nil {
		return nil, err
	}

	// register clients
	c.walletServiceClient = pb.NewWalletServiceClient(conn)

	return c, nil
}

func (c *Client) connect(address, cert string, noTLS bool) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	if noTLS {
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
	} else {
		creds, err := credentials.NewClientTLSFromFile(cert, "")
		if err != nil {
			return nil, err
		}

		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(
				creds,
			),
		}

		conn, err = grpc.Dial(address, opts...)
		if err != nil {
			return nil, err
		}
	}

	return conn, nil
}

type SendResult struct {
	TransactionHash string `json:"transaction_hash"`
}

func (c *Client) Send(amount int64, sourceAccount uint32, destinationAddress, passphrase string) (*SendResult, error) {
	// decode destination address
	addr, err := dcrutil.DecodeAddress(destinationAddress)
	if err != nil {
		return nil, fmt.Errorf("error decoding destination address: %s", err.Error())
	}

	// get script
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	// construct transaction
	constructRequest := &pb.ConstructTransactionRequest{
		SourceAccount: sourceAccount,
		NonChangeOutputs: []*pb.ConstructTransactionRequest_Output{{
			Destination: &pb.ConstructTransactionRequest_OutputDestination{
				Script:        pkScript,
				ScriptVersion: 0,
			},
			Amount: amount,
		}},
	}

	ctx := context.Background()

	constructResponse, err := c.walletServiceClient.ConstructTransaction(ctx, constructRequest)
	if err != nil {
		return nil, fmt.Errorf("error constructing transaction: %s", err.Error())
	}

	// sign transaction
	signRequest := &pb.SignTransactionRequest{
		Passphrase:            []byte(passphrase),
		SerializedTransaction: constructResponse.UnsignedTransaction,
	}

	signResponse, err := c.walletServiceClient.SignTransaction(ctx, signRequest)
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %s", err.Error())
	}

	// publish transaction
	publishRequest := &pb.PublishTransactionRequest{
		SignedTransaction: signResponse.Transaction,
	}

	publishResponse, err := c.walletServiceClient.PublishTransaction(ctx, publishRequest)
	if err != nil {
		return nil, fmt.Errorf("error publishing transaction: %s", err.Error())
	}

	response := &SendResult{
		TransactionHash: string(publishResponse.TransactionHash),
	}

	return response, nil
}

type AccountBalanceResult struct {
	AccountName     string         `json:"account_name"`
	AccountNumber   uint32         `json:"account_number"`
	Total           dcrutil.Amount `json:"total"`
	Spendable       dcrutil.Amount `json:"spendable"`
	LockedByTickets dcrutil.Amount `json:"locked_by_tickets"`
	VotingAuthority dcrutil.Amount `json:"voting_authority"`
	Unconfirmed     dcrutil.Amount `json:"unconfirmed"`
}

func (c *Client) Balance() ([]AccountBalanceResult, error) {
	ctx := context.Background()
	accounts, err := c.walletServiceClient.Accounts(ctx, &pb.AccountsRequest{})
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts: %s", err.Error())
	}

	balanceResult := make([]AccountBalanceResult, len(accounts.Accounts))
	for i, v := range accounts.Accounts {
		req := &pb.BalanceRequest{
			AccountNumber:         v.AccountNumber,
			RequiredConfirmations: 0,
		}

		res, err := c.walletServiceClient.Balance(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("error fetching balance for account: %d :%s", v.AccountNumber, err.Error())
		}

		balanceResult[i] = AccountBalanceResult{
			AccountNumber:   v.AccountNumber,
			AccountName:     v.AccountName,
			Total:           dcrutil.Amount(res.Total),
			Spendable:       dcrutil.Amount(res.Spendable),
			LockedByTickets: dcrutil.Amount(res.LockedByTickets),
			VotingAuthority: dcrutil.Amount(res.VotingAuthority),
			Unconfirmed:     dcrutil.Amount(res.Unconfirmed),
		}
	}

	return balanceResult, nil
}

type ReceiveResult struct {
	Address string
}

func (c *Client) Receive(accountNumber uint32) (*ReceiveResult, error) {
	ctx := context.Background()

	req := &pb.NextAddressRequest{
		Account:   accountNumber,
		GapPolicy: pb.NextAddressRequest_GAP_POLICY_WRAP,
		Kind:      pb.NextAddressRequest_BIP0044_EXTERNAL,
	}

	r, err := c.walletServiceClient.NextAddress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error generating receive address: %s", err.Error())
	}

	res := &ReceiveResult{
		Address: r.Address,
	}

	return res, nil
}

func (c *Client) ValidateAddress(address string) (bool, error) {
	req := &pb.ValidateAddressRequest{
		Address: address,
	}

	r, err := c.walletServiceClient.ValidateAddress(context.Background(), req)
	if err != nil {
		return false, err
	}

	return r.IsValid, nil
}
