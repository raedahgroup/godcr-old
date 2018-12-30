package commands

import (
	"fmt"

	"github.com/decred/dcrd/dcrutil"
	"github.com/raedahgroup/godcr/cli/termio"
	"github.com/raedahgroup/godcr/walletrpcclient"
)

// BalanceCommand displays the user's account balance.
type BalanceCommand struct {
	CommanderStub
	Detailed bool `short:"d" long:"detailed" description:"Display detailed account balance report"`
}

// Run runs the `balance` command, displaying the user's account balance.
func (balanceCommand BalanceCommand) Run(client *walletrpcclient.Client, args []string) error {
	accountBalances, err := client.Balance()
	if err != nil {
		return err
	}

	if balanceCommand.Detailed {
		showDetailedBalance(accountBalances)
	} else {
		showBalanceSummary(accountBalances)
	}

	return nil
}

func showDetailedBalance(accountBalances []*walletrpcclient.AccountBalanceResult) {
	columns := []string{
		"Account",
		"Total",
		"Spendable",
		"Locked By Tickets",
		"Voting Authority",
		"Unconfirmed",
	}
	rows := make([][]interface{}, len(accountBalances))
	for i, account := range accountBalances {
		rows[i] = []interface{}{
			account.AccountName,
			account.Total,
			account.Spendable,
			account.LockedByTickets,
			account.VotingAuthority,
			account.Unconfirmed,
		}
	}

	termio.PrintTabularResult(termio.StdoutWriter, columns, rows)
}

func showBalanceSummary(accountBalances []*walletrpcclient.AccountBalanceResult) {
	summarizeBalance := func(total, spendable dcrutil.Amount) string {
		if total == spendable {
			return total.String()
		} else {
			return fmt.Sprintf("Total %s (Spendable %s)", total.String(), spendable.String())
		}
	}

	if len(accountBalances) == 1 {
		commandOutput := summarizeBalance(accountBalances[0].Total, accountBalances[0].Spendable)
		termio.PrintStringResult(commandOutput)
	} else {
		commandOutput := make([]string, len(accountBalances))
		for i, accountBalance := range accountBalances {
			balanceText := summarizeBalance(accountBalance.Total, accountBalance.Spendable)
			commandOutput[i] = fmt.Sprintf("%s \t %s", accountBalance.AccountName, balanceText)
		}
		termio.PrintStringResult(commandOutput...)
	}
}
