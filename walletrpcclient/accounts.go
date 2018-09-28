package walletrpcclient

import (
	"context"

	pb "github.com/decred/dcrwallet/rpc/walletrpc"
)

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
