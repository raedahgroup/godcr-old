package walletrpcclient

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/decred/dcrwallet/rpc/walletrpc"
	"google.golang.org/grpc"
)

var (
	balanceColumns = []string{
		"Total",
		"Spendable",
		"Locked By Tickets",
		"Voting Authority",
		"Unconfirmed",
	}
)

// balance gets the balance of an account by its account number
// requires two parameters. 1. Account number 2. Required confirmations
// returns an error if any of the required parameters are absent.
// returns an error if any of the parameters passed in cannot be converted to their required types
// for transport
func balance(conn *grpc.ClientConn, ctx context.Context, opts []string) (*Response, error) {
	c := pb.NewWalletServiceClient(conn)

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
	r, err := c.Balance(ctx, req)
	if err != nil {
		return nil, err
	}

	res := &Response{
		Columns: balanceColumns,
	}

	row := fmt.Sprintf("%d \t %d \t %d \t %d \t %d",
		r.Total,
		r.Spendable,
		r.LockedByTickets,
		r.VotingAuthority,
		r.Unconfirmed,
	)

	res.Result = append(res.Result, row)
	return res, nil
}

func init() {
	RegisterHandler("balance", "Show account balance for account number", balance)
}
