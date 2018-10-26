package walletrpcclient

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/decred/dcrd/txscript"

	"github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/decred/dcrd/dcrutil"
	pb "github.com/decred/dcrwallet/rpc/walletrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type (
	Response struct {
		Columns []string
		Result  [][]interface{}
		Qrcode  bool
	}
	Handler func(ctx context.Context, args []string) (*Response, error)

	Client struct {
		funcMap      map[string]Handler
		commands     map[string]string
		descriptions map[string]string
		wc           pb.WalletServiceClient
		vc           pb.VersionServiceClient
	}
)

const (
	requiredConfirmations int32 = 0
)

func New() *Client {
	client := &Client{
		funcMap:      make(map[string]Handler),
		commands:     make(map[string]string),
		descriptions: make(map[string]string),
	}

	client.registerHandlers()
	return client
}

func (c *Client) Connect(address, cert string, noTLS bool) error {
	var conn *grpc.ClientConn
	var err error

	if noTLS {
		conn, err = grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return err
		}
	} else {
		creds, err := credentials.NewClientTLSFromFile(cert, "")
		if err != nil {
			return err
		}

		// dial options
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(
				creds,
			),
		}

		conn, err = grpc.Dial(address, opts...)
		if err != nil {
			return err
		}
	}

	c.wc = pb.NewWalletServiceClient(conn)
	c.vc = pb.NewVersionServiceClient(conn)
	return nil
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
func (c *Client) RegisterHandler(key, command, description string, h Handler) {
	if _, ok := c.funcMap[key]; ok {
		panic("trying to register a handler twice: " + key)
	}

	c.funcMap[key] = h
	c.commands[key] = command
	c.descriptions[key] = description
}

func (c *Client) sendTransaction(ctx context.Context, opts []string) (*Response, error) {
	sourceAccount, err := getSendSourceAccount(c.wc, ctx)
	if err != nil {
		return nil, err
	}

	destinationAddressStr, err := getSendDestinationAddress(c.wc, ctx)
	if err != nil {
		return nil, err
	}

	// decode address
	addr, err := dcrutil.DecodeAddress(destinationAddressStr)
	if err != nil {
		return nil, err
	}

	amount, err := getSendAmount(c.wc, ctx)
	if err != nil {
		return nil, err
	}

	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		return nil, err
	}

	cReq := &pb.ConstructTransactionRequest{
		SourceAccount: sourceAccount,
		NonChangeOutputs: []*pb.ConstructTransactionRequest_Output{{
			Destination: &pb.ConstructTransactionRequest_OutputDestination{
				Script:        pkScript,
				ScriptVersion: 0,
			},
			Amount: amount,
		}},
	}

	cRes, err := c.wc.ConstructTransaction(ctx, cReq)
	if err != nil {
		return nil, fmt.Errorf("Error constructing transaction: %s", err.Error())
	}

	passphrase, err := getWalletPassphrase(c.wc, ctx)
	if err != nil {
		return nil, err
	}

	// Sign transaction
	sReq := &pb.SignTransactionRequest{
		Passphrase:            []byte(passphrase),
		SerializedTransaction: cRes.UnsignedTransaction,
	}

	sRes, err := c.wc.SignTransaction(ctx, sReq)
	if err != nil {
		return nil, fmt.Errorf("error signing transaction: %s", err.Error())
	}

	// publish transaction
	pReq := &pb.PublishTransactionRequest{
		SignedTransaction: sRes.Transaction,
	}

	pRes, err := c.wc.PublishTransaction(ctx, pReq)
	if err != nil {
		return nil, fmt.Errorf("error publishing transaction: %s", err.Error())
	}

	res := &Response{
		Columns: []string{"Result", "Hash"},
	}

	resultRow := []interface{}{
		"Successfull",
		string(pRes.TransactionHash),
	}

	res.Result = [][]interface{}{resultRow}

	return res, nil
}

