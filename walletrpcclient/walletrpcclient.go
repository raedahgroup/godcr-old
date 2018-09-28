package walletrpcclient

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pb "github.com/decred/dcrwallet/rpc/walletrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	Response struct {
		Columns []string
		Result  [][]interface{}
	}
	Handler func(ctx context.Context, args []string) (*Response, error)

	Client struct {
		funcMap  map[string]Handler
		commands map[string]string
		wc       pb.WalletServiceClient
	}
)

func New(address, cert string, noTLS bool) (*Client, error) {
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

		// dial options
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

	client := &Client{
		wc:       pb.NewWalletServiceClient(conn),
		funcMap:  make(map[string]Handler),
		commands: make(map[string]string),
	}

	client.registerHandlers()

	return client, nil
}

// IsCommandSupported returns a boolean whose value depends on if a command is registered as suppurted along
// with it's func handler
func (c *Client) IsCommandSupported(command string) bool {
	_, ok := c.funcMap[command]
	return ok
}

// RunCommand takes a command and tries to call the appropriate handler to call a gRPC service
// This should only be called after verifying that the command is supported using the IsCommandSupported
// function.
func (c *Client) RunCommand(command string, opts []string) (*Response, error) {
	handler := c.funcMap[command]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	res, err := handler(ctx, opts)
	return res, err
}

// RegisterHandler registers a command, its description and its handler
func (c *Client) RegisterHandler(key, description string, h Handler) {
	if _, ok := c.funcMap[key]; ok {
		panic("trying to register a handler twice: " + key)
	}

	c.funcMap[key] = h
	c.commands[key] = description
}

// accounts lists all the accounts and their total balances in a wallet
// no options are required to perform this operation
func (c *Client) accounts(ctx context.Context, opts []string) (*Response, error) {
	r, err := c.wc.Accounts(ctx, &pb.AccountsRequest{})
	if err != nil {
		return nil, err
	}

	accountsColumn := []string{
		"Account Number",
		"Account Name",
		"Total Balance",
	}

	res := &Response{
		Columns: accountsColumn,
		Result:  [][]interface{}{},
	}

	for _, v := range r.Accounts {
		row := []interface{}{v.AccountNumber, v.AccountName, v.TotalBalance}
		res.Result = append(res.Result, row)
	}

	return res, nil
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

// balance gets the balance of an account by its account number
// requires two parameters. 1. Account number 2. Required confirmations
// returns an error if any of the required parameters are absent.
// returns an error if any of the parameters passed in cannot be converted to their required types
// for transport
func (c *Client) balance(ctx context.Context, opts []string) (*Response, error) {
	// check if passed options are complete
	if len(opts) < 2 {
		return nil, fmt.Errorf("command 'balance' requires 2 params. %d found", len(opts))
	}

	accountNumber, err := strconv.ParseUint(opts[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error getting account number from options: %s", err.Error())
	}

	requiredConfirmations, err := strconv.ParseInt(opts[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error getting required confirmations from options: %s", err.Error())
	}
	req := &pb.BalanceRequest{
		AccountNumber:         uint32(accountNumber),
		RequiredConfirmations: int32(requiredConfirmations),
	}
	r, err := c.wc.Balance(ctx, req)
	if err != nil {
		return nil, err
	}

	balanceColumns := []string{
		"Total",
		"Spendable",
		"Locked By Tickets",
		"Voting Authority",
		"Unconfirmed",
	}

	res := &Response{
		Columns: balanceColumns,
	}

	row := []interface{}{
		r.Total,
		r.Spendable,
		r.LockedByTickets,
		r.VotingAuthority,
		r.Unconfirmed,
	}

	res.Result = append(res.Result, row)
	return res, nil
}

func (c *Client) registerHandlers() {
	c.RegisterHandler("send", "Send DCR to address", c.sendTransaction)
	c.RegisterHandler("accounts", "List all accounts", c.accounts)
	c.RegisterHandler("balance", "Show account balance for account number", c.balance)
}
