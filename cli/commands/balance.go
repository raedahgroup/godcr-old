package commands

import (
	"fmt"
	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/dcrcli/cli"
	"github.com/raedahgroup/dcrcli/walletrpcclient"
	"github.com/raedahgroup/dcrcli/walletsource"
)

// BalanceCommand displays the user's account balance.
type BalanceCommand struct {
	Detailed bool `short:"d" long:"detailed" description:"Display detailed account balance report"`
}

// Execute runs the `balance` command, displaying the user's account balance.
func (balanceCommand BalanceCommand) Execute(args []string) error {
	accounts, err := cli.WalletSource.AccountsOverview()
	if err != nil {
		return err
	}

	if balanceCommand.Detailed {
		showDetailedBalance(accounts)
	} else {
		showBalanceSummary(accounts)
	}

	return nil
}

func showDetailedBalance(accounts []*walletsource.Account) {
	res := &cli.Response{
		Columns: []string{
			"Account",
			"Total",
			"Spendable",
			"Locked By Tickets",
			"Voting Authority",
			"Unconfirmed",
		},
		Result: make([][]interface{}, len(accounts)),
	}
	for i, account := range accounts {
		res.Result[i] = []interface{}{
			account.Name,
			account.Balance.Total,
			account.Balance.Spendable,
			account.Balance.LockedByTickets,
			account.Balance.VotingAuthority,
			account.Balance.Unconfirmed,
		}
	}

	cli.PrintResult(cli.StdoutWriter, res)
}

func showBalanceSummary(accounts []*walletsource.Account) {
	summarizeBalance := func(total, spendable dcrutil.Amount) string {
		if total == spendable {
			return total.String()
		} else {
			return fmt.Sprintf("Total %s (Spendable %s)", total.String(), spendable.String())
		}
	}

	if len(accounts) == 1 {
		commandOutput := summarizeBalance(accounts[0].Balance.Total, accounts[0].Balance.Spendable)
		cli.PrintStringResult(commandOutput)
	} else {
		commandOutput := make([]string, len(accounts))
		for i, account := range accounts {
			balanceText := summarizeBalance(account.Balance.Total, account.Balance.Spendable)
			commandOutput[i] = fmt.Sprintf("%s \t %s", account.Name, balanceText)
		}
		cli.PrintStringResult(commandOutput...)
	}
}