// balance gets the balance of an account by its account number
// requires at least one option (AccountNumber).
// the second paramter (minConf) is optional and defaults to 0 if not set
// returns an error if any of the parameters passed in cannot be converted to their required types
// for transport
func (c *Client) balance(ctx context.Context, opts []string) (*Response, error) {
	accountsRes, err := c.wc.Accounts(ctx, &pb.AccountsRequest{})
	if err != nil {
		return nil, fmt.Errorf("error fetching accounts. err: %s", err.Error())
	}

	balances := make([][]interface{}, len(accountsRes.Accounts))
	for i, v := range accountsRes.Accounts {
		balanceReq := &pb.BalanceRequest{
			AccountNumber:         v.AccountNumber,
			RequiredConfirmations: 0,
		}

		balanceRes, err := c.wc.Balance(ctx, balanceReq)
		if err != nil {
			return nil, fmt.Errorf("error fetching balance for account: %d. err: %s", v.AccountNumber, err.Error())
		}

		balances[i] = []interface{}{
			v.AccountName,
			dcrutil.Amount(balanceRes.Total),
			dcrutil.Amount(balanceRes.Spendable),
			dcrutil.Amount(balanceRes.LockedByTickets),
			dcrutil.Amount(balanceRes.VotingAuthority),
			dcrutil.Amount(balanceRes.Unconfirmed),
		}
	}

	balanceColumns := []string{
		"Account",
		"Total",
		"Spendable",
		"Locked By Tickets",
		"Voting Authority",
		"Unconfirmed",
	}

	res := &Response{
		Columns: balanceColumns,
		Result:  balances,
	}

	return res, nil
}

// receive returns a generated address, and generates a qr code for recieving funds
// requires no parameter
func (c *Client) receive(ctx context.Context, opts []string) (*Response, error) {
	if len(opts) == 0 {
		return nil, errors.New("command 'receive' requires at least 1 param. 0 found")
	}

	accountNumber, err := strconv.ParseUint(opts[0], 0, 32)
	if err != nil {
		return nil, fmt.Errorf("error parsing account number. err:%s", err.Error())
	}

	// TODO this should be optionally supplied by the user
	gapPolicy := pb.NextAddressRequest_GAP_POLICY_WRAP
	// this shouldn't
	kind := pb.NextAddressRequest_BIP0044_EXTERNAL

	req := &pb.NextAddressRequest{
		Account:   uint32(accountNumber),
		GapPolicy: gapPolicy,
		Kind:      kind,
	}

	r, err := c.wc.NextAddress(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error getting receive address: err: %s", err.Error())
	}

	res := &Response{
		Columns: []string{
			"Address",
			"QR Code",
		},
		Result: [][]interface{}{
			[]interface{}{
				r.Address,
				"",
			},
		},
	}
	obj := qrcodeTerminal.New()
	obj.Get(r.Address).Print()

	return res, nil
}

// walletVersion fetches and returns version of wallet we are connected to
func (c *Client) walletVersion(ctx context.Context, opts []string) (*Response, error) {
	r, err := c.vc.Version(ctx, &pb.VersionRequest{})
	if err != nil {
		return nil, err
	}

	res := &Response{
		Columns: []string{
			"Version",
		},
		Result: [][]interface{}{
			[]interface{}{r.VersionString},
		},
	}
	return res, nil
}

// listCommands lists all supported commands
// requires no parameter
func (c *Client) listCommands(ctx context.Context, opts []string) (*Response, error) {
	res := &Response{
		Columns: []string{"Command", "Description"},
	}
	for i, v := range c.commands {
		item := []interface{}{
			v,
			c.descriptions[i],
		}
		res.Result = append(res.Result, item)
	}
	return res, nil
}

func (c *Client) registerHandlers() {
	c.RegisterHandler("listcommands", "-l", "List all supported commands", c.listCommands)
	c.RegisterHandler("send", "send", "Send DCR to address. Multi-step", c.sendTransaction)
	c.RegisterHandler("walletversion", "walletversion", "Show version of wallet", c.walletVersion)
	c.RegisterHandler("balance", "balance", "Check balance of an account", c.balance)
	c.RegisterHandler("receive", "receive", "Generate address to receive funds", c.receive)
}
