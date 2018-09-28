package walletrpcclient

import (
	"context"
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

func (c *Client) registerHandlers() {
	c.RegisterHandler("send", "Send DCR to address", c.sendTransaction)
	c.RegisterHandler("accounts", "List all accounts", c.accounts)
	c.RegisterHandler("balance", "Show account balance for account number", c.balance)
}
