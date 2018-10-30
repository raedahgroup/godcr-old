package walletrpcclient

import (
	"context"
	"errors"
	"fmt"
	"strconv"

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
	Handler func(args []string, params map[string]interface{}) (*Response, error)

	Client struct {
		funcMap      map[string]Handler
		commands     map[string]string
		descriptions map[string]string
		wc           pb.WalletServiceClient
		vc           pb.VersionServiceClient
	}
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

// listCommands lists all supported commands
// requires no parameter
func (c *Client) ListSupportedCommands() *Response {
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
	return res
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
func (c *Client) RunCommand(command string, args []string, params map[string]interface{}) (*Response, error) {
	handler := c.funcMap[command]
	res, err := handler(args, params)
	return res, err
}

// listCommands calls the cmdListCommands. This requires no parameters to run
func (c *Client) listCommands(args []string, params map[string]interface{}) (*Response, error) {
	return c.cmdListCommands(context.Background())
}

// receive calls the cmdReceive function.
// parameters:
// 1. where the command is called from. either web or terminal  (required); defaults to terminal
// 2. accountNumber (required when called from web)
func (c *Client) receive(args []string, params map[string]interface{}) (*Response, error) {
	var accountNumber uint32
	caller := "terminal"

	if c, ok := params["caller"]; ok {
		caller = c.(string)
	}

	if caller == "web" {
		if acc, ok := params["accountNumber"]; ok {
			accountNumber = acc.(uint32)
		} else {
			return nil, errors.New("account number is required")
		}
	} else {
		if len(args) == 0 {
			return nil, errors.New("command 'receive' requires at least 1 param. 0 found \nUsage:\n  receive \"accountnumber\"")
		}
		acc, err := strconv.ParseUint(args[0], 0, 32)
		if err != nil {
			return nil, fmt.Errorf("error parsing account number. err:%s", err.Error())
		}
		accountNumber = uint32(acc)
	}

	return c.cmdReceive(context.Background(), accountNumber)
}

func (c *Client) send(args []string, params map[string]interface{}) (*Response, error) {
	var sourceAccount uint32
	var destinationAddress string
	var sendAmount int64
	var passphrase string
	var err error

	caller := "terminal"
	if c, ok := params["caller"]; ok {
		caller = c.(string)
	}

	ctx := context.Background()
	if caller == "web" {
		if sAcc, ok := params["sourceAccount"]; ok {
			sourceAccount = sAcc.(uint32)
		} else {
			return nil, errors.New("source account number is required")
		}

		if addr, ok := params["destinationAddress"]; ok {
			destinationAddress = addr.(string)
		} else {
			return nil, errors.New("destination address is required")
		}

		if amt, ok := params["amount"]; ok {
			sendAmount = amt.(int64)
		} else {
			return nil, errors.New("send amoount is required")
		}

		if pass, ok := params["passphrase"]; ok {
			passphrase = pass.(string)
		} else {
			return nil, errors.New("passphrase is required")
		}
	} else {
		sourceAccount, err = getSendSourceAccount(c.wc, ctx)
		if err != nil {
			return nil, err
		}

		destinationAddress, err = getSendDestinationAddress(c.wc, ctx)
		if err != nil {
			return nil, err
		}

		sendAmount, err = getSendAmount(c.wc, ctx)
		if err != nil {
			return nil, err
		}

		passphrase, err = getWalletPassphrase(c.wc, ctx)
		if err != nil {
			return nil, err
		}
	}

	return c.cmdSendTransaction(ctx, sourceAccount, destinationAddress, sendAmount, passphrase)
}

func (c *Client) balance(args []string, params map[string]interface{}) (*Response, error) {
	return c.cmdBalance(context.Background())
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

func (c *Client) registerHandlers() {
	c.RegisterHandler("listcommands", "-l", "List all supported commands", c.listCommands)
	c.RegisterHandler("receive", "receive", "Generate address to receive funds", c.receive)
	c.RegisterHandler("send", "send", "Send DCR to address. Multi-step", c.send)
	c.RegisterHandler("balance", "balance", "Check balance of an account", c.balance)
}
