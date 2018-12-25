package commands

import (
	"github.com/raedahgroup/dcrcli/cli"
)

// BalanceCommand displays the user's account balance.
type BalanceCommand struct{}

// Execute runs the `balance` command, displaying the user's account balance.
func (b BalanceCommand) Execute(args []string) error {
	balances, err := cli.WalletClient.Balance()
	if err != nil {
		return err
	}

	res := &cli.Response{
		Columns: []string{
			"Account",
			"Total",
			"Spendable",
			"Locked By Tickets",
			"Voting Authority",
			"Unconfirmed",
		},
		Result: make([][]interface{}, len(balances)),
	}

	for i, v := range balances {
		res.Result[i] = []interface{}{
			v.AccountName,
			v.Total,
			v.Spendable,
			v.LockedByTickets,
			v.VotingAuthority,
			v.Unconfirmed,
		}
	}

	cli.PrintResult(cli.StdoutWriter, res)
	return nil
}
