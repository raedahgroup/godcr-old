package walletrpcclient

import (
	"context"
	"fmt"

	pb "github.com/decred/dcrwallet/rpc/walletrpc"
	"google.golang.org/grpc"
)

var (
	accountsColumn = []string{
		"Account Number",
		"Account Name",
		"Total Balance",
	}
)

// accounts lists all the accounts and their total balances in a wallet
// no options are required to perform this operation
func accounts(conn *grpc.ClientConn, ctx context.Context, opts []string) (*Response, error) {
	c := pb.NewWalletServiceClient(conn)
	r, err := c.Accounts(ctx, &pb.AccountsRequest{})
	if err != nil {
		return nil, err
	}

	res := &Response{
		Columns: accountsColumn,
	}

	for _, v := range r.Accounts {
		row := fmt.Sprintf("%d \t %s \t %d", v.AccountNumber, v.AccountName, v.TotalBalance)
		res.Result = append(res.Result, row)
	}

	return res, nil
}

func init() {
	RegisterHandler("accounts", "List all accounts", accounts)
}
